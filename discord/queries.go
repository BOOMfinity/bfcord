package discord

import (
	"context"
	"github.com/BOOMfinity/bfcord/discord/permissions"
	"github.com/BOOMfinity/golog"
	"github.com/andersfylling/snowflake/v5"
)

type QueryOptions[V any] interface {
	NoCache() V
	NoAPI() V
	Reason(reason string) V
}

type UserQuery interface {
	QueryOptions[UserQuery]
	// Get returns User type from cache (if used) or directly from Discord API
	Get() (user *User, err error)
	// Send creates DM channel with user and sends defined message
	Send() (msg CreateMessageBuilder, err error)
	// CreateDM creates a private channel with user
	CreateDM() (ch *Channel, err error)
	ID() snowflake.ID
}

type WebhookQuery interface {
	Fetch() (wh *Webhook, err error)
	Execute() WebhookExecuteBuilder
	Delete() (err error)
	DeleteMessage(id snowflake.ID) (err error)
	EditMessage(id snowflake.ID) WebhookUpdateMessageBuilder
}

type LowLevelClientQuery interface {
	QueryOptions[LowLevelClientQuery]
	CreateMessage(channel snowflake.ID, data MessageCreate) (msg *Message, err error)
	UpdateMessage(channel snowflake.ID, message snowflake.ID, data MessageCreate) (msg *Message, err error)
	ExecuteWebhook(id snowflake.ID, token string, data WebhookExecute, wait bool, thread snowflake.ID) (msg *Message, err error)
	UpdateWebhookMessage(id snowflake.ID, token string, message snowflake.ID, data MessageCreate, thread snowflake.ID) (msg *Message, err error)
	Message(method string, url string, _data any) (msg *Message, err error)
	UpdateChannel(id snowflake.ID, data ChannelUpdate) (ch *Channel, err error)
	CreateGuildChannel(guild snowflake.ID, data ChannelUpdate) (ch *Channel, err error)
	UpdateGuildMember(guild snowflake.ID, member snowflake.ID, data MemberUpdate) (m *MemberWithUser, err error)
	UpdateGuild(guild snowflake.ID, data GuildUpdate) (g *Guild, err error)
	SendDM(channel snowflake.ID) (msg CreateMessageBuilder)
	StartThread(channel snowflake.ID, message snowflake.ID, data ThreadCreate) (ch *Channel, err error)
	CreateOrUpdate(guild, role snowflake.ID, data RoleCreate) (r *Role, err error)
	CreateForumMessage(id snowflake.ID, data ForumMessageCreate) (d *ChannelWithMessage, err error)
}

type ChannelMessagesQuery interface {
	Around(id snowflake.ID, limit uint16) (msgs []Message, err error)
	After(ctx context.Context, id snowflake.ID, limit uint16) (msgs []Message, err error)
	Before(ctx context.Context, id snowflake.ID, limit uint16) (msgs []Message, err error)
	Latest(limit uint16) (msgs []Message, err error)
	ID() snowflake.ID
}

type ChannelQuery interface {
	QueryOptions[ChannelQuery]
	Message(id snowflake.ID) MessageQuery
	SendMessage() CreateMessageBuilder
	Get() (ch *Channel, err error)
	Edit() UpdateChannelTypeSelector
	Delete() error
	Messages() ChannelMessagesQuery
	Bulk(ids []snowflake.ID) error
	// EditPermissions
	// DeletePermission(id snowflake.ID) error
	Invites() (invites []InviteWithMeta, err error)
	//CreateInvite() (inv *Invite, err error)
	Follow(target snowflake.ID) error
	Pinned() (msg []Message, err error)
	StartThread(name string) CreateThreadTypeSelector
	StartForumThread(name string) CreateForumMessageBuilder
	Join() error
	AddMember(id snowflake.ID) error
	Leave() error
	RemoveMember(id snowflake.ID) error
	GetThreadMember(id snowflake.ID) (tm *ThreadMember, err error)
	Stage() StageQuery
	ID() snowflake.ID
}

type StageQuery interface {
	QueryOptions[StageQuery]
	Create(topic string, notify bool) (stage *StageInstance, err error)
	Get() (stage *StageInstance, err error)
	Modify(topic string) (stage *StageInstance, err error)
	Delete() error
}

type InviteQuery interface {
}

type MessageQuery interface {
	Edit() MessageBuilder
	Delete() (err error)
	React(emoji string) error
	Reaction(emoji string) MessageReactionQuery
	RemoveAllReactions() (err error)
	Get() (msg *Message, err error)
	CrossPost() error
	Pin() error
	UnPin() error
	StartThread(name string) CreateThreadChannelBuilder
	ChannelID() snowflake.ID
	ID() snowflake.ID
}

type MessageReactionQuery interface {
	QueryOptions[MessageReactionQuery]
	RemoveOwn() (err error)
	After(limit uint64, after snowflake.ID) (users []User, err error)
	Range(limit uint64, after snowflake.ID) (users []User, err error)
	All(limit uint64) (users []User, err error)
	Remove(userID snowflake.ID) (err error)
	RemoveAll() (err error)
	Emoji() string
	Message() snowflake.ID
	Channel() snowflake.ID
}

type GuildMemberQuery interface {
	QueryOptions[GuildMemberQuery]
	Get() (member *MemberWithUser, err error)
	Ban(days uint8) (err error)
	Unban() (err error)
	AddRole(role snowflake.ID) (err error)
	RemoveRole(role snowflake.ID) (err error)
	Kick() (err error)
	Edit() UpdateGuildMemberBuilder
	Permissions() (perm permissions.Permission, err error)
	PermissionsIn(channel snowflake.ID) (perm permissions.Permission, err error)
	VoiceState() (state VoiceState, err error)
	ID() snowflake.ID
	GuildID() snowflake.ID
}

type RoleQuery interface {
	QueryOptions[RoleQuery]
	Get() (role *Role, err error)
	Edit() RoleBuilder
	Delete() error
}

type GuildQuery interface {
	QueryOptions[GuildQuery]
	Get() (guild *Guild, err error)
	Delete() (err error)
	Channels() (channels []Channel, err error)
	CreateChannel(name string) GuildChannelBuilder
	Edit() GuildBuilder
	UpdateChannelPositions(positions *GuildChannelPositionsBuilder) (err error)
	ActiveThreads() (threads []Channel, err error)
	Member(id snowflake.ID) GuildMemberQuery
	Search(query string, limit uint16) (members []MemberWithUser, err error)
	SetCurrentNick(nick string) (err error)
	Roles() (roles []Role, err error)
	VoiceStates() (states Slice[VoiceState], err error)
	Invites() (invites []InviteWithMeta, err error)
	Role(id snowflake.ID) RoleQuery
	CreateRole() RoleBuilder
	UpdateRolePositions(roles RolePositions) error
	// Members gets the guild members from API and store them in cache (if enabled)
	//
	// limit variable has no upper value, so if you set it to more than 1000 (theoretically Discord API limit) the library will just make appropriate number of requests.
	// if limit is set to -1, bfcord will try to fetch all guild members.
	Members(limit int, after snowflake.ID) (members []MemberWithUser, err error)
}

type ClientQuery interface {
	// User returns user-specific Discord API methods
	User(id snowflake.ID) UserQuery
	// CurrentUser returns current logged user details
	CurrentUser() (user *User, err error)
	Channel(id snowflake.ID) ChannelQuery
	Guild(id snowflake.ID) GuildQuery
	// Logger returns instance of global logger
	Log() golog.Logger

	LowLevel() LowLevelClientQuery

	Webhook(id snowflake.ID, token string) WebhookQuery
}
