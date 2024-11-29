package discord

import "github.com/andersfylling/snowflake/v5"

type ScheduledEvent struct {
	ID                 snowflake.ID               `json:"id,omitempty"`
	GuildID            snowflake.ID               `json:"guild_id,omitempty"`
	ChannelID          snowflake.ID               `json:"channel_id,omitempty"`
	CreatorID          snowflake.ID               `json:"creator_id,omitempty"`
	Name               string                     `json:"name,omitempty"`
	Description        string                     `json:"description,omitempty"`
	ScheduledStartTime Timestamp                  `json:"scheduled_start_time,omitempty"`
	ScheduledEndTime   Timestamp                  `json:"scheduled_end_time,omitempty"`
	PrivacyLevel       ScheduledEventPrivacyLevel `json:"privacy_level,omitempty"`
	Status             ScheduledEventStatus       `json:"status,omitempty"`
	EntityType         ScheduledEventEntityType   `json:"entity_type,omitempty"`
	EntityID           snowflake.ID               `json:"entity_id,omitempty"`
	EntityMetadata     ScheduledEventEntity       `json:"entity_metadata,omitempty"`
	Creator            User                       `json:"creator,omitempty"`
	UserCount          uint                       `json:"user_count,omitempty"`
	Image              string                     `json:"image,omitempty"`
}

type ScheduledEventPrivacyLevel uint

const (
	ScheduledEventGuildOnly ScheduledEventPrivacyLevel = iota + 2
)

type ScheduledEventStatus uint

const (
	ScheduledEventStatusScheduled ScheduledEventStatus = iota + 1
	ScheduledEventStatusActive
	ScheduledEventStatusCompleted
	ScheduledEventStatusCanceled
)

type ScheduledEventEntityType uint

const (
	ScheduledEventEntityTypeStageInstance ScheduledEventEntityType = iota + 1
	ScheduledEventEntityTypeVoice
	ScheduledEventEntityTypeExternal
)

type ScheduledEventEntity struct {
	Location string `json:"location,omitempty"`
}
