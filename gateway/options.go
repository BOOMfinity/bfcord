package gateway

import (
	"github.com/BOOMfinity/bfcord/api"
	"github.com/BOOMfinity/bfcord/gateway/intents"
	"github.com/BOOMfinity/golog"
)

type Options struct {
	logger     golog.Logger
	apiClient  *api.Client
	identify   Identify
	bufferSize int
}

type Option func(v *Options)

func WithIntents(intents intents.Intent) Option {
	return func(v *Options) {
		v.identify.Intents = intents
	}
}

func WithProperties(props IdentifyProperties) Option {
	return func(v *Options) {
		v.identify.Properties = props
	}
}

func WithLogger(log golog.Logger) Option {
	return func(v *Options) {
		v.logger = log
	}
}

func WithBufferSize(size int) Option {
	return func(v *Options) {
		v.bufferSize = size
	}
}

func WithApiClient(client *api.Client) Option {
	return func(v *Options) {
		v.apiClient = client
	}
}
