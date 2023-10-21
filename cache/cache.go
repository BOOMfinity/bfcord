package cache

import (
	"github.com/BOOMfinity/go-utils/sets"
	"sync"

	"github.com/BOOMfinity/bfcord/discord"
	"github.com/andersfylling/snowflake/v5"
)

var _ = (Store)(&DefaultStore{})

type Store interface {
	Guilds() SafeStoreCustom[snowflake.ID, discord.Guild]
	Members() SafeStoreCustom[snowflake.ID, SafeStoreCustom[snowflake.ID, discord.Member]]
	Channels() SafeStoreCustom[snowflake.ID, SafeStoreCustom[snowflake.ID, discord.Channel]]
	Reactions() SafeStoreCustom[snowflake.ID, SafeStoreCustom[string, discord.MessageReaction]]
	Presences() SafeStoreCustom[snowflake.ID, SafeStoreCustom[snowflake.ID, discord.BasePresence]]
	Messages() SafeStoreCustom[snowflake.ID, sets.Set[snowflake.ID, discord.BaseMessage]]
	Private() SafeStoreCustom[snowflake.ID, discord.Channel]
	Users() SafeStoreCustom[snowflake.ID, discord.User]
	VoiceStates() SafeStoreCustom[snowflake.ID, SafeStoreCustom[snowflake.ID, discord.VoiceState]]
	SetChannelGuild(channel, guild snowflake.ID)
	ChannelGuild(channel snowflake.ID) snowflake.ID
}

type SafeStoreCustom[K comparable, V any] interface {
	Get(key K) (val V, found bool)
	UnsafeGet(key K) V
	GetOrSet(key K, set func() V) V
	ToSlice() []V
	Filter(fn func(item V) bool) []V
	Each(fn func(item V))
	Find(fn func(item V) bool) (V, bool)
	Sort(fn func(a, b V) bool) []V
	Has(key K) bool
	Set(key K, value V)
	Delete(key K) bool
	Update(key K, fn func(value V) V) (ok bool)
	Size() int
}

func newSnowflakeStore[V any]() SafeStoreCustom[snowflake.ID, V] {
	return newSafeStoreCustom[snowflake.ID, V]()
}

func newSafeStoreCustom[K comparable, V any]() SafeStoreCustom[K, V] {
	return NewSafeMap[K, V](0)
}

func NewDefaultStore() *DefaultStore {
	store := new(DefaultStore)
	store.users = NewSafeMap[snowflake.ID, discord.User](0)
	store.guilds = NewSafeMap[snowflake.ID, discord.Guild](0)
	store.private = NewSafeMap[snowflake.ID, discord.Channel](0)
	store.reactions = NewSafeMapWithInitializer[snowflake.ID, SafeStoreCustom[string, discord.MessageReaction]](0, newSafeStoreCustom[string, discord.MessageReaction])
	store.members = NewSafeMapWithInitializer[snowflake.ID, SafeStoreCustom[snowflake.ID, discord.Member]](0, newSnowflakeStore[discord.Member])
	store.channels = NewSafeMapWithInitializer[snowflake.ID, SafeStoreCustom[snowflake.ID, discord.Channel]](0, newSnowflakeStore[discord.Channel])
	store.presences = NewSafeMapWithInitializer[snowflake.ID, SafeStoreCustom[snowflake.ID, discord.BasePresence]](0, newSnowflakeStore[discord.BasePresence])
	store.messages = NewSafeMapWithInitializer[snowflake.ID, sets.Set[snowflake.ID, discord.BaseMessage]](0, func() sets.Set[snowflake.ID, discord.BaseMessage] {
		return sets.NewLimitedCustomSet[snowflake.ID, discord.BaseMessage](func(item discord.BaseMessage) snowflake.ID {
			return item.ID
		}, 100)
	})
	store.aliases = map[snowflake.Snowflake]snowflake.Snowflake{}
	store.voiceStates = NewSafeMapWithInitializer[snowflake.ID, SafeStoreCustom[snowflake.ID, discord.VoiceState]](0, newSnowflakeStore[discord.VoiceState])

	return store
}

type DefaultStore struct {
	users     SafeStoreCustom[snowflake.ID, discord.User]
	guilds    SafeStoreCustom[snowflake.ID, discord.Guild]
	private   SafeStoreCustom[snowflake.ID, discord.Channel]
	reactions SafeStoreCustom[snowflake.ID, SafeStoreCustom[string, discord.MessageReaction]]
	// roles       SafeStore[SafeStore[discord.Role]]
	members     SafeStoreCustom[snowflake.ID, SafeStoreCustom[snowflake.ID, discord.Member]]
	channels    SafeStoreCustom[snowflake.ID, SafeStoreCustom[snowflake.ID, discord.Channel]]
	presences   SafeStoreCustom[snowflake.ID, SafeStoreCustom[snowflake.ID, discord.BasePresence]]
	voiceStates SafeStoreCustom[snowflake.ID, SafeStoreCustom[snowflake.ID, discord.VoiceState]]
	messages    SafeStoreCustom[snowflake.ID, sets.Set[snowflake.ID, discord.BaseMessage]]
	// emojis      SafeStore[SafeStore[discord.Emoji]]
	aliases map[snowflake.ID]snowflake.ID
	m       sync.RWMutex
}

func (d *DefaultStore) Messages() SafeStoreCustom[snowflake.ID, sets.Set[snowflake.ID, discord.BaseMessage]] {
	return d.messages
}

func (d *DefaultStore) VoiceStates() SafeStoreCustom[snowflake.ID, SafeStoreCustom[snowflake.ID, discord.VoiceState]] {
	return d.voiceStates
}

func (d *DefaultStore) Reactions() SafeStoreCustom[snowflake.ID, SafeStoreCustom[string, discord.MessageReaction]] {
	return d.reactions
}

func (d *DefaultStore) Presences() SafeStoreCustom[snowflake.ID, SafeStoreCustom[snowflake.ID, discord.BasePresence]] {
	return d.presences
}

func (d *DefaultStore) Private() SafeStoreCustom[snowflake.ID, discord.Channel] {
	return d.private
}

func (d *DefaultStore) Guilds() SafeStoreCustom[snowflake.ID, discord.Guild] {
	return d.guilds
}

func (d *DefaultStore) Members() SafeStoreCustom[snowflake.ID, SafeStoreCustom[snowflake.ID, discord.Member]] {
	return d.members
}

func (d *DefaultStore) Channels() SafeStoreCustom[snowflake.ID, SafeStoreCustom[snowflake.ID, discord.Channel]] {
	return d.channels
}

func (d *DefaultStore) Users() SafeStoreCustom[snowflake.ID, discord.User] {
	return d.users
}

func (d *DefaultStore) SetChannelGuild(channel, guild snowflake.ID) {
	d.m.Lock()
	defer d.m.Unlock()
	d.aliases[channel] = guild
}

func (d *DefaultStore) ChannelGuild(channel snowflake.ID) snowflake.ID {
	d.m.RLock()
	defer d.m.RUnlock()
	if id, ok := d.aliases[channel]; ok {
		return id
	}
	return 0
}
