package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/BOOMfinity/bfcord/slash"
	"github.com/BOOMfinity/go-utils/rate"
	"sync"
	"time"

	"github.com/BOOMfinity/bfcord/api"
	"github.com/BOOMfinity/bfcord/cache"
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/bfcord/gateway"
	"github.com/BOOMfinity/bfcord/gateway/events"
	"github.com/BOOMfinity/bfcord/internal/slices"
	"github.com/BOOMfinity/golog"
	"github.com/andersfylling/snowflake/v5"
	"go.uber.org/atomic"
)

type Client interface {
	discord.ClientQuery
	API() *api.Client
	Start(ctx context.Context) error
	AvgLatency() uint16
	Get(id uint16) *gateway.Shard
	All() []*gateway.Shard
	Sub() *EventManager
	Store() cache.Store
	Presence() gateway.PresenceSet
	Slash() slash.Query
	CalcShard(guild snowflake.ID) uint16
	FetchMembers(ctx context.Context, data gateway.RequestMembers) (members []discord.MemberWithUser, presences []discord.BasePresence, err error)
	ChangeVoiceState(ctx context.Context, opts gateway.ChangeVoiceStateOptions) (*gateway.VoiceStateUpdateEvent, *gateway.VoiceServerUpdateEvent, error)

	Wait()
}

type client struct {
	store cache.Store
	slash slash.Query
	*api.Client
	manager *EventManager
	limiter *rate.Limiter
	config  *options
	logger  golog.Logger
	token   string
	shards  []*gateway.Shard
	current snowflake.ID
	m       sync.RWMutex
}

func (v *client) Channel(id snowflake.ID) discord.ChannelQuery {
	resolver := &channelResolver{ChannelQuery: api.NewChannelQuery(v.API(), id), bot: v}
	resolver.resolverOptions = resolverOptions[discord.ChannelQuery]{data: resolver}
	return resolver
}

func (v *client) Guild(id snowflake.ID) discord.GuildQuery {
	resolver := &guildResolver{GuildQuery: api.NewGuildQuery(v.API(), id), bot: v}
	resolver.resolverOptions = resolverOptions[discord.GuildQuery]{data: resolver}
	return resolver
}

func (v *client) FetchMembers(ctx context.Context, data gateway.RequestMembers) (members []discord.MemberWithUser, presences []discord.BasePresence, err error) {
	shardID := v.CalcShard(data.GuildID)
	shard := v.Get(shardID)
	if shard == nil {
		return nil, nil, errors.New("invalid guild id")
	}
	members, presences, err = shard.Gateway().RequestMembers(ctx, data)
	if err != nil {
		return
	}
	if v.Store() == nil {
		return
	}
	for i := range members {
		v.Store().Users().Set(members[i].UserID, members[i].User)
		v.Store().Members().UnsafeGet(members[i].GuildID).Set(members[i].UserID, members[i].Member)
	}
	for i := range presences {
		v.Store().Presences().UnsafeGet(presences[i].GuildID).Set(presences[i].UserID, presences[i])
	}
	return
}

// ChangeVoiceState can be used to join or leave voice channel.
//
// See: https://discord.com/developers/docs/topics/voice-connections#establishing-a-voice-websocket-connection
//
// To leave voice channel, set opts.ChannelID to 0 - method will return nil pointers.
func (v *client) ChangeVoiceState(ctx context.Context, opts gateway.ChangeVoiceStateOptions) (state *gateway.VoiceStateUpdateEvent, server *gateway.VoiceServerUpdateEvent, err error) {
	shard := v.Get(v.CalcShard(opts.GuildID))
	if opts.GuildID.Valid() && opts.ChannelID.IsZero() { // short path for leaving
		return nil, nil, shard.Gateway().ChangeVoiceState(opts)
	}

	var stateEvtHandler, serverEvtHandler *EventHandler
	stateChan, serverChan := make(chan *gateway.VoiceStateUpdateEvent), make(chan *gateway.VoiceServerUpdateEvent)

	stateEvtHandler = v.manager.VoiceStateUpdate(func(_ Client, _ *gateway.Shard, ev gateway.VoiceStateUpdateEvent) {
		if ev.GuildID == opts.GuildID && ev.ChannelID == opts.ChannelID && ev.UserID == v.current {
			stateEvtHandler.Close()
			stateChan <- &ev
		}
	})
	serverEvtHandler = v.manager.VoiceServerUpdate(func(_ Client, _ *gateway.Shard, ev gateway.VoiceServerUpdateEvent) {
		if ev.GuildID == opts.GuildID {
			serverEvtHandler.Close()
			serverChan <- &ev
		}
	})

	_ = shard.Gateway().ChangeVoiceState(opts)

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer func() {
		cancel()
		close(stateChan)
		close(serverChan)
		stateEvtHandler.Close()
		serverEvtHandler.Close()
	}()

f:
	for server == nil || state == nil {
		select {
		case state = <-stateChan:
			if server != nil {
				break f
			}
		case server = <-serverChan:
			if state != nil {
				break f
			}
		case <-ctx.Done():
			err = ctx.Err()
			break f
		}
	}

	return
}

func (v *client) CalcShard(guild snowflake.ID) uint16 {
	return uint16(uint64(guild>>22) % uint64(v.config.ShardCount))
}

func (v *client) User(id snowflake.ID) discord.UserQuery {
	resolver := &userResolver{UserQuery: api.NewUserQuery(v.API(), id), bot: v}
	resolver.resolverOptions = resolverOptions[discord.UserQuery]{data: resolver}
	return resolver
}

func (v *client) API() *api.Client {
	return v.Client
}

func (v *client) Slash() slash.Query {
	return v.slash
}

func (v *client) Presence() gateway.PresenceSet {
	return &presenceSet{c: v}
}

func (v *client) All() (x []*gateway.Shard) {
	v.m.RLock()
	defer v.m.RUnlock()
	x = append(x, v.shards...)
	return
}

func (v *client) Store() cache.Store {
	return v.store
}

func (v *client) Sub() *EventManager {
	return v.manager
}

func (v *client) AvgLatency() uint16 {
	v.m.RLock()
	all := make([]uint16, 0, len(v.shards))
	for i := range v.shards {
		all = append(all, v.shards[i].Latency())
	}
	v.m.RUnlock()
	return uint16(slices.Sum(all) / len(all))
}

func (v *client) Wait() {
	for {
		time.Sleep(24 * time.Hour)
	}
}

func (v *client) Get(id uint16) *gateway.Shard {
	v.m.RLock()
	defer v.m.RUnlock()
	for _, shard := range v.shards {
		if shard.ID() == id {
			return shard
		}
	}
	return nil
}

func (v *client) CurrentUser() (user *discord.User, err error) {
	if v.current.Valid() && v.Store() != nil {
		user, found := v.Store().Users().Get(v.current)
		if found {
			return &user, nil
		}
	}
	return v.User(v.current).NoCache().Get()
}

func (v *client) Start(ctx context.Context) error {
	v.Log().Info().Send("Spawning %v shard(-s)", len(v.config.Shards))
	_ = v.limiter.Wait(context.Background())
	for i := range v.config.Shards {
		id := v.config.Shards[i]
		gtw := gateway.New(v.token, id, v.config.ShardCount, v.config.GatewayOptions...)
		shard := gateway.NewShard(gtw)
		v.shards = append(v.shards, shard)
		gtw.OnData(v.handle)
		member := gtw.EventChannel().Join()
		err := gtw.Connect(ctx)
		if err != nil {
			panic(err)
		}
		timeout := time.NewTimer(v.config.Timeout)
	waiting:
		for {
			select {
			case msg, more := <-member.Out:
				if !more {
					return errors.New("channel closed")
				}
				if msg.Data() == events.ShardReady && !v.config.BlockUntilPrefetch {
					break waiting
				}
				if msg.Data() == events.ShardPrefetchCompleted && v.config.BlockUntilPrefetch {
					break waiting
				}
			case <-timeout.C:
				return errors.New("timed out")
			}
		}
		timeout.Stop()
		member.Close()
		v.Log().Info().Send("Shard #%v has been successfully connected", id)
		_ = v.limiter.Wait(context.Background())
	}
	Execute(v.manager, func(_h ReadyEvent) {
		_h(v)
	})
	return nil
}

func (v *client) Log() golog.Logger {
	return v.logger
}

// New creates a client with default settings	 (automatic sharding, default cache). To override these settings, use Options
func New(token string, opt ...Option) (Client, error) {
	def := &options{AutoSharding: true, Logger: golog.New("bfcord"), BlockUntilPrefetch: true, Timeout: 45 * time.Second, Store: cache.NewDefaultStore()}
	gtwLog := def.Logger.Module("gateway")
	def.GatewayOptions = append(def.GatewayOptions, gateway.WithLogger(gtwLog), gateway.WithApiClient(api.NewClient(token, api.WithLogger(gtwLog.Module("api")))))
	for i := range opt {
		opt[i](def)
	}
	c := new(client)
	c.token = token
	c.logger = def.Logger
	c.Client = api.NewClient(token, def.APIOptions...)
	c.config = def
	data, err := c.Client.SessionData()
	if err != nil {
		return nil, err
	}
	if def.AutoSharding {
		def.ShardCount = data.Shards
	}
	c.limiter = rate.NewLimiter(5050*time.Millisecond, data.Limits.MaxConcurrency)
	if len(def.Shards) == 0 {
		for len(def.Shards) < int(def.ShardCount) {
			def.Shards = append(def.Shards, uint16(len(def.Shards)))
		}
	}
	c.shards = make([]*gateway.Shard, 0, len(def.Shards))
	c.manager = &EventManager{id: atomic.NewUint64(1)}
	user, err := c.API().CurrentUser()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch current user: %w", err)
	}
	c.slash = slash.NewClientWithAppID(c.token, user.ID)
	if def.Store != nil {
		c.store = def.Store
	}
	return c, nil
}

type presenceSet struct {
	c Client
}

func (v *presenceSet) SetCustom(data gateway.PresenceUpdate) {
	for _, shard := range v.c.All() {
		shard.Gateway().Presence().SetCustom(data)
	}
}

func (v *presenceSet) Set(status discord.PresenceStatus, ac discord.Activity) {
	for _, shard := range v.c.All() {
		shard.Gateway().Presence().Set(status, ac)
	}
}
