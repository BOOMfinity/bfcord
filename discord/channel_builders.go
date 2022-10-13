package discord

import (
	"github.com/BOOMfinity/bfcord/discord/permissions"
	"github.com/andersfylling/snowflake/v5"
)

type CreateThreadChannelBuilder interface {
	BuilderExecutorReason[*Channel]
	ExpandableCreateThreadChannelBuilder[CreateThreadChannelBuilder]
}

type ExpandableCreateThreadChannelBuilder[B any] interface {
	Archived(archived bool) B
	AutoArchiveDuration(dur ThreadArchiveDuration) B
	Locked(locked bool) B
	Invitable(invitable bool) B
	RateLimitPerUser(limit uint32) B
	Builder() B
}

type UpdateThreadChannelBuilder interface {
	BuilderExecutorReason[*Channel]
	ExpandableUpdateThreadChannelBuilder[UpdateThreadChannelBuilder]
}

type ExpandableUpdateThreadChannelBuilder[B any] interface {
	ExpandableCreateThreadChannelBuilder[B]

	Name(name string) B
}

type GuildChannelBuilder interface {
	BuilderExecutorReason[*Channel]
	ExpandableGuildChannelBuilder[GuildChannelBuilder]
}

type ExpandableGuildChannelBuilder[B any] interface {
	Type(t ChannelType) B
	Topic(topic string) B
	Bitrate(bitrate uint64) B
	UserLimit(limit uint16) B
	RateLimitPerUser(limit uint32) B
	Position(pos int) B
	Parent(id snowflake.ID) B
	NSFW(isNSFW bool) B
	Overwrites(perms []permissions.Overwrite) B
	Builder() B
}

type UpdateGuildChannelBuilder interface {
	BuilderExecutorReason[*Channel]
	ExpandableUpdateGuildChannelBuilder[UpdateGuildChannelBuilder]
}

type ExpandableUpdateGuildChannelBuilder[B any] interface {
	ExpandableGuildChannelBuilder[UpdateGuildChannelBuilder]

	Name(name string) B
}

type CreateThreadTypeSelector interface {
	Public() CreateThreadChannelBuilder
	Private() CreateThreadChannelBuilder
	News() CreateThreadChannelBuilder
}

type UpdateChannelTypeSelector interface {
	Thread() UpdateThreadChannelBuilder
	Guild() UpdateGuildChannelBuilder
}
