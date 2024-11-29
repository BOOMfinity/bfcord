package client

import (
	"context"
	"fmt"
	"github.com/BOOMfinity/bfcord/api"
	"sync"
	"sync/atomic"
	"time"

	"github.com/BOOMfinity/golog/v2"
	"github.com/andersfylling/snowflake/v5"

	"github.com/BOOMfinity/bfcord/client/cache"
	"github.com/BOOMfinity/bfcord/client/events"
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/bfcord/utils"
	"github.com/BOOMfinity/bfcord/ws"
)

type Session interface {
	api.Client

	API() api.Client
	Events() events.SessionDispatcher
	Cache() cache.Store
	Log() golog.Logger
	ShardID(guild snowflake.ID) uint16

	Get(id uint16) Shard
	Ping() int
	Shards() []uint16
	ShardCount() uint16
	// Alive checks if shard is connected to the Discord Gateway. It is good to check when you want to make sure the cache is up-to-date.
	//
	// If no shard id is provided, then all shards of current Session instance are checked. It is not recommended in most cases.
	//
	// If gateway reach the max reconnection limit, it will wait 5 minutes before next try. Of course, you won't receive any new events you have to handle but if you have the interval jobs (or dashboard), you should be sure the specific shard is connected.
	Alive(shard ...uint16) bool
	// Unavailable checks if guild is not available due to outage or has not been loaded from lazy GUILD_CREATE events yet.
	Unavailable(id snowflake.ID) bool
	UnavailableCount() int
	FetchMembers(ctx context.Context, params ws.RequestGuildMembersParams) ([]discord.MemberWithUser, []discord.Presence, error)

	PermissionsIn(guild, channel, member snowflake.ID) (discord.Permission, error)
	SortedMemberRoles(guild, member snowflake.ID) ([]discord.Role, error)

	Shutdown()
	Start()
}

type sessionImpl struct {
	api.Client

	events      events.SessionDispatcher
	shards      []Shard
	shardCount  uint16
	log         golog.Logger
	cache       cache.Store
	handlers    utils.SimpleMap[string, handleDispatchFn]
	mut         sync.RWMutex
	unavailable utils.SimpleMap[uint16, utils.SimpleMap[snowflake.ID, ws.UnavailableGuild]]

	metrics struct {
		events    atomic.Uint64
		totalTime atomic.Uint64
	}
}

func (s *sessionImpl) SortedMemberRoles(guild, member snowflake.ID) ([]discord.Role, error) {
	roles, err := s.Guild(guild).Roles()
	if err != nil {
		return nil, fmt.Errorf("cannot get guild roles: %w", err)
	}
	m, err := s.Guild(guild).Member(member).Get()
	if err != nil {
		return nil, fmt.Errorf("cannot get guild member: %w", err)
	}
	return SortedMemberRoles(roles, m.Roles), nil
}

func (s *sessionImpl) UnavailableCount() (c int) {
	s.unavailable.Each(func(_ uint16, v utils.SimpleMap[snowflake.ID, ws.UnavailableGuild]) {
		c += v.Size()
	})
	return
}

func (s *sessionImpl) User(id snowflake.ID) api.UserClient {
	return userClient{
		UserClient: s.API().User(id),
		id:         id,
		sess:       s,
	}
}

func (s *sessionImpl) Channel(id snowflake.ID) api.ChannelClient {
	return channelClient{
		ChannelClient: s.API().Channel(id),
		id:            id,
		sess:          s,
	}
}

func (s *sessionImpl) Guild(id snowflake.ID) api.GuildClient {
	return guildClient{
		GuildClient: s.API().Guild(id),
		id:          id,
		sess:        s,
	}
}

func (s *sessionImpl) API() api.Client {
	return s.Client
}

func (s *sessionImpl) FetchMembers(ctx context.Context, params ws.RequestGuildMembersParams) ([]discord.MemberWithUser, []discord.Presence, error) {
	id := s.ShardID(params.GuildID)
	shard := s.Get(id)
	members, presences, err := shard.FetchMembers(ctx, params)
	if s.Cache() != nil {
		for _, v := range members {
			s.Cache().Members().Get(params.GuildID).Set(v.User.ID, v.Member)
			s.Cache().Users().Set(v.User.ID, v.User)
		}
		for _, v := range presences {
			s.Cache().Presences().Get(params.GuildID).Set(v.User.ID, v)
		}
	}
	return members, presences, err
}

func (s *sessionImpl) Alive(shard ...uint16) bool {
	if len(shard) == 0 {
		shard = s.Shards()
	}
	for _, v := range shard {
		if s.Get(v).Status() != ws.StatusConnected {
			return false
		}
	}
	return true
}

func (s *sessionImpl) Unavailable(id snowflake.ID) bool {
	shard := s.ShardID(id)
	if data, ok := s.unavailable.Get(shard); ok {
		if _, ok = data.Get(id); ok {
			return true
		}
	}
	return false
}

func (s *sessionImpl) Events() events.SessionDispatcher {
	return s.events
}

func (s *sessionImpl) Cache() cache.Store {
	return s.cache
}

func (s *sessionImpl) Log() golog.Logger {
	return s.log
}

func (s *sessionImpl) ShardID(guild snowflake.ID) uint16 {
	return uint16(guild>>22) % s.shardCount
}

func (s *sessionImpl) ShardCount() uint16 {
	return s.shardCount
}

func (s *sessionImpl) Shards() (shards []uint16) {
	s.mut.RLock()
	defer s.mut.RUnlock()
	shards = make([]uint16, len(s.shards))
	for _, shard := range s.shards {
		shards = append(shards, shard.ID())
	}
	return
}

func (s *sessionImpl) Get(id uint16) Shard {
	s.mut.RLock()
	defer s.mut.RUnlock()
	for _, shard := range s.shards {
		if shard.ID() == id {
			return shard
		}
	}
	return nil
}

func (s *sessionImpl) Ping() (avg int) {
	s.mut.RLock()
	defer s.mut.RUnlock()
	for _, shard := range s.shards {
		avg += shard.Ping()
	}
	avg = avg / len(s.shards)
	return
}

func (s *sessionImpl) Shutdown() {
	s.mut.RLock()
	defer s.mut.RUnlock()
	for _, shard := range s.shards {
		shard.Disconnect()
	}
}

func (s *sessionImpl) registerEventHandlers() {
	s.handlers.Set("READY", readyEventHandler)
	s.handlers.Set("GUILD_CREATE", guildCreateEventHandler)
	s.handlers.Set("GUILD_UPDATE", guildUpdateEventHandler)
	s.handlers.Set("GUILD_DELETE", guildDeleteEventHandler)
	s.handlers.Set("GUILD_BAN_ADD", guildBan(true))
	s.handlers.Set("GUILD_BAN_REMOVE", guildBan(false))
	s.handlers.Set("CHANNEL_CREATE", channelCreateEventHandler)
	s.handlers.Set("CHANNEL_UPDATE", channelUpdateEventHandler)
	s.handlers.Set("CHANNEL_DELETE", channelDeleteEventHandler)
	s.handlers.Set("CHANNEL_PINS_UPDATE", channelPinsUpdateEventHandler)
	s.handlers.Set("MESSAGE_CREATE", messageCreateEventHandler)
	s.handlers.Set("MESSAGE_UPDATE", messageUpdateEventHandler)
	s.handlers.Set("MESSAGE_DELETE", messageDeleteEventHandler)
	s.handlers.Set("THREAD_CREATE", threadCreateEventHandler)
	s.handlers.Set("THREAD_UPDATE", threadUpdateEventHandler)
	s.handlers.Set("THREAD_DELETE", threadDeleteEventHandler)
	s.handlers.Set("THREAD_LIST_SYNC", threadListSyncEventHandler)
	s.handlers.Set("THREAD_MEMBERS_UPDATE", threadMembersUpdateEventHandler)
	s.handlers.Set("GUILD_ROLE_ADD", guildRoleEventHandler(false))
	s.handlers.Set("GUILD_ROLE_UPDATE", guildRoleEventHandler(true))
	s.handlers.Set("GUILD_ROLE_DELETE", guildRoleDeleteEventHandler)
	s.handlers.Set("GUILD_SCHEDULED_EVENT_CREATE", guildScheduledEventHandler("create"))
	s.handlers.Set("GUILD_SCHEDULED_EVENT_UPDATE", guildScheduledEventHandler("update"))
	s.handlers.Set("GUILD_SCHEDULED_EVENT_DELETE", guildScheduledEventHandler("delete"))
	s.handlers.Set("GUILD_SCHEDULED_EVENT_USER_ADD", guildScheduledUserEventHandler(false))
	s.handlers.Set("GUILD_SCHEDULED_EVENT_USER_REMOVE", guildScheduledUserEventHandler(true))
	s.handlers.Set("GUILD_MEMBER_ADD", guildMemberAddEventHandler)
	s.handlers.Set("GUILD_MEMBER_UPDATE", guildMemberUpdateEventHandler)
	s.handlers.Set("GUILD_MEMBER_REMOVE", guildMemberRemoveEventHandler)
	s.handlers.Set("INVITE_CREATE", inviteCreateEventHandler)
	s.handlers.Set("INVITE_DELETE", inviteDeleteEventHandler)
	s.handlers.Set("INTERACTION_CREATE", handle[discord.Interaction](func(log golog.Logger, sess Session, _ ws.InternalDispatchEvent, _ Shard, data *discord.Interaction) {
		sess.Events().InteractionCreate().Sender(func(handler events.InteractionCreateEvent) {
			handler(data, sess.Interaction(data.ID, data.Token))
		})
	}))
	s.handlers.Set("VOICE_STATE_UPDATE", handleVoiceStateUpdate)
	s.handlers.Set("VOICE_SERVER_UPDATE", handleVoiceServerUpdate)
}

func (s *sessionImpl) Start() {
	go s.metricsService()
	s.mut.RLock()
	defer s.mut.RUnlock()
	var wg sync.WaitGroup
	wg.Add(len(s.shards))
	for _, shard := range s.shards {
		go func() {
			go s.handleEvents(shard)
			defer wg.Done()
			if err := shard.Connect(context.Background()); err != nil {
				panic(fmt.Errorf("failed to start shard #%d: %w", shard.ID(), err))
			}
		}()
		time.Sleep(150 * time.Millisecond)
	}
	wg.Wait()
}

func New() Creator {
	return &creatorImpl{
		cache:      cache.NewDefault(nil),
		log:        golog.New("bfcord"),
		shardCount: 1,
		shards:     []uint16{0},
		intents:    ws.GatewayIntentDefault,
	}
}
