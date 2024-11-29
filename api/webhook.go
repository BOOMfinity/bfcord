package api

import (
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/bfcord/internal/httpc"
	"github.com/BOOMfinity/golog/v2"
	"github.com/andersfylling/snowflake/v5"
	"github.com/valyala/fasthttp"
	"net/url"
)

var _globalHTTP = httpc.NewClient("", golog.New("webhooks"))

type WebhookClient discord.Webhook

func (w *WebhookClient) Fetch() (err error) {
	id := w.ID
	token := w.Token
	*w, err = httpc.NewJSONRequest[WebhookClient](_globalHTTP, func(b httpc.RequestBuilder) error {
		return b.Execute("webhooks", w.ID.String(), w.Token)
	})
	w.ID = id
	w.Token = token
	return
}

func (w *WebhookClient) Delete() error {
	return httpc.NewRequest(_globalHTTP, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodDelete)
		return b.Execute("webhooks", w.ID.String(), w.Token)
	})
}

func (w *WebhookClient) Modify(params ModifyWebhookParams, reason ...string) (err error) {
	id := w.ID
	token := w.Token
	*w, err = httpc.NewJSONRequest[WebhookClient](_globalHTTP, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPatch)
		b.Body(params)
		b.Reason(reason...)
		return b.Execute("webhooks", w.ID.String(), w.Token)
	})
	w.ID = id
	w.Token = token
	return
}

func (w *WebhookClient) Execute(params WebhookExecuteParams) (msg discord.Message, err error) {
	if params.Wait {
		msg, err = httpc.NewJSONRequest[discord.Message](_globalHTTP, func(b httpc.RequestBuilder) error {
			b.Method(fasthttp.MethodPost)
			uploadFiles(b, params, params.Attachments)
			return b.Execute("webhooks", w.ID.String(), w.Token)
		})
	} else {
		err = httpc.NewRequest(_globalHTTP, func(b httpc.RequestBuilder) error {
			b.Method(fasthttp.MethodPost)
			uploadFiles(b, params, params.Attachments)
			return b.Execute("webhooks", w.ID.String(), w.Token)
		})
	}
	return
}

func (w *WebhookClient) GetMessage(id snowflake.ID, thread ...snowflake.ID) (discord.Message, error) {
	return httpc.NewJSONRequest[discord.Message](_globalHTTP, func(b httpc.RequestBuilder) error {
		values := url.Values{}
		if len(thread) > 0 {
			values.Set("thread_id", thread[0].String())
		}
		return b.Execute("webhooks", w.ID.String(), w.Token, "messages", id.String()+"?"+values.Encode())
	})
}

func (w *WebhookClient) DeleteMessage(id snowflake.ID) error {
	return httpc.NewRequest(_globalHTTP, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodDelete)
		return b.Execute("webhooks", w.ID.String(), w.Token, "messages", id.String())
	})
}

func (w *WebhookClient) EditMessage(id snowflake.ID, params EditMessageParams, thread ...snowflake.ID) (discord.Message, error) {
	return httpc.NewJSONRequest[discord.Message](_globalHTTP, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPatch)
		values := url.Values{}
		if len(thread) > 0 {
			values.Set("thread_id", thread[0].String())
		}
		uploadFiles(b, params, params.Attachments)
		return b.Execute("webhooks", w.ID.String(), w.Token, "messages", id.String()+"?"+values.Encode())
	})
}

func NewWebhook(id snowflake.ID, token string) WebhookClient {
	return WebhookClient{
		ID:    id,
		Token: token,
	}
}
