package discord

import (
	"github.com/BOOMfinity/bfcord/discord/components"
	"github.com/andersfylling/snowflake/v5"
)

type MessageBuilder interface {
	BuilderExecutor[Message]
	BaseMessageBuilder[MessageBuilder]
}

type WebhookExecuteBuilder interface {
	Execute() (msg *Message, err error)
	ExpandableWebhookExecuteBuilder[WebhookExecuteBuilder]
}

type WebhookUpdateMessageBuilder interface {
	Execute() (msg *Message, err error)
	ExpandableWebhookUpdateMessageBuilder[WebhookUpdateMessageBuilder]
}

type ExpandableWebhookUpdateMessageBuilder[B any] interface {
	BaseMessageBuilder[B]
	Thread(id snowflake.ID) B
}

type ExpandableWebhookExecuteBuilder[B any] interface {
	BaseMessageBuilder[B]
	AvatarURL(url string) B
	Username(name string) B
	Thread(id snowflake.ID) B
	NoWait() B
}

type BaseMessageBuilder[B any] interface {
	Content(str string) B
	Embed(embed MessageEmbed) B
	Embeds(embeds []MessageEmbed) B
	// Deprecated: Use ActionRow or AutoActionRows instead.
	Components(list []components.Component) B
	// ActionRow appends new row with given components.
	ActionRow(items ...components.ActionRowItem) B
	// AutoActionRows will automatically split components into rows.
	AutoActionRows(items ...components.ActionRowItem) B
	File(f MessageFile) B
	Files(f []MessageFile) B
	KeepFiles(files []Attachment) B
	AllowedMentions(allowed MessageAllowedMentions) B
	DoNotKeepFiles() B

	ClearEmbeds() B
	ClearFiles() B
	ClearContent() B
	ClearComponents() B
	ClearAllowedMentions() B

	Raw() MessageCreate
	Builder() B
}

type BuilderExecutorReason[R any] interface {
	Execute(api ClientQuery, reason ...string) (*R, error)
}

type BuilderExecutor[R any] interface {
	Execute(api ClientQuery) (*R, error)
}

type CreateMessageBuilder interface {
	BuilderExecutor[Message]
	_createMessageBuilder[CreateMessageBuilder]
}

type _createMessageBuilder[B any] interface {
	BaseMessageBuilder[B]
	Reference(ref MessageReference) B
	TTS() B
	SuppressEmbeds() B
}

type CreateForumMessageBuilder interface {
	BuilderExecutor[ChannelWithMessage]
	_createForumMessageBuilder[CreateForumMessageBuilder]
}

type _createForumMessageBuilder[B any] interface {
	Content(str string) B
	Embed(embed MessageEmbed) B
	Embeds(embeds []MessageEmbed) B
	// ActionRow appends new row with given components.
	ActionRow(items ...components.ActionRowItem) B
	// AutoActionRows will automatically split components into rows.
	AutoActionRows(items ...components.ActionRowItem) B
	File(f MessageFile) B
	Files(f []MessageFile) B
	RawForum() ForumMessageCreate
	Builder() B
	AutoArchiveDuration(dur ThreadArchiveDuration) B
	RateLimitPerUser(limit uint32) B
	Tags(t []snowflake.ID) B
}
