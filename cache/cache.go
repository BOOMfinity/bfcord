package cache

import (
	"github.com/BOOMfinity/go-utils/sets"
	"sync"

	"github.com/BOOMfinity/bfcord/discord"
	"github.com/andersfylling/snowflake/v5"
)

var _ = (Store)(&DefaultStore{})

type Store interface {
	Guilds() SafeStore[discord.BaseGuild]
	Members() SafeStore[SafeStore[discord.Member]]
	Channels() SafeStore[SafeStore[discord.Channel]]
	Reactions() SafeStore[SafeStoreCustom[string, discord.MessageReaction]]
	Presences() SafeStore[SafeStore[discord.BasePresence]]
	Messages() SafeStore[sets.Set[snowflake.ID, discord.BaseMessage]]
	Private() SafeStore[discord.Channel]
	Users() SafeStore[discord.User]
	VoiceStates() SafeStore[SafeStore[discord.VoiceState]]
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
	Size() int
}

type SafeStore[V any] SafeStoreCustom[snowflake.ID, V]

func newSafeStore[V any]() SafeStore[V] {
	return newSafeStoreCustom[snowflake.ID, V]()
}

func newSafeStoreCustom[K comparable, V any]() SafeStoreCustom[K, V] {
	return NewSafeMap[K, V](0)
}

func NewDefaultStore() *DefaultStore {
	store := new(DefaultStore)
	store.users = NewSafeMap[snowflake.ID, discord.User](0)
	store.guilds = NewSafeMap[snowflake.ID, discord.BaseGuild](0)
	store.private = NewSafeMap[snowflake.ID, discord.Channel](0)
	store.reactions = NewSafeMapWithInitializer[snowflake.ID, SafeStoreCustom[string, discord.MessageReaction]](0, newSafeStoreCustom[string, discord.MessageReaction])
	store.members = NewSafeMapWithInitializer[snowflake.ID, SafeStore[discord.Member]](0, newSafeStore[discord.Member])
	store.channels = NewSafeMapWithInitializer[snowflake.ID, SafeStore[discord.Channel]](0, newSafeStore[discord.Channel])
	store.presences = NewSafeMapWithInitializer[snowflake.ID, SafeStore[discord.BasePresence]](0, newSafeStore[discord.BasePresence])
	store.messages = NewSafeMapWithInitializer[snowflake.ID, sets.Set[snowflake.ID, discord.BaseMessage]](0, func() sets.Set[snowflake.ID, discord.BaseMessage] {
		return sets.NewLimitedCustomSet[snowflake.ID, discord.BaseMessage](func(item discord.BaseMessage) snowflake.ID {
			return item.ID
		}, 100)
	})
	store.aliases = map[snowflake.Snowflake]snowflake.Snowflake{}
	store.voiceStates = NewSafeMapWithInitializer[snowflake.ID, SafeStore[discord.VoiceState]](0, newSafeStore[discord.VoiceState])

	return store
}

type DefaultStore struct {
	users       SafeStore[discord.User]
	guilds      SafeStore[discord.BaseGuild]
	private     SafeStore[discord.Channel]
	reactions   SafeStore[SafeStoreCustom[string, discord.MessageReaction]]
	roles       SafeStore[SafeStore[discord.Role]]
	members     SafeStore[SafeStore[discord.Member]]
	channels    SafeStore[SafeStore[discord.Channel]]
	presences   SafeStore[SafeStore[discord.BasePresence]]
	voiceStates SafeStore[SafeStore[discord.VoiceState]]
	messages    SafeStore[sets.Set[snowflake.ID, discord.BaseMessage]]
	emojis      SafeStore[SafeStore[discord.Emoji]]
	aliases     map[snowflake.ID]snowflake.ID
	m           sync.RWMutex
}

func (d *DefaultStore) Messages() SafeStore[sets.Set[snowflake.ID, discord.BaseMessage]] {
	return d.messages
}

func (d *DefaultStore) VoiceStates() SafeStore[SafeStore[discord.VoiceState]] {
	return d.voiceStates
}

func (d *DefaultStore) Reactions() SafeStore[SafeStoreCustom[string, discord.MessageReaction]] {
	return d.reactions
}

func (d *DefaultStore) Presences() SafeStore[SafeStore[discord.BasePresence]] {
	return d.presences
}

func (d *DefaultStore) Private() SafeStore[discord.Channel] {
	return d.private
}

func (d *DefaultStore) Guilds() SafeStore[discord.BaseGuild] {
	return d.guilds
}

func (d *DefaultStore) Members() SafeStore[SafeStore[discord.Member]] {
	return d.members
}

func (d *DefaultStore) Channels() SafeStore[SafeStore[discord.Channel]] {
	return d.channels
}

func (d *DefaultStore) Users() SafeStore[discord.User] {
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
