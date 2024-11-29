package cache

import (
	. "github.com/andersfylling/snowflake/v5"

	. "github.com/BOOMfinity/bfcord/discord"
)

type Default struct {
	users           Map[ID, User]
	guilds          Map[ID, Guild]
	messages        SubMap[ID, Map[ID, Message]]
	channels        Map[ID, Channel]
	presences       SubMap[ID, Map[ID, Presence]]
	members         SubMap[ID, Map[ID, Member]]
	scheduledEvents SubMap[ID, Map[ID, ScheduledEvent]]
	voice           SubMap[ID, Map[ID, VoiceState]]
}

func (d *Default) Users() Map[ID, User] {
	return d.users
}

func (d *Default) Guilds() Map[ID, Guild] {
	return d.guilds
}

func (d *Default) Messages() SubMap[ID, Map[ID, Message]] {
	return d.messages
}

func (d *Default) Channels() Map[ID, Channel] {
	return d.channels
}

func (d *Default) Presences() SubMap[ID, Map[ID, Presence]] {
	return d.presences
}

func (d *Default) Members() SubMap[ID, Map[ID, Member]] {
	return d.members
}

func (d *Default) ScheduledEvents() SubMap[ID, Map[ID, ScheduledEvent]] {
	return d.scheduledEvents
}

func (d *Default) VoiceStates() SubMap[ID, Map[ID, VoiceState]] {
	return d.voice
}

func NewDefault(cfg *DefaultConfig) Store {
	if cfg == nil {
		cfg = &DefaultConfig{
			Users:           10_000,
			Guilds:          250,
			Channels:        100,
			MessageLimit:    100,
			PrivateChannels: 100,
			Presences:       100,
			Members:         100,
			Messages:        150,
		}
	}
	def := new(Default)

	def.users = NewMap[ID, User](cfg.Users)
	def.guilds = NewMap[ID, Guild](cfg.Guilds)
	def.channels = NewMap[ID, Channel](cfg.Channels)
	def.messages = NewSubMap[ID, Map[ID, Message]](func() Map[ID, Message] {
		return NewLimitedMap[ID, Message](cfg.MessageLimit, cfg.Messages)
	})
	def.members = NewSubMap[ID, Map[ID, Member]](func() Map[ID, Member] {
		return NewMap[ID, Member](cfg.Members)
	})
	def.scheduledEvents = NewSubMap[ID, Map[ID, ScheduledEvent]](func() Map[ID, ScheduledEvent] {
		return NewMap[ID, ScheduledEvent](0)
	})
	def.presences = NewSubMap[ID, Map[ID, Presence]](func() Map[ID, Presence] {
		return NewMap[ID, Presence](cfg.Presences)
	})
	def.voice = NewSubMap[ID, Map[ID, VoiceState]](func() Map[ID, VoiceState] {
		return NewMap[ID, VoiceState](0)
	})

	return def
}

// DefaultConfig is used to set up Default implementation of Store.
type DefaultConfig struct {
	Users           uint
	Guilds          uint
	Channels        uint
	Messages        uint
	PrivateChannels uint
	Presences       uint
	Members         uint
	MessageLimit    int
}
