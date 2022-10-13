package discord

import (
	"github.com/BOOMfinity/bfcord/internal/timeconv"
	"github.com/andersfylling/snowflake/v5"
)

// Presence
//
// Reference: https://discord.com/developers/docs/topics/gateway#presence-update-presence-update-event-fields
type Presence struct {
	BasePresence
	User User `json:"user"`
}

type BasePresence struct {
	ClientStatus PresenceClientStatus `json:"client_status"`
	Status       PresenceStatus       `json:"status"`
	Activities   []Activity           `json:"activities"`
	GuildID      snowflake.ID         `json:"guild_id"`
	UserID       snowflake.ID         `json:"user_id"`
}

type PresenceClientStatus struct {
	Desktop PresenceStatus `json:"desktop"`
	Mobile  PresenceStatus `json:"mobile"`
	Web     PresenceStatus `json:"web"`
}

// Activity
//
// Reference: https://discord.com/developers/docs/topics/gateway#activity-object-activity-structure
type Activity struct {
	Emoji         *Emoji          `json:"emoji,omitempty"`
	Party         *ActivityParty  `json:"party,omitempty"`
	Assets        *ActivityAssets `json:"assets,omitempty"`
	Name          string          `json:"name"`
	Details       string          `json:"details,omitempty"`
	State         string          `json:"state,omitempty"`
	Url           string          `json:"url,omitempty"`
	ApplicationID snowflake.ID    `json:"application_id,omitempty"`
	CreatedAt     int64           `json:"created_at,omitempty"`
	Flags         ActivityFlags   `json:"flags,omitempty"`
	Type          ActivityType    `json:"type"`
	Instance      bool            `json:"instance,omitempty"`
}

type ActivityButton struct {
	Label string `json:"label"`
	Url   string `json:"url"`
}

type ActivityFlags uint16

const (
	ActivityFlagInstance ActivityFlags = 1 << iota
	ActivityFlagJoin
	ActivityFlagSpectate
	ActivityFlagJoinRequest
	ActivityFlagSync
	ActivityFlagPlay
	ActivityFlagPartyFriends
	ActivityFlagPartyVoice
	ActivityFlagEmbedded
)

type ActivityAssets struct {
	LargeImage string `json:"large_image"`
	LargeText  string `json:"large_text"`
	SmallImage string `json:"small_image"`
	SmallText  string `json:"small_text"`
}

type ActivityParty struct {
	ID   string `json:"id"`
	Size []int  `json:"size"`
}

type ActivityTimestamps struct {
	Start timeconv.Timestamp `json:"start"`
	End   timeconv.Timestamp `json:"end"`
}

type ActivityType uint8

const (
	ActivityGame = iota
	ActivityStreaming
	ActivityListening
	ActivityWatching
	ActivityCustom
	ActivityCompeting
)

type PresenceStatus string

const (
	StatusOnline  = "online"
	StatusIdle    = "idle"
	StatusDND     = "dnd"
	StatusOffline = "offline"
)
