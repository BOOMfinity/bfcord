package discord

import (
	"io"

	"github.com/BOOMfinity/bfcord/discord/components"
	"github.com/BOOMfinity/bfcord/internal/timeconv"
	"github.com/andersfylling/snowflake/v5"
)

// Message
//
// Reference: https://discord.com/developers/docs/resources/channel#message-object
type Message struct {
	Reactions ReactionStore `json:"reactions"`
	Member    Member        `json:"member"`
	Author    User          `json:"author"`
	BaseMessage
}

func (v *Message) Fetch(api ClientQuery) error {
	msg, err := v.fetch(api)
	if err != nil {
		return err
	}
	v.fetchPatch(*msg)
	v.Reactions = msg.Reactions
	v.Member = msg.Member
	v.Author = msg.Author
	return nil
}

func (v *Message) Patch() {
	v.BaseMessage.AuthorID = v.Author.ID
	for i := range v.Reactions {
		v.Reactions[i].MessageID = v.ID
		v.Reactions[i].ChannelID = v.ChannelID
		v.Reactions[i].Emoji.GuildID = v.GuildID
	}
	v.Member.UserID = v.Author.ID
	v.Member.GuildID = v.GuildID
}

type BaseMessage struct {
	EditedTimestamp timeconv.Timestamp      `json:"edited_timestamp"`
	Timestamp       timeconv.Timestamp      `json:"timestamp"`
	Reference       *MessageReference       `json:"message_reference"`
	Activity        MessageActivity         `json:"activity"`
	Nonce           string                  `json:"nonce"`
	Content         string                  `json:"content"`
	Attachments     []Attachment            `json:"attachments"`
	Components      []components.Component  `json:"components"`
	Embeds          []MessageEmbed          `json:"embeds"`
	Mentions        []User                  `json:"mentions"`
	MentionRoles    []snowflake.ID          `json:"mention_roles"`
	MentionChannels []MessageChannelMention `json:"mention_channels"`
	Interaction     MessageInteraction      `json:"interaction"`
	ApplicationID   snowflake.ID            `json:"application_id"`
	GuildID         snowflake.ID            `json:"guild_id"`
	ChannelID       snowflake.ID            `json:"channel_id"`
	WebhookID       snowflake.ID            `json:"webhook_id"`
	AuthorID        snowflake.ID            `json:"author_id"`
	ID              snowflake.ID            `json:"id"`
	Type            MessageType             `json:"type"`
	MentionEveryone bool                    `json:"mention_everyone"`
	TTS             bool                    `json:"tts"`
	Pinned          bool                    `json:"pinned"`
}

type MessageInteraction struct {
	Name string       `json:"name"`
	User User         `json:"user"`
	ID   snowflake.ID `json:"id"`
	Type uint8        `json:"type"`
}

func (v BaseMessage) fetch(api ClientQuery) (*Message, error) {
	// TODO: Ignore cache
	return api.Channel(v.ChannelID).Message(v.ID).Get()
}

func (v *BaseMessage) fetchPatch(msg Message) {
	v.Content = msg.Content
	v.Type = msg.Type
	v.Embeds = msg.Embeds
	v.Attachments = msg.Attachments
	v.Pinned = msg.Pinned
}

func (v BaseMessage) Fetch(api ClientQuery) error {
	msg, err := v.fetch(api)
	if err != nil {
		return err
	}
	v.fetchPatch(*msg)
	return nil
}

func (v BaseMessage) IsGuild() bool {
	return v.GuildID.Valid()
}

func (v BaseMessage) API(client ClientQuery) MessageQuery {
	return client.Channel(v.ChannelID).Message(v.ID)
}

func (v BaseMessage) Guild(client ClientQuery) GuildQuery {
	return client.Guild(v.GuildID)
}

func (v BaseMessage) Author(client ClientQuery) UserQuery {
	return client.User(v.AuthorID)
}

func (v BaseMessage) Member(client ClientQuery) GuildMemberQuery {
	return client.Guild(v.GuildID).Member(v.AuthorID)
}

func (v BaseMessage) Edit(client ClientQuery) MessageBuilder {
	return client.Channel(v.ChannelID).Message(v.ID).Edit()
}

func (v BaseMessage) Channel(client ClientQuery) ChannelQuery {
	return client.Channel(v.ChannelID)
}

func (v BaseMessage) Reply(client ClientQuery, ref bool) CreateMessageBuilder {
	bl := v.Channel(client).SendMessage()
	if ref {
		bl.Reference(MessageReference{
			MessageID: v.ID,
			ChannelID: v.ChannelID,
			GuildID:   v.GuildID,
		})
	}
	return bl
}

type MessageActivity struct {
	PartyID string              `json:"party_id"`
	Type    MessageActivityType `json:"type"`
}

type MessageActivityType uint8

const (
	MessageActivityJoin MessageActivityType = iota + 1
	MessageActivitySpectate
	MessageActivityListen
	MessageActivityJoinRequest
)

type MessageType uint8

const (
	MessageTypeDefault MessageType = iota
	MessageTypeRecipientAdd
	MessageTypeRecipientRemove
	MessageTypeCall
	MessageTypeChannelNameChange
	MessageTypeChannelIconChange
	MessageTypeChannelPinnedMessage
	MessageTypeGuildMemberJoin
	MessageTypeGuildSubscription
	MessageTypeGuildSubscriptionTier1
	MessageTypeGuildSubscriptionTier2
	MessageTypeGuildSubscriptionTier3
	MessageTypeChannelFollowAdd
	MessageTypeDiscoveryDisqualified
	MessageTypeDiscoveryReQualified
	MessageTypeDiscoveryInitialWarning
	MessageTypeDiscoveryFinalWarning
	MessageTypeThreadCreated
	MessageTypeReply
	MessageTypeChatInputCommand
	MessageTypeThreadStarterMessage
	MessageTypeGuildInviteReminder
	MessageTypeContextMenuCommand
)

// MessageEmbed
//
// Reference: https://discord.com/developers/docs/resources/channel#embed-object
type MessageEmbed struct {
	Timestamp   *timeconv.Timestamp `json:"timestamp,omitempty"`
	Author      *EmbedAuthor        `json:"author,omitempty"`
	Footer      *EmbedFooter        `json:"footer,omitempty"`
	Provider    *EmbedProvider      `json:"provider,omitempty"`
	Type        EmbedType           `json:"type,omitempty"`
	Description string              `json:"description,omitempty"`
	Url         string              `json:"url,omitempty"`
	Title       string              `json:"title,omitempty"`
	Fields      []EmbedField        `json:"fields,omitempty"`
	Thumbnail   *EmbedMedia         `json:"thumbnail,omitempty"`
	Video       *EmbedMedia         `json:"video,omitempty"`
	Image       *EmbedMedia         `json:"image,omitempty"`
	Color       int64               `json:"color,omitempty"`
}

type EmbedFooter struct {
	Text         string `json:"text"`
	IconUrl      string `json:"icon_url,omitempty"`
	ProxyIconUrl string `json:"proxy_icon_url,omitempty"`
}

type EmbedMedia struct {
	Url      string `json:"url,omitempty"`
	ProxyUrl string `json:"proxy_url,omitempty"`
	Height   int    `json:"height,omitempty"`
	Width    int    `json:"width,omitempty"`
}

type EmbedProvider struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type EmbedAuthor struct {
	Name         string `json:"name"`
	Url          string `json:"url,omitempty"`
	IconUrl      string `json:"icon_url,omitempty"`
	ProxyIconUrl string `json:"proxy_icon_url,omitempty"`
}

type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type EmbedType string

const (
	EmbedTypeRich    EmbedType = "rich"
	EmbedTypeImage   EmbedType = "image"
	EmbedTypeVideo   EmbedType = "video"
	EmbedTypeGif     EmbedType = "gifv"
	EmbedTypeArticle EmbedType = "article"
	EmbedTypeLink    EmbedType = "link"
)

type MessageChannelMention struct {
	Name    string       `json:"name"`
	ID      snowflake.ID `json:"id"`
	GuildID snowflake.ID `json:"guild_id"`
	Type    ChannelType  `json:"type"`
}

type MessageReference struct {
	MessageID snowflake.ID `json:"message_id"`
	ChannelID snowflake.ID `json:"channel_id,omitempty"`
	GuildID   snowflake.ID `json:"guild_id,omitempty"`
}

type MessageCreate struct {
	Content          *string                 `json:"content,omitempty"`
	TTS              *bool                   `json:"tts,omitempty"`
	MessageReference *MessageReference       `json:"message_reference,omitempty"`
	Embeds           *[]MessageEmbed         `json:"embeds,omitempty"`
	Files            *[]MessageFile          `json:"-"`
	Attachments      *[]Attachment           `json:"attachments,omitempty"`
	Components       *[]components.Component `json:"components,omitempty"`
	AllowedMentions  *MessageAllowedMentions `json:"allowed_mentions,omitempty"`
}

// MessageAllowedMentions
//
// Reference: https://discord.com/developers/docs/resources/channel#allowed-mentions-object-allowed-mentions-structure
type MessageAllowedMentions struct {
	// Supported values: roles, users, everyone
	Parse       []string       `json:"parse"`
	Roles       []snowflake.ID `json:"roles,omitempty"`
	Users       []snowflake.ID `json:"users,omitempty"`
	RepliedUser bool           `json:"replied_user,omitempty"`
}

type ForumMessageCreate struct {
	Message             MessageCreate          `json:"message"`
	Name                *string                `json:"name,omitempty"`
	AutoArchiveDuration *ThreadArchiveDuration `json:"auto_archive_duration,omitempty"`
	RateLimitPerUser    uint32                 `json:"rate_limit_per_user,omitempty"`
	AppliedTags         []snowflake.ID         `json:"applied_tags,omitempty"`
}

type MessageFile struct {
	Reader      io.Reader `json:"-"`
	Name        string    `json:"-"`
	Description string    `json:"-"`
	Url         string    `json:"-"`
	Ephemeral   bool      `json:"-"`
	Base64      bool      `json:"-"`
}
