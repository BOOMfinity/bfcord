package permissions

import (
	bytes2 "bytes"
	"strconv"
)

var nullBytes = []byte("null")

type Permission uint64

func (b *Permission) Add(a ...Permission) {
	for i := range a {
		*b |= a[i]
	}
}

func (b *Permission) Clear(a ...Permission) {
	for i := range a {
		*b &^= a[i]
	}
}

func (x Permission) Has(perm Permission) bool {
	return x&perm == perm
}

func (x Permission) Administrator() bool {
	return x.Has(Administrator)
}

func (b Permission) MarshalJSON() (dst []byte, err error) {
	str := []byte(strconv.FormatUint(uint64(b), 10))
	dst = append(dst, '"')
	dst = append(dst, str...)
	dst = append(dst, '"')
	return
}

func (b *Permission) UnmarshalJSON(bytes []byte) (err error) {
	if bytes2.Equal(bytes, nullBytes) {
		*b = 0
		return nil
	}
	bytes = bytes[1:]
	bytes = bytes[:len(bytes)-1]
	d, err := strconv.ParseUint(string(bytes), 10, 64)
	if err != nil {
		return
	}
	*b = Permission(d)
	return
}

func (b Permission) Serialize() (data map[Permission]bool) {
	data = make(map[Permission]bool, len(permissions))
	for perm := range permissions {
		data[perm] = b.Has(perm)
	}
	return
}

const (
	CreateInstantInvite Permission = 1 << iota
	KickMembers
	BanMembers
	Administrator
	ManageChannels
	ManageGuild
	AddReactions
	ViewAuditLog
	PrioritySpeaker
	Stream
	ViewChannel
	SendMessages
	SendTTSMessages
	ManageMessages
	EmbedLinks
	AttachFiles
	ReadMessageHistory
	MentionEveryone
	UseExternalEmojis
	ViewGuildInsights
	Connect
	Speak
	MuteMembers
	DeafenMembers
	MoveMembers
	UseVAD
	ChangeNickname
	ManageNicknames
	ManageRoles
	ManageWebhooks
	ManageEmojisAndStickers
	UseApplicationCommands
	RequestToSpeak
	ManageEvents
	ManageThreads
	CreatePublicThreads
	CreatePrivateThreads
	UseExternalStickers
	SendMessagesInThreads
	UseEmbeddedActivities
	ModerateMembers

	All = CreateInstantInvite |
		KickMembers |
		BanMembers |
		Administrator |
		ManageChannels |
		ManageGuild |
		AddReactions |
		ViewAuditLog |
		PrioritySpeaker |
		Stream |
		ViewChannel |
		SendMessages |
		SendTTSMessages |
		ManageMessages |
		EmbedLinks |
		AttachFiles |
		ReadMessageHistory |
		MentionEveryone |
		UseExternalEmojis |
		ViewGuildInsights |
		Connect |
		Speak |
		MuteMembers |
		DeafenMembers |
		MoveMembers |
		UseVAD |
		ChangeNickname |
		ManageNicknames |
		ManageRoles |
		ManageWebhooks |
		ManageEmojisAndStickers |
		UseApplicationCommands |
		RequestToSpeak |
		ManageEvents |
		ManageThreads |
		CreatePublicThreads |
		CreatePrivateThreads |
		UseExternalStickers |
		SendMessagesInThreads |
		UseEmbeddedActivities |
		ModerateMembers
)

var permissions = map[Permission]string{
	CreateInstantInvite:     "createInstantInvite",
	KickMembers:             "kickMembers",
	BanMembers:              "banMembers",
	Administrator:           "administrator",
	ManageChannels:          "manageChannels",
	ManageGuild:             "manageGuild",
	AddReactions:            "addReactions",
	ViewAuditLog:            "viewAuditLog",
	PrioritySpeaker:         "prioritySpeaker",
	Stream:                  "stream",
	ViewChannel:             "viewChannel",
	SendMessages:            "sendMessages",
	SendTTSMessages:         "sendTTSMessages",
	ManageMessages:          "manageMessages",
	EmbedLinks:              "embedLinks",
	AttachFiles:             "attachFiles",
	ReadMessageHistory:      "readMessageHistory",
	MentionEveryone:         "mentionEveryone",
	UseExternalEmojis:       "useExternalEmojis",
	ViewGuildInsights:       "viewGuildInsights",
	Connect:                 "connect",
	Speak:                   "speak",
	MuteMembers:             "muteMembers",
	DeafenMembers:           "deafenMembers",
	MoveMembers:             "moveMembers",
	UseVAD:                  "useVAD",
	ChangeNickname:          "changeNickname",
	ManageNicknames:         "manageNicknames",
	ManageRoles:             "manageRoles",
	ManageWebhooks:          "manageWebhooks",
	ManageEmojisAndStickers: "manageEmojisAndStickers",
	UseApplicationCommands:  "useApplicationCommands",
	RequestToSpeak:          "requestToSpeak",
	ManageEvents:            "manageEvents",
	ManageThreads:           "manageThreads",
	CreatePublicThreads:     "createPublicThreads",
	CreatePrivateThreads:    "createPrivateThreads",
	UseExternalStickers:     "useExternalStickers",
	SendMessagesInThreads:   "sendMessagesInThreads",
	UseEmbeddedActivities:   "useEmbeddedActivities",
	ModerateMembers:         "moderateMembers",
}

// String returns string representation of permission
func (p Permission) String() string {
	return permissions[p]
}
