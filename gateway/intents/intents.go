package intents

type Intent int

const (
	Guilds Intent = 1 << iota
	GuildMembers
	GuildBans
	GuildEmojisAndStickers
	GuildIntegrations
	GuildWebhooks
	GuildInvites
	GuildVoiceStates
	GuildPresences
	GuildMessages
	GuildMessageReactions
	GuildMessageTyping
	DirectMessages
	DirectMessageReactions
	DirectMessageTyping
	MessageContent
	GuildScheduledEvents

	Default = Guilds |
		GuildMembers |
		GuildBans |
		GuildEmojisAndStickers |
		GuildIntegrations |
		GuildWebhooks |
		GuildInvites |
		GuildVoiceStates |
		GuildMessages |
		GuildMessageReactions |
		GuildScheduledEvents
)
