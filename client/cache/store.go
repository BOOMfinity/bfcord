package cache

import (
	"errors"

	. "github.com/andersfylling/snowflake/v5"

	. "github.com/BOOMfinity/bfcord/discord"
)

var (
	ErrNotFound = errors.New("resource not found in cache")
)

type Store interface {
	Users() Map[ID, User]
	Guilds() Map[ID, Guild]
	Messages() SubMap[ID, Map[ID, Message]]
	Channels() Map[ID, Channel]
	Presences() SubMap[ID, Map[ID, Presence]]
	Members() SubMap[ID, Map[ID, Member]]
	ScheduledEvents() SubMap[ID, Map[ID, ScheduledEvent]]
	VoiceStates() SubMap[ID, Map[ID, VoiceState]]
}
