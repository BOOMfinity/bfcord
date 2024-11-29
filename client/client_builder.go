package client

import (
	"fmt"
	"github.com/BOOMfinity/bfcord/api"
	"time"

	"github.com/BOOMfinity/go-utils/rate"
	"github.com/BOOMfinity/golog/v2"
	"github.com/andersfylling/snowflake/v5"

	"github.com/BOOMfinity/bfcord/client/cache"
	"github.com/BOOMfinity/bfcord/client/events"
	"github.com/BOOMfinity/bfcord/utils"
	"github.com/BOOMfinity/bfcord/ws"
)

type Creator interface {
	Cache(c cache.Store) Creator
	Logger(log golog.Logger) Creator
	EnableAutoSharding() Creator
	ShardCount(c uint16) Creator
	Shards(c uint16, s ...uint16) Creator
	Concurrency(c int) Creator
	Intents(i ws.GatewayIntent) Creator
	Build(token string) (Session, error)
}

type creatorImpl struct {
	cache        cache.Store
	log          golog.Logger
	autoSharding bool
	shardCount   uint16
	shards       []uint16
	intents      ws.GatewayIntent
	concurrency  int
}

func (ctr *creatorImpl) Concurrency(c int) Creator {
	ctr.concurrency = c
	return ctr
}

func (ctr *creatorImpl) Intents(i ws.GatewayIntent) Creator {
	ctr.intents = i
	return ctr
}

func (ctr *creatorImpl) Cache(c cache.Store) Creator {
	ctr.cache = c
	return ctr
}

func (ctr *creatorImpl) Logger(log golog.Logger) Creator {
	ctr.log = log
	return ctr
}

func (ctr *creatorImpl) EnableAutoSharding() Creator {
	ctr.autoSharding = true
	return ctr
}

func (ctr *creatorImpl) ShardCount(c uint16) Creator {
	ctr.shardCount = c
	ctr.shards = make([]uint16, 0, ctr.shardCount)
	for i := range c {
		ctr.shards = append(ctr.shards, i)
	}
	return ctr
}

func (ctr *creatorImpl) Shards(c uint16, s ...uint16) Creator {
	ctr.shards = s
	ctr.shardCount = c
	return ctr
}

func (ctr *creatorImpl) Build(token string) (Session, error) {
	if token == "" {
		return nil, fmt.Errorf("token required")
	}
	rest := api.NewClient(ctr.log.Module("api"), token, api.WithCacheProxy(proxyImpl{ctr.cache}))
	sess := new(sessionImpl)
	sess.handlers = utils.NewSimpleMap[string, handleDispatchFn]()
	sess.unavailable = utils.NewSimpleMap[uint16, utils.SimpleMap[snowflake.ID, ws.UnavailableGuild]]()
	sess.events = events.NewSessionDispatcher(ctr.log.Module("dispatcher"))
	sess.log = ctr.log
	sess.cache = ctr.cache
	sess.Client = rest
	{
		ctr.log.Debug().Send("Fetching current user")
		user, err := rest.GetCurrentUser()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch current user: %w", err)
		}
		ctr.log.Debug().Send("Identified as %s (%d)", user.Username, user.ID)
	}
	ctr.log.Debug().Send("Fetching bot gateway info")
	gateway, err := rest.GatewayInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bot gateway: %w", err)
	}
	if ctr.autoSharding {
		ctr.log.Debug().Send("Using recommended shard count from bot gateway info (%d)", gateway.Shards)
		ctr.shardCount = gateway.Shards
		ctr.shards = make([]uint16, 0, ctr.shardCount)
		for i := range ctr.shardCount {
			ctr.shards = append(ctr.shards, i)
		}
	} else {
		ctr.log.Debug().Send("Using used-defined sharding parameters (%d total shards, %v shard ids to connect)", ctr.shardCount, ctr.shards)
	}
	if ctr.concurrency == 0 {
		ctr.log.Debug().Send("Concurrency limit: %d (from API)", gateway.Limit.MaxConcurrency)
		ctr.concurrency = gateway.Limit.MaxConcurrency
	} else {
		ctr.log.Debug().Send("Concurrency limit: %d (used-defined)", ctr.concurrency)
	}
	ctr.log.Debug().Send("Intents: %d", ctr.intents)

	sess.shards = make([]Shard, len(ctr.shards))
	sess.shardCount = ctr.shardCount

	limiter := rate.NewLimiter(5250*time.Millisecond, ctr.concurrency)

	for i, id := range ctr.shards {
		shard := &shardImpl{
			ping:        -1,
			unavailable: utils.NewSimpleMap[snowflake.ID, ws.UnavailableGuild](),
			Gateway: ws.NewGateway(ws.Config{
				Logger:        ctr.log.Module("gateway"),
				Intents:       ctr.intents,
				ID:            id,
				URL:           gateway.URL,
				Compression:   false,
				Token:         token,
				ShardCount:    ctr.shardCount,
				GlobalLimiter: limiter,
			}),
		}

		go shard.backgroundJob()

		sess.shards[i] = shard
	}

	ctr.log.Debug().Send("Registering event handlers")
	sess.registerEventHandlers()

	ctr.log.Debug().Send("Configured session to run %d of %d shard(s), with concurrency set to %d and shard ids: %s", len(ctr.shards), ctr.shardCount, ctr.concurrency, join(ctr.shards, ","))

	return sess, nil
}

func join[T any](arr []T, s string) (str string) {
	for i := range arr {
		if i == 0 {
			str = fmt.Sprint(arr[0])
			continue
		}
		str += s + " " + fmt.Sprint(arr[i])
	}
	return
}
