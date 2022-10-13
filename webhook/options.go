package webhook

import (
	"github.com/BOOMfinity/bfcord/api"
	"github.com/BOOMfinity/golog"
)

type Option func(opt *options)

type options struct {
	logger golog.Logger
	client *api.Client
}

func WithLogger(l golog.Logger) Option {
	return func(opt *options) {
		opt.logger = l
	}
}

func WithClient(client *api.Client) Option {
	return func(opt *options) {
		opt.client = client
	}
}
