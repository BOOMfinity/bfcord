// Low-Level Discord API requests

package api

import (
	"bytes"
	"context"
	"fmt"
	"github.com/BOOMfinity/bfcord/errs"
	"reflect"
	"sync"
	"time"

	"github.com/BOOMfinity/bfcord/api/limits"
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/golog"
	"github.com/andersfylling/snowflake/v5"
	"github.com/segmentio/encoding/json"
	"github.com/valyala/fasthttp"
)

var contentJson = []byte("application/json")
var discordPrefix = []byte("https://discord.com/api/")

var requestDataPool = sync.Pool{
	New: func() interface{} {
		return &RequestData{}
	},
}

type Client struct {
	manager      *limits.Manager
	logger       golog.Logger
	token        string
	staticHeader []byte
	retryDelay   time.Duration
	timeout      time.Duration
	maxRetries   uint8
}

func (v *Client) Webhook(id snowflake.ID, token string) discord.WebhookQuery {
	return webhookQuery{api: v, id: id, token: token}
}

func (v *Client) Channel(id snowflake.ID) discord.ChannelQuery {
	return NewChannelQuery(v, id)
}

func (v *Client) LowLevel() discord.LowLevelClientQuery {
	d := &lowLevelQuery{client: v}
	d.emptyOptions = emptyOptions[discord.LowLevelClientQuery]{data: d}
	return d
}

func (v Client) Log() golog.Logger {
	return v.logger
}

func (v Client) acquireRequestData() *RequestData {
	data := requestDataPool.Get().(*RequestData)
	data.timeout = v.timeout
	data.retries = v.maxRetries
	data.retryDelay = v.retryDelay
	return data
}

func (v Client) New(auth bool) *fasthttp.Request {
	req := fasthttp.AcquireRequest()
	req.Header.AddBytesV(fasthttp.HeaderContentType, contentJson)
	if auth && len(v.staticHeader) > 0 {
		req.Header.AddBytesV(fasthttp.HeaderAuthorization, v.staticHeader)
	}
	return req
}

func (v Client) DoNoResp(req *fasthttp.Request, options ...Option) error {
	res, err := v.Do(req, options...)
	if res != nil {
		fasthttp.ReleaseResponse(res)
	}
	return err
}

func (v Client) DoResult(req *fasthttp.Request, result any, options ...Option) error {
	if reflect.ValueOf(result).Kind() != reflect.Ptr {
		return errs.ResultMustBePointer
	}
	res, err := v.Do(req, options...)
	if err != nil {
		return err
	}
	defer fasthttp.ReleaseResponse(res)
	err = json.Unmarshal(res.Body(), &result)
	if err != nil {
		return err
	}
	return nil
}

func (v Client) DoBytes(req *fasthttp.Request, options ...Option) ([]byte, error) {
	res, err := v.Do(req, options...)
	if err != nil {
		return nil, err
	}
	defer fasthttp.ReleaseResponse(res)
	return res.Body(), nil
}

func (v Client) Do(req *fasthttp.Request, options ...Option) (*fasthttp.Response, error) {
	defer fasthttp.ReleaseRequest(req)
	data := v.acquireRequestData()
	defer requestDataPool.Put(data)
	for i := range options {
		options[i](data)
	}
	res := fasthttp.AcquireResponse()
	for i := 0; i < int(data.retries+1); i++ {
		if i > 0 {
			time.Sleep(data.retryDelay)
		}
		if bytes.HasPrefix(req.URI().FullURI(), discordPrefix) {
			ctx, cancel := context.WithTimeout(context.Background(), v.timeout)
			err := v.manager.Wait(ctx, string(req.URI().FullURI()))
			if err != nil {
				cancel()
				return nil, err
			}
			cancel()
		}
		_time := time.Now()
		err := fasthttp.DoTimeout(req, res, data.timeout)
		if err != nil {
			return nil, err
		}
		if res.StatusCode() < 400 {
			v.logger.Debug().Any(string(req.URI().Path())).Any(time.Since(_time)).Send("Response body: %vB", len(res.Body()))
		}
		if (res.StatusCode() >= 400 && res.StatusCode() < 500) && !bytes.HasPrefix(req.URI().FullURI(), discordPrefix) {
			v.logger.Warn().Any(string(req.URI().Path())).Any(time.Since(_time)).Send("Received status %v with %vB of body size", res.StatusCode(), len(res.Body()))
		}
		if res.StatusCode() >= 500 {
			v.logger.Error().Any(string(req.URI().Path())).Any(time.Since(_time)).Send("Internal server error :(")
			continue
		}
		if bytes.HasPrefix(req.URI().FullURI(), discordPrefix) {
			err = v.manager.ParseRequest(req, res)
			if err != nil {
				return nil, err
			}
			if res.StatusCode() == fasthttp.StatusTooManyRequests {
				v.logger.Debug().Any(string(req.URI().Path())).Any(time.Since(_time)).Send("Just got rate limited :(")
				data.retries--
				continue
			}
			if res.StatusCode() >= 400 && json.Valid(res.Body()) {
				switch res.StatusCode() {
				case 401:
					err = errs.HTTPUnauthorized
				case 404:
					err = errs.HTTPNotFound
				default:
					dcErr := new(errs.DiscordError)
					err = json.Unmarshal(res.Body(), &dcErr)
					if err != nil {
						return nil, fmt.Errorf("failed to parse discord error: %w", err)
					}
					err = dcErr
				}
				v.logger.Warn().Any(string(req.URI().Path())).Any(time.Since(_time)).Send(err.Error())
				return nil, err
			}
		}
		return res, nil
	}
	return nil, errs.TooManyRetries
}

func (v *Client) User(id snowflake.ID) discord.UserQuery {
	return NewUserQuery(v, id)
}

func (v *Client) Guild(id snowflake.ID) discord.GuildQuery {
	return NewGuildQuery(v, id)
}

func (v Client) CurrentUser() (user *discord.User, err error) {
	req := v.New(true)
	req.SetRequestURI(fmt.Sprintf("%v/users/@me", FullApiUrl))
	err = v.DoResult(req, &user)
	return
}

func (v Client) GatewayURL() (url string, err error) {
	gateway := struct {
		Url string `json:"url"`
	}{}
	req := v.New(false)
	req.SetRequestURI(FullApiUrl + "/gateway")
	err = v.DoResult(req, &gateway)
	url = gateway.Url
	return
}

type SessionInfo struct {
	Url    string `json:"url"`
	Shards uint16 `json:"shards"`
	Limits struct {
		Total          int `json:"total"`
		Remaining      int `json:"remaining"`
		ResetAfter     int `json:"reset_after"`
		MaxConcurrency int `json:"max_concurrency"`
	} `json:"session_start_limit"`
}

func (v Client) SessionData() (data SessionInfo, err error) {
	req := v.New(true)
	req.SetRequestURI(FullApiUrl + "/gateway/bot")
	err = v.DoResult(req, &data, WithRetryDelay(10*time.Second))
	return
}

func NewClient(token string, options ...Option) *Client {
	def := &RequestData{
		timeout:          10 * time.Second,
		retryDelay:       5 * time.Second,
		retries:          3,
		prefix:           FullApiUrl,
		authHeaderPrefix: "Bot",
	}
	for i := range options {
		options[i](def)
	}
	if def.logger == nil {
		def.logger = golog.NewWithLevel("api", golog.LevelWarn)
	}
	var header []byte
	if token != "" {
		header = append([]byte(fmt.Sprintf("%v ", def.authHeaderPrefix)), token...)
	}
	return &Client{
		logger:       def.logger,
		token:        token,
		retryDelay:   def.retryDelay,
		timeout:      def.timeout,
		maxRetries:   def.retries,
		staticHeader: header,
		manager:      limits.NewManager(def.prefix),
	}
}
