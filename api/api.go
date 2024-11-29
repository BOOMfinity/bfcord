package api

import (
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/andersfylling/snowflake/v5"
)

type GuildClient interface {
	UpdateChannelPositions(positions []GuildChannelPosition, reason ...string) error
	UpdateRolePositions(positions []GuildRolePosition, reason ...string) error
	CreateChannel(params GuildChannelParams, reason ...string) (discord.Channel, error)
	Get() (discord.Guild, error)
	Modify(params ModifyGuildParams, reason ...string) (discord.Guild, error)
	Delete(reason ...string) error
	Channels() ([]discord.Channel, error)
	ActiveThreads() ([]discord.Channel, error)
	Member(id snowflake.ID) MemberClient
	Members(params GuildMembersParams) ([]discord.Member, error)
	Search(query string, limit ...uint) ([]discord.Member, error)
	ModifyCurrentMember(nick string, reason ...string) (discord.Member, error)
	BulkBan(ids []snowflake.ID, seconds uint, reason ...string) (GuildBanAddResponse, error)
	Roles() ([]discord.Role, error)
	Role(id snowflake.ID) RoleClient
	CreateRole(params CreateRoleParams, reason ...string) (discord.Role, error)
	CurrentUserVoiceState() (discord.VoiceState, error)
	ModifyCurrentUserVoiceState(params ModifyCurrentUserVoiceStateParams) (discord.VoiceState, error)
	Emojis() ([]discord.Emoji, error)
	Emoji(id snowflake.ID) EmojiClient
	CreateEmoji(params CreateEmojiParams, reason ...string) (discord.Emoji, error)
	Events() ([]discord.ScheduledEvent, error)
	CreateEvent(params CreateScheduledEventParams, reason ...string) (discord.ScheduledEvent, error)
	Event(id snowflake.ID) GuildEventClient
}

type ChannelClient interface {
	JoinThread() error
	AddThreadMember(id snowflake.ID) error
	LeaveThread() error
	RemoveThreadMember(id snowflake.ID) error
	ThreadMember(id snowflake.ID, withMember bool) (discord.ThreadMember, error)
	AddRecipient(id snowflake.ID, userToken, userNickname string) error
	RemoveRecipient(id snowflake.ID) error
	Modify(params ModifyChannelParams, reason ...string) (discord.Channel, error)
	StartThread(data StartThreadWithoutMessageParams, reason ...string) (discord.Channel, error)
	StartForumMediaThread(data StartForumOrMediaThreadParams, reason ...string) (discord.Channel, error)
	UpdateChannelPermissions(id snowflake.ID, data UpdateChannelPermissionsParams, reason ...string) error
	DeleteChannelPermission(id snowflake.ID, reason ...string) error
	FollowAnnouncementChannel(webhook snowflake.ID, reason ...string) (discord.FollowedChannel, error)
	Pins() ([]discord.Message, error)
	Invites() ([]discord.Invite, error)
	CreateInvite(data CreateChannelInviteParams, reason ...string) (discord.Invite, error)
	Delete(reason ...string) error
	Webhooks() ([]discord.Webhook, error)
	CreateWebhook(params CreateWebhookParams, reason ...string) (discord.Webhook, error)
	Get() (dst discord.Channel, err error)
	SendMessage(params CreateMessageParams) (discord.Message, error)
	Messages() MessagesQuery
	Message(id snowflake.ID) MessageClient
	BulkDelete(messages []snowflake.ID, reason ...string) error
}

type EmojiClient interface {
	Get() (discord.Emoji, error)
	Modify(params ModifyEmojiParams, reason ...string) (discord.Emoji, error)
	Delete(reason ...string) error
}

type GuildEventClient interface {
	Get(withUserCount bool) (discord.ScheduledEvent, error)
	Modify(params ModifyScheduledEventParams, reason ...string) (discord.ScheduledEvent, error)
	Delete(reason ...string) error
	Users() GuildEventQuery
}

type GuildEventQuery interface {
	Before(id snowflake.ID, withMember bool, limit uint) ([]ScheduledEventUser, error)
	After(id snowflake.ID, withMember bool, limit uint) ([]ScheduledEventUser, error)
	Latest(withMember bool, limit uint) ([]ScheduledEventUser, error)
}

type InteractionClient interface {
	Pong() error
	Response() InteractionResponseClient
	Update(params EditMessageParams) error
	DeferredUpdate() error
	DeferredComponentUpdate() error
	Reply(params InteractionMessageParams) error
	SendFollowUp(params FollowUpParams) (discord.Message, error)
	FollowUp(id snowflake.ID) FollowUpClient
	AutoComplete(choices []discord.CommandChoice) error
	TextInput(params TextInputParams) error
}

type FollowUpClient interface {
	Get() (discord.Message, error)
	Update(params EditMessageParams) (discord.Message, error)
	Delete() error
}

type InteractionResponseClient interface {
	Get() (discord.Message, error)
	Delete() error
	Edit(params EditMessageParams) (discord.Message, error)
}

type MemberClient interface {
	Get() (discord.MemberWithUser, error)
	Modify(params ModifyGuildMemberParams, reason ...string) (discord.Member, error)
	AddRole(id snowflake.ID, reason ...string) error
	RemoveRole(id snowflake.ID, reason ...string) error
	Kick(reason ...string) error
	CreateBan(seconds uint, reason ...string) error
	RemoveBan(reason ...string) error
	VoiceState() (discord.VoiceState, error)
	ModifyVoiceState(params ModifyUserVoiceStateParams) (discord.VoiceState, error)
}

type MessageClient interface {
	Answer(id uint, params MessagePollVotersParams) ([]discord.User, error)
	EndPoll() (discord.Message, error)
	Reaction(emoji string) ReactionClient
	StartThread(data StartThreadParams, reason ...string) (discord.Channel, error)
	Get() (dst discord.Message, _ error)
	Delete(reason ...string) error
	Pin(reason ...string) error
	Unpin(reason ...string) error
	Update(params EditMessageParams) (discord.Message, error)
	CrossPost() error
	DeleteAllReactions() error
}

type MessagesQuery interface {
	Latest(limit uint) ([]discord.Message, error)
	Before(id snowflake.ID, limit uint) ([]discord.Message, error)
	After(id snowflake.ID, limit uint) ([]discord.Message, error)
	Around(id snowflake.ID, limit uint) ([]discord.Message, error)
}

type ReactionClient interface {
	Reactions(opts MessageReactionsParams) (_ []discord.User, err error)
	React() error
	DeleteOwn() error
	Delete(user snowflake.ID) error
	DeleteAll() error
}

type RoleClient interface {
	Get() (discord.Role, error)
	Modify(params CreateRoleParams, reason ...string) (discord.Role, error)
	Delete(reason ...string) error
}

type SlashClient interface {
	Command(id snowflake.ID) SlashCommandClient
	Guild(id snowflake.ID) GuildSlashCommandClient
	List() ([]discord.Command, error)
	Create(params discord.CreateCommand) (discord.Command, error)
	Bulk(cmds []discord.CreateCommand) ([]discord.Command, error)
}

type SlashCommandClient interface {
	Get() (discord.Command, error)
	Delete() error
	Modify(params discord.CreateCommand) (discord.Command, error)
	ModifyPermissions(params discord.CommandPermissions) (discord.CommandPermissions, error)
	Permissions() (discord.CommandPermissions, error)
}

type GuildSlashCommandClient interface {
	PermissionList() ([]discord.CommandPermissions, error)
	List() ([]discord.Command, error)
	Create(params discord.CreateCommand) (discord.Command, error)
	Bulk(cmds []discord.CreateCommand) ([]discord.Command, error)
	Command(id snowflake.ID) SlashCommandClient
}

type StageClient interface {
	Get() (discord.StageInstance, error)
	Create(params CreateStageInstanceParams, reason ...string) (discord.StageInstance, error)
	Modify(params ModifyStageInstanceParams, reason ...string) (discord.StageInstance, error)
	Delete(reason ...string) error
}

type UserClient interface {
	Get() (u discord.User, _ error)
	CreateDM(recipient snowflake.ID) (discord.Channel, error)
}
