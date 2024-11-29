package discord

import "github.com/BOOMfinity/go-utils/nullable"

type Invite struct {
	InviteType               InviteType                        `json:"type,omitempty"`
	Code                     string                            `json:"code,omitempty"`
	Guild                    Guild                             `json:"guild"`
	Channel                  Channel                           `json:"channel"`
	Inviter                  User                              `json:"inviter"`
	TargetType               InviteTargetType                  `json:"target_type,omitempty"`
	TargetUser               User                              `json:"target_user"`
	TargetApplication        Application                       `json:"target_application"`
	ApproximatePresenceCount int                               `json:"approximate_presence_count,omitempty"`
	ApproximateMemberCount   int                               `json:"approximate_member_count,omitempty"`
	ExpiresAt                nullable.Nullable[Timestamp]      `json:"expires_at"`
	GuildScheduledEvent      nullable.Nullable[ScheduledEvent] `json:"guild_scheduled_event"`
}

type InviteType uint

const (
	InviteGuild InviteType = iota
	InviteGroupDM
	InviteFriend
)

type InviteTargetType uint

const (
	InviteTargetStream InviteTargetType = iota + 1
	InviteTargetEmbeddedApplication
)
