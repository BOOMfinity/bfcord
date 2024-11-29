package httpc

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/BOOMfinity/golog/v2"
	"github.com/segmentio/encoding/json"
	"github.com/valyala/fasthttp"
)

type RequestOption func(req *fasthttp.Request) error
type ResponseOption func(res *fasthttp.Response) error

type RequestBuilder interface {
	Method(m string) RequestBuilder
	Header(name, value string) RequestBuilder
	Body(v any) RequestBuilder
	Reason(...string) RequestBuilder
	Parse(v any) RequestBuilder
	Retries(n uint) RequestBuilder
	NoAuth() RequestBuilder
	Debug() RequestBuilder
	Multipart(fn func(writer *multipart.Writer) error) RequestBuilder
	OnRequest(opts ...RequestOption) RequestBuilder
	OnResponse(opts ...ResponseOption) RequestBuilder
	Execute(segments ...string) error
}

type requestBuilderImpl struct {
	log       golog.Logger
	client    *Client
	limit     uint
	reqOpts   []RequestOption
	resOpts   []ResponseOption
	doNotAuth bool
}

func (b *requestBuilderImpl) OnRequest(opts ...RequestOption) RequestBuilder {
	b.reqOpts = append(b.reqOpts, opts...)
	return b
}

func (b *requestBuilderImpl) OnResponse(opts ...ResponseOption) RequestBuilder {
	b.resOpts = append(b.resOpts, opts...)
	return b
}

func (b *requestBuilderImpl) Reason(v ...string) RequestBuilder {
	if len(v) > 0 {
		b.reqOpts = append(b.reqOpts, func(req *fasthttp.Request) error {
			req.Header.Set("X-Audit-Log-Reason", strings.Join(v, " "))
			return nil
		})
	}
	return b
}

func (b *requestBuilderImpl) Multipart(fn func(writer *multipart.Writer) error) RequestBuilder {
	b.reqOpts = append(b.reqOpts, func(req *fasthttp.Request) error {
		writer := multipart.NewWriter(req.BodyWriter())
		req.Header.SetContentType(writer.FormDataContentType())
		req.Header.SetMultipartFormBoundary(writer.Boundary())
		if err := fn(writer); err != nil {
			return fmt.Errorf("failed to execute user multipart function: %w", err)
		}
		if err := writer.Close(); err != nil {
			return fmt.Errorf("failed to close multipart writer: %w", err)
		}
		return nil
	})
	return b
}

func (b *requestBuilderImpl) Header(name, value string) RequestBuilder {
	b.reqOpts = append(b.reqOpts, func(req *fasthttp.Request) error {
		req.Header.Add(name, value)
		return nil
	})
	return b
}

func (b *requestBuilderImpl) Debug() RequestBuilder {
	b.log.Level(golog.LevelTrace)
	return b
}

func (b *requestBuilderImpl) Retries(n uint) RequestBuilder {
	b.limit = n
	return b
}

func (b *requestBuilderImpl) Parse(v any) RequestBuilder {
	b.resOpts = append(b.resOpts, func(res *fasthttp.Response) error {
		if err := json.Unmarshal(res.Body(), v); err != nil {
			return fmt.Errorf("failed to unmarshal response struct: %w", err)
		}
		return nil
	})
	return b
}

func (b *requestBuilderImpl) Body(v any) RequestBuilder {
	b.reqOpts = append(b.reqOpts, func(req *fasthttp.Request) error {
		data, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		req.SetBody(data)
		return nil
	})
	return b
}

func (b *requestBuilderImpl) Method(m string) RequestBuilder {
	b.reqOpts = append(b.reqOpts, func(req *fasthttp.Request) error {
		req.Header.SetMethod(m)
		return nil
	})
	return b
}

func (b *requestBuilderImpl) NoAuth() RequestBuilder {
	b.doNotAuth = true
	return b
}

func (b *requestBuilderImpl) Execute(segments ...string) error {
	url := ResolvePath(segments...)
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)
	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.SetContentType("application/json")
	if !b.doNotAuth && b.client.token != "" {
		req.Header.Set("Authorization", "Bot "+b.client.token)
	}
	for _, opt := range b.reqOpts {
		if err := opt(req); err != nil {
			return errors.Join(ErrFailedToParseRequestOptions, err)
		}
	}
	for r := (uint)(0); r < b.limit; r++ {
		_ = b.client.limiter.Wait(context.Background())
		if err := b.client.buckets.Acquire(context.Background(), url); err != nil {
			return fmt.Errorf("failed to acquire the bucket: %w", err)
		}
		b.log.Trace().Param("retries", r).Send("Sending request to the '%s' with %d bytes of body", url, len(req.Body()))
		if err := fasthttp.Do(req, res); err != nil {
			return fmt.Errorf("failed to send the request: %w", err)
		}
		if res.StatusCode() == fasthttp.StatusTooManyRequests {
			if err := b.client.buckets.Release(url, &res.Header); err != nil {
				return fmt.Errorf("failed to release the bucket: %w", err)
			}
			continue
		}
		if res.StatusCode() >= 400 && res.StatusCode() < 500 {
			if err := b.client.buckets.Release(url, &res.Header); err != nil {
				return fmt.Errorf("failed to release the bucket: %w", err)
			}
			if json.Valid(res.Body()) {
				var dcErr DiscordError
				if err := json.Unmarshal(res.Body(), &dcErr); err != nil {
					return fmt.Errorf("failed to parse Discord error: %w", err)
				}
				if dcErr.Code == 0 {
					return fmt.Errorf("something went wrong: %s", string(res.Body()))
				}
				return &dcErr
			}
			return fmt.Errorf("something went wrong: %s", string(res.Body()))
		}
		if res.StatusCode() >= 500 {
			return fmt.Errorf("there is a server error on Discord side: %s", string(res.Body()))
		}
		if err := b.client.buckets.Release(url, &res.Header); err != nil {
			return fmt.Errorf("failed to release the bucket: %w", err)
		}
		for _, opt := range b.resOpts {
			if err := opt(res); err != nil {
				return errors.Join(ErrFailedToParseResponseOptions, err)
			}
		}
		return nil
	}
	return ErrMaxRetriesReached
}

func NewRequest(c *Client, fn func(b RequestBuilder) error) error {
	id := c.id.Add(1)
	return fn(&requestBuilderImpl{
		log:    c.log.Copy().Scope(fmt.Sprint(id)),
		client: c,
		limit:  5,
	})
}

func NewJSONRequest[T any](c *Client, fn func(b RequestBuilder) error) (dst T, err error) {
	return dst, NewRequest(c, func(b RequestBuilder) error {
		b.Parse(&dst)
		return fn(b)
	})
}
