package webhook

import (
	"github.com/BOOMfinity/bfcord/api"
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/andersfylling/snowflake/v5"
)

type Webhook struct {
	api *api.Client
	discord.WebhookQuery
}

func (w *Webhook) Reset(id snowflake.ID, token string) *Webhook {
	w.WebhookQuery = w.api.Webhook(id, token)
	return w
}

func New(id snowflake.ID, token string, opts ...Option) *Webhook {
	opt := &options{}
	for i := range opts {
		opts[i](opt)
	}
	w := new(Webhook)
	if opt.client == nil {
		opt.client = api.NewClient("", api.WithLogger(opt.logger))
	}
	w.api = opt.client
	w.Reset(id, token)
	return w
}
