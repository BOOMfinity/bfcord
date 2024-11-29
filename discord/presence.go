package discord

import (
	"github.com/BOOMfinity/bfcord/utils"
	"github.com/andersfylling/snowflake/v5"
)

type Presence struct {
	User         User         `json:"user,omitempty"`
	GuildID      snowflake.ID `json:"guild_id,omitempty"`
	Status       UserStatus   `json:"status,omitempty"`
	Activities   []Activity   `json:"activities,omitempty"`
	ClientStatus ClientStatus `json:"client_status"`
}

type ClientStatus struct {
	OS      string `json:"os,omitempty"`
	Browser string `json:"browser,omitempty"`
	Device  string `json:"device,omitempty"`
}

type UserStatus string

const (
	UserStatusOnline    UserStatus = "online"
	UserStatusIdle      UserStatus = "idle"
	UserStatusDoNotDist UserStatus = "dnd"
	UserStatusOffline   UserStatus = "offline"
)

type Activity struct {
	Name          string                          `json:"name,omitempty"`
	Type          ActivityType                    `json:"type,omitempty"`
	URL           string                          `json:"url,omitempty"`
	CreatedAt     UnixTimestamp                   `json:"created_at"`
	Timestamps    ActivityTimestamps              `json:"timestamps"`
	ApplicationID snowflake.ID                    `json:"application_id,omitempty"`
	Details       string                          `json:"details,omitempty"`
	State         string                          `json:"state,omitempty"`
	Emoji         utils.Nullable[Emoji]           `json:"emoji,omitempty"`
	Party         utils.Nullable[ActivityParty]   `json:"party,omitempty"`
	Assets        utils.Nullable[ActivityAssets]  `json:"assets,omitempty"`
	Secrets       utils.Nullable[ActivitySecrets] `json:"secrets,omitempty"`
	Instance      bool                            `json:"instance,omitempty"`
	Flags         ActivityFlag                    `json:"flags,omitempty"`
	Buttons       []string                        `json:"buttons,omitempty"`
}

type ActivityButton struct {
	Label string `json:"label,omitempty"`
	URL   string `json:"url,omitempty"`
}

type ActivityFlag uint16

const (
	ActivityFlagInstance ActivityFlag = 1 << iota
	ActivityFlagJoin
	ActivityFlagSpectate
	ActivityFlagJoinRequest
	ActivityFlagSync
	ActivityFlagPlay
	ActivityFlagPartyPrivacyFriends
	ActivityFlagPartyPrivacyVoiceChannel
	ActivityFlagEmbedded
)

type ActivitySecrets struct {
	Join     string `json:"join,omitempty"`
	Spectate string `json:"spectate,omitempty"`
	Match    string `json:"match,omitempty"`
}

type ActivityAssets struct {
	LargeImage string `json:"large_image,omitempty"`
	LargeText  string `json:"large_text,omitempty"`
	SmallImage string `json:"small_image,omitempty"`
	SmallText  string `json:"small_text,omitempty"`
}

type ActivityParty struct {
	ID   string `json:"id,omitempty"`
	Size [2]int `json:"size,omitempty"`
}

type ActivityTimestamps struct {
	Start UnixTimestamp `json:"start"`
	End   UnixTimestamp `json:"end"`
}

type ActivityType uint8

const (
	ActivityTypeGame ActivityType = iota
	ActivityTypeStreaming
	ActivityTypeListening
	ActivityTypeWatching
	ActivityTypeCustom
	ActivityTypeCompeting
)
