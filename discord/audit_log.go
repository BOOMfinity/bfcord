package discord

import "github.com/andersfylling/snowflake/v5"

type AuditLog struct {
	AuditLogEntries      []AuditLogEntry
	GuildScheduledEvents []ScheduledEvent
	Threads              []Channel
	Users                []User
	Webhooks             []Webhook
	// TODO: auto moderation, commands
}

type AuditLogEntry struct {
	TargetID   snowflake.ID
	Changes    []AuditLogChange
	UserID     snowflake.ID
	ID         snowflake.ID
	ActionType AuditLogActionType
	Options    AuditLogOptions
	Reason     string
}

type AuditLogChange struct {
	NewValue any
	OldValue any
	Key      string
}

type AuditLogOptions struct {
	ApplicationID                 snowflake.ID
	AutoModerationRuleName        string
	AutoModerationRuleTriggerType string
	ChannelID                     snowflake.ID
	Count                         int
	DeleteMemberDays              int
	ID                            snowflake.ID
	MembersRemoved                int
	MessageID                     snowflake.ID
	RoleName                      string
	Type                          PermissionOverwriteType
	IntegrationType               string
}

type AuditLogActionType uint

const (
	AuditLogActionGuildUpdate                             AuditLogActionType = 1
	AuditLogActionChannelCreate                                              = 10
	AuditLogActionChannelUpdate                                              = 11
	AuditLogActionChannelDelete                                              = 12
	AuditLogActionOverwriteCreate                                            = 13
	AuditLogActionOverwriteUpdate                                            = 14
	AuditLogActionOverwriteDelete                                            = 15
	AuditLogActionMemberKick                                                 = 20
	AuditLogActionMemberPrune                                                = 21
	AuditLogActionMemberBanAdd                                               = 22
	AuditLogActionMemberBanRemove                                            = 23
	AuditLogActionMemberUpdate                                               = 24
	AuditLogActionMemberRoleUpdate                                           = 25
	AuditLogActionMemberMove                                                 = 26
	AuditLogActionMemberDisconnect                                           = 27
	AuditLogActionBotAdd                                                     = 28
	AuditLogActionRoleCreate                                                 = 30
	AuditLogActionRoleUpdate                                                 = 31
	AuditLogActionRoleDelete                                                 = 32
	AuditLogActionInviteCreate                                               = 40
	AuditLogActionInviteUpdate                                               = 41
	AuditLogActionInviteDelete                                               = 42
	AuditLogActionWebhookCreate                                              = 50
	AuditLogActionWebhookUpdate                                              = 51
	AuditLogActionWebhookDelete                                              = 52
	AuditLogActionEmojiCreate                                                = 60
	AuditLogActionEmojiUpdate                                                = 61
	AuditLogActionEmojiDelete                                                = 62
	AuditLogActionMessageDelete                                              = 72
	AuditLogActionMessageBulkDelete                                          = 73
	AuditLogActionMessagePin                                                 = 74
	AuditLogActionMessageUnpin                                               = 75
	AuditLogActionIntegrationCreate                                          = 80
	AuditLogActionIntegrationUpdate                                          = 81
	AuditLogActionIntegrationDelete                                          = 82
	AuditLogActionStageCreate                                                = 83
	AuditLogActionStageUpdate                                                = 84
	AuditLogActionStageDelete                                                = 85
	AuditLogActionStickerCreate                                              = 90
	AuditLogActionStickerUpdate                                              = 91
	AuditLogActionStickerDelete                                              = 92
	AuditLogActionScheduledCreate                                            = 100
	AuditLogActionScheduledUpdate                                            = 101
	AuditLogActionScheduledDelete                                            = 102
	AuditLogActionThreadCreate                                               = 110
	AuditLogActionThreadUpdate                                               = 111
	AuditLogActionThreadDelete                                               = 112
	AuditLogActionCommandPermissionUpdate                                    = 121
	AuditLogActionAutoModerationRuleCreate                                   = 140
	AuditLogActionAutoModerationRuleUpdate                                   = 141
	AuditLogActionAutoModerationRuleDelete                                   = 142
	AuditLogActionAutoModerationBlockMessage                                 = 143
	AuditLogActionAutoModerationFlagToChannel                                = 144
	AuditLogActionAutoModerationUserCommunicationDisabled                    = 145
	AuditLogActionMonetizationRequestCreated                                 = 150
	AuditLogActionMonetizationTermsAccepted                                  = 151
	AuditLogActionOnboardingPromptCreate                                     = 163
	AuditLogActionOnboardingPromptUpdate                                     = 164
	AuditLogActionOnboardingPromptDelete                                     = 165
	AuditLogActionOnboardingCreate                                           = 166
	AuditLogActionOnboardingUpdate                                           = 167
	AuditLogActionHomeSettingsCreate                                         = 190
	AuditLogActionHomeSettingsUpdate                                         = 191
)
