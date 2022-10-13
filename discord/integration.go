package discord

import "github.com/andersfylling/snowflake/v5"

type Integration struct {
	Account struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"account"`
	Type              string                    `json:"type"`
	SyncedAt          string                    `json:"synced_at"`
	Name              string                    `json:"name"`
	User              User                      `json:"user"`
	GuildID           snowflake.ID              `json:"guild_id"`
	SubscriberCount   int                       `json:"subscriber_count"`
	RoleID            snowflake.ID              `json:"role_id"`
	ID                snowflake.ID              `json:"id"`
	ExpireGracePeriod int                       `json:"expire_grace_period"`
	EnableEmoticons   bool                      `json:"enable_emoticons"`
	ExpireBehavior    IntegrationExpireBehavior `json:"expire_behavior"`
	Syncing           bool                      `json:"syncing"`
	Enabled           bool                      `json:"enabled"`
	Revoked           bool                      `json:"revoked"`
}

type IntegrationExpireBehavior uint8

const (
	ExpireBehaviorRemoveRole IntegrationExpireBehavior = iota
	ExpireBehaviorKickUser
)
