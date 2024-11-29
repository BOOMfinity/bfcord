package discord

import (
	"github.com/andersfylling/snowflake/v5"
)

type Reaction struct {
	Count        uint                 `json:"count,omitempty"`
	CountDetails ReactionCountDetails `json:"count_details"`
	Me           bool                 `json:"me,omitempty"`
	MeBurst      bool                 `json:"me_burst,omitempty"`
	Emoji        Emoji                `json:"emoji,omitempty"`
	BurstColors  []string             `json:"burst_colors,omitempty"`
}

type ReactionCountDetails struct {
	Burst  uint `json:"burst,omitempty"`
	Normal uint `json:"normal,omitempty"`
}

type Emoji struct {
	ID            snowflake.ID   `json:"id,omitempty"`
	Name          string         `json:"name,omitempty"`
	Roles         []snowflake.ID `json:"roles,omitempty"`
	User          User           `json:"user,omitempty"`
	RequireColons bool           `json:"require_colons,omitempty"`
	Managed       bool           `json:"managed,omitempty"`
	Animated      bool           `json:"animated,omitempty"`
	Available     bool           `json:"available,omitempty"`
}

type ReactionType uint

const (
	ReactionNormal ReactionType = iota
	ReactionBurst
)
