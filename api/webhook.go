package api

import (
	"fmt"
	"github.com/BOOMfinity/bfcord/api/builders"
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/andersfylling/snowflake/v5"
	"github.com/valyala/fasthttp"
)

type webhookQuery struct {
	api   *Client
	token string
	id    snowflake.ID
}

func (w webhookQuery) url() string {
	if w.token == "" {
		return fmt.Sprintf("%v/webhooks/%v", FullApiUrl, w.id)
	} else {
		return fmt.Sprintf("%v/webhooks/%v/%v", FullApiUrl, w.id, w.token)
	}
}

func (w webhookQuery) Fetch() (wh discord.Webhook, err error) {
	req := w.api.New(true)
	req.SetRequestURI(w.url())
	err = w.api.DoResult(req, &wh)
	return
}

func (w webhookQuery) Execute() discord.WebhookExecuteBuilder {
	return builders.NewWebhookExecuteBuilder(w.id, w.token, w.api)
}

func (w webhookQuery) Delete() (err error) {
	req := w.api.New(true)
	req.SetRequestURI(w.url())
	req.Header.SetMethod(fasthttp.MethodDelete)
	return w.api.DoNoResp(req)
}

func (w webhookQuery) DeleteMessage(id snowflake.ID) (err error) {
	req := w.api.New(true)
	req.SetRequestURI(fmt.Sprintf("%v/webhooks/%v/%v/messages/%v", FullApiUrl, w.id, w.token, id))
	req.Header.SetMethod(fasthttp.MethodDelete)
	return w.api.DoNoResp(req)
}

func (w webhookQuery) EditMessage(id snowflake.ID) discord.WebhookUpdateMessageBuilder {
	return builders.NewWebhookUpdateMessageBuilder(w.id, w.token, id, w.api)
}
