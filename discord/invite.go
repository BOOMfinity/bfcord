package discord

import "github.com/BOOMfinity/bfcord/internal/timeconv"

type InviteWithMeta struct {
	CreatedAt timeconv.Timestamp `json:"created_at,omitempty"`
	Invite
	Uses      uint32 `json:"uses,omitempty"`
	MaxUses   uint32 `json:"max_uses,omitempty"`
	MaxAge    uint32 `json:"max_age,omitempty"`
	Temporary bool   `json:"temporary,omitempty"`
}

type Invite struct {
	ExpiresAt           timeconv.Timestamp   `json:"expires_at,omitempty"`
	Guild               *BaseGuild           `json:"guild,omitempty"`
	Inviter             *User                `json:"inviter,omitempty"`
	GuildScheduledEvent *GuildScheduledEvent `json:"guild_scheduled_event,omitempty"`
	TargetUser          *User                `json:"target_user,omitempty"`
	Channel             *Channel             `json:"channel,omitempty"`
	StageInstance       *StageInstance       `json:"stage_instance,omitempty"`
	Code                string               `json:"code,omitempty"`
	PresenceCount       uint32               `json:"approximate_presence_count,omitempty"`
	MemberCount         uint32               `json:"approximate_member_count,omitempty"`
	TargetType          InviteTargetType     `json:"target_type,omitempty"`
}

type InviteTargetType uint8

const (
	TargetTypeStream InviteTargetType = iota + 1
	TargetTypeEmbeddedApplication
)
