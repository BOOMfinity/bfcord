package client

import (
	"github.com/BOOMfinity/bfcord/api"
	"github.com/BOOMfinity/bfcord/cache"
	"github.com/BOOMfinity/bfcord/gateway"
	"github.com/BOOMfinity/bfcord/gateway/intents"
	"github.com/BOOMfinity/golog"
	"time"
)

type Options struct {
	Store              cache.Store
	Logger             golog.Logger
	Shards             []uint16
	IgnoredEvents      []string
	GatewayOptions     []gateway.Option
	APIOptions         []api.Option
	Timeout            time.Duration
	ShardCount         uint16
	AutoSharding       bool
	BlockUntilPrefetch bool
}

type Option func(v *Options)

func WithAPIOptions(opts ...api.Option) Option {
	return func(v *Options) {
		v.APIOptions = opts
	}
}

func WithConnectionTimeout(timeout time.Duration) Option {
	return func(v *Options) {
		v.Timeout = timeout
	}
}

func WithStore(store cache.Store) Option {
	return func(v *Options) {
		v.Store = store
	}
}

func WithDisabledPrefetchBlock() Option {
	return func(v *Options) {
		v.BlockUntilPrefetch = false
	}
}

func WithGatewayOpts(options ...gateway.Option) Option {
	return func(v *Options) {
		v.GatewayOptions = append(v.GatewayOptions, options...)
	}
}

func WithShardCount(count uint16) Option {
	return func(v *Options) {
		v.ShardCount = count
		v.AutoSharding = false
	}
}

func WithShards(shards []uint16) Option {
	return func(v *Options) {
		v.Shards = shards
		v.AutoSharding = false
	}
}

func WithIntents(intents intents.Intent) Option {
	return WithGatewayOpts(gateway.WithIntents(intents))
}

func WithDebug(enabled bool) Option {
	return func(v *Options) {
		v.Logger.SetLevel(golog.LevelDebug)
	}
}

func WithLogger(logger golog.Logger) Option {
	return func(v *Options) {
		v.Logger = logger
	}
}
