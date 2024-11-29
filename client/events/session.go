package events

import "github.com/BOOMfinity/golog/v2"

type sessionDispatcher struct {
	ready                    Dispatcher[ReadyEvent]
	guildCreate              Dispatcher[GuildCreateEvent]
	guildUpdate              Dispatcher[GuildUpdateEvent]
	guildDelete              Dispatcher[GuildDeleteEvent]
	channelCreate            Dispatcher[ChannelCreateEvent]
	channelUpdate            Dispatcher[ChannelUpdateEvent]
	channelDelete            Dispatcher[ChannelDeleteEvent]
	channelPinsUpdate        Dispatcher[ChannelPinsUpdateEvent]
	messageCreate            Dispatcher[MessageCreateEvent]
	messageUpdate            Dispatcher[MessageUpdateEvent]
	messageDelete            Dispatcher[MessageDeleteEvent]
	threadCreate             Dispatcher[ThreadCreateEvent]
	threadUpdate             Dispatcher[ThreadUpdateEvent]
	threadDelete             Dispatcher[ThreadDeleteEvent]
	threadListSync           Dispatcher[ThreadListSyncEvent]
	threadMembersUpdate      Dispatcher[ThreadMembersUpdateEvent]
	guildRoleAdd             Dispatcher[GuildRoleAddEvent]
	guildRoleUpdate          Dispatcher[GuildRoleUpdateEvent]
	guildRoleDelete          Dispatcher[GuildRoleDeleteEvent]
	guildScheduledCreate     Dispatcher[GuildScheduledCreateEvent]
	guildScheduledUpdate     Dispatcher[GuildScheduledUpdateEvent]
	guildScheduledDelete     Dispatcher[GuildScheduledDeleteEvent]
	guildScheduledUserAdd    Dispatcher[GuildScheduledUserAddEvent]
	guildScheduledUserRemove Dispatcher[GuildScheduledUserRemoveEvent]
	guildMemberAdd           Dispatcher[GuildMemberAddEvent]
	guildMemberRemove        Dispatcher[GuildMemberRemoveEvent]
	guildMemberUpdate        Dispatcher[GuildMemberUpdateEvent]
	inviteCreate             Dispatcher[InviteCreateEvent]
	inviteDelete             Dispatcher[InviteDeleteEvent]
	guildBanAdd              Dispatcher[GuildBanAddEvent]
	guildBanRemove           Dispatcher[GuildBanRemoveEvent]
	interactionCreate        Dispatcher[InteractionCreateEvent]
	voiceStateUpdate         Dispatcher[VoiceStateUpdateEvent]
	voiceServerUpdate        Dispatcher[VoiceServerUpdateEvent]
}

func (s *sessionDispatcher) Ready() Dispatcher[ReadyEvent] {
	return s.ready
}

func (s *sessionDispatcher) GuildCreate() Dispatcher[GuildCreateEvent] {
	return s.guildCreate
}

func (s *sessionDispatcher) GuildDelete() Dispatcher[GuildDeleteEvent] {
	return s.guildDelete
}

func (s *sessionDispatcher) ChannelCreate() Dispatcher[ChannelCreateEvent] {
	return s.channelCreate
}

func (s *sessionDispatcher) ChannelUpdate() Dispatcher[ChannelUpdateEvent] {
	return s.channelUpdate
}

func (s *sessionDispatcher) ChannelPinsUpdate() Dispatcher[ChannelPinsUpdateEvent] {
	return s.channelPinsUpdate
}

func (s *sessionDispatcher) ChannelDelete() Dispatcher[ChannelDeleteEvent] {
	return s.channelDelete
}

func (s *sessionDispatcher) MessageCreate() Dispatcher[MessageCreateEvent] {
	return s.messageCreate
}

func (s *sessionDispatcher) MessageUpdate() Dispatcher[MessageUpdateEvent] {
	return s.messageUpdate
}

func (s *sessionDispatcher) MessageDelete() Dispatcher[MessageDeleteEvent] {
	return s.messageDelete
}

func (s *sessionDispatcher) GuildUpdate() Dispatcher[GuildUpdateEvent] {
	return s.guildUpdate
}

func (s *sessionDispatcher) ThreadCreate() Dispatcher[ThreadCreateEvent] {
	return s.threadCreate
}

func (s *sessionDispatcher) ThreadUpdate() Dispatcher[ThreadUpdateEvent] {
	return s.threadUpdate
}

func (s *sessionDispatcher) ThreadDelete() Dispatcher[ThreadDeleteEvent] {
	return s.threadDelete
}

func (s *sessionDispatcher) ThreadListSync() Dispatcher[ThreadListSyncEvent] {
	return s.threadListSync
}

func (s *sessionDispatcher) ThreadMembersUpdate() Dispatcher[ThreadMembersUpdateEvent] {
	return s.threadMembersUpdate
}

func (s *sessionDispatcher) GuildRoleAdd() Dispatcher[GuildRoleAddEvent] {
	return s.guildRoleAdd
}

func (s *sessionDispatcher) GuildRoleUpdate() Dispatcher[GuildRoleUpdateEvent] {
	return s.guildRoleUpdate
}

func (s *sessionDispatcher) GuildRoleDelete() Dispatcher[GuildRoleDeleteEvent] {
	return s.guildRoleDelete
}

func (s *sessionDispatcher) GuildScheduledCreate() Dispatcher[GuildScheduledCreateEvent] {
	return s.guildScheduledCreate
}

func (s *sessionDispatcher) GuildScheduledUpdate() Dispatcher[GuildScheduledUpdateEvent] {
	return s.guildScheduledUpdate
}

func (s *sessionDispatcher) GuildScheduledDelete() Dispatcher[GuildScheduledDeleteEvent] {
	return s.guildScheduledDelete
}

func (s *sessionDispatcher) GuildScheduledUserAdd() Dispatcher[GuildScheduledUserAddEvent] {
	return s.guildScheduledUserAdd
}

func (s *sessionDispatcher) GuildScheduledUserRemove() Dispatcher[GuildScheduledUserRemoveEvent] {
	return s.guildScheduledUserRemove
}

func (s *sessionDispatcher) GuildMemberAdd() Dispatcher[GuildMemberAddEvent] {
	return s.guildMemberAdd
}

func (s *sessionDispatcher) GuildMemberRemove() Dispatcher[GuildMemberRemoveEvent] {
	return s.guildMemberRemove
}

func (s *sessionDispatcher) GuildMemberUpdate() Dispatcher[GuildMemberUpdateEvent] {
	return s.guildMemberUpdate
}

func (s *sessionDispatcher) InviteCreate() Dispatcher[InviteCreateEvent] {
	return s.inviteCreate
}

func (s *sessionDispatcher) InviteDelete() Dispatcher[InviteDeleteEvent] {
	return s.inviteDelete
}

func (s *sessionDispatcher) GuildBanAdd() Dispatcher[GuildBanAddEvent] {
	return s.guildBanAdd
}

func (s *sessionDispatcher) GuildBanRemove() Dispatcher[GuildBanRemoveEvent] {
	return s.guildBanRemove
}

func (s *sessionDispatcher) InteractionCreate() Dispatcher[InteractionCreateEvent] {
	return s.interactionCreate
}

func (s *sessionDispatcher) VoiceStateUpdate() Dispatcher[VoiceStateUpdateEvent] {
	return s.voiceStateUpdate
}

func (s *sessionDispatcher) VoiceServerUpdate() Dispatcher[VoiceServerUpdateEvent] {
	return s.voiceServerUpdate
}

func NewSessionDispatcher(log golog.Logger) SessionDispatcher {
	return &sessionDispatcher{
		ready:                    NewDispatcher[ReadyEvent](log),
		guildCreate:              NewDispatcher[GuildCreateEvent](log),
		guildUpdate:              NewDispatcher[GuildUpdateEvent](log),
		guildDelete:              NewDispatcher[GuildDeleteEvent](log),
		channelCreate:            NewDispatcher[ChannelCreateEvent](log),
		channelUpdate:            NewDispatcher[ChannelUpdateEvent](log),
		channelDelete:            NewDispatcher[ChannelDeleteEvent](log),
		channelPinsUpdate:        NewDispatcher[ChannelPinsUpdateEvent](log),
		messageCreate:            NewDispatcher[MessageCreateEvent](log),
		messageUpdate:            NewDispatcher[MessageUpdateEvent](log),
		messageDelete:            NewDispatcher[MessageDeleteEvent](log),
		threadCreate:             NewDispatcher[ThreadCreateEvent](log),
		threadUpdate:             NewDispatcher[ThreadUpdateEvent](log),
		threadDelete:             NewDispatcher[ThreadDeleteEvent](log),
		threadListSync:           NewDispatcher[ThreadListSyncEvent](log),
		threadMembersUpdate:      NewDispatcher[ThreadMembersUpdateEvent](log),
		guildRoleAdd:             NewDispatcher[GuildRoleAddEvent](log),
		guildRoleUpdate:          NewDispatcher[GuildRoleUpdateEvent](log),
		guildRoleDelete:          NewDispatcher[GuildRoleDeleteEvent](log),
		guildScheduledCreate:     NewDispatcher[GuildScheduledCreateEvent](log),
		guildScheduledUpdate:     NewDispatcher[GuildScheduledUpdateEvent](log),
		guildScheduledDelete:     NewDispatcher[GuildScheduledDeleteEvent](log),
		guildScheduledUserAdd:    NewDispatcher[GuildScheduledUserAddEvent](log),
		guildScheduledUserRemove: NewDispatcher[GuildScheduledUserRemoveEvent](log),
		guildMemberAdd:           NewDispatcher[GuildMemberAddEvent](log),
		guildMemberRemove:        NewDispatcher[GuildMemberRemoveEvent](log),
		guildMemberUpdate:        NewDispatcher[GuildMemberUpdateEvent](log),
		inviteCreate:             NewDispatcher[InviteCreateEvent](log),
		inviteDelete:             NewDispatcher[InviteDeleteEvent](log),
		guildBanAdd:              NewDispatcher[GuildBanAddEvent](log),
		guildBanRemove:           NewDispatcher[GuildBanRemoveEvent](log),
		interactionCreate:        NewDispatcher[InteractionCreateEvent](log),
		voiceStateUpdate:         NewDispatcher[VoiceStateUpdateEvent](log),
		voiceServerUpdate:        NewDispatcher[VoiceServerUpdateEvent](log),
	}
}
