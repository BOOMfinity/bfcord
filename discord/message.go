package discord

import (
	"github.com/andersfylling/snowflake/v5"

	"github.com/BOOMfinity/bfcord/utils"
)

type Message struct {
	ID                  snowflake.ID                               `json:"id,omitempty"`
	ChannelID           snowflake.ID                               `json:"channel_id,omitempty"`
	Author              User                                       `json:"author,omitempty"`
	Content             string                                     `json:"content,omitempty"`
	Timestamp           Timestamp                                  `json:"timestamp"`
	EditedTimestamp     Timestamp                                  `json:"edited_timestamp"`
	TTS                 bool                                       `json:"tts,omitempty"`
	MentionEveryone     bool                                       `json:"mention_everyone,omitempty"`
	Mentions            []User                                     `json:"mentions,omitempty"`
	MentionRoles        []snowflake.ID                             `json:"mention_roles,omitempty"`
	MentionChannels     []MessageChannelMention                    `json:"mention_channels,omitempty"`
	Attachments         []Attachment                               `json:"attachments,omitempty"`
	Embeds              []MessageEmbed                             `json:"embeds,omitempty"`
	Reactions           []Reaction                                 `json:"reactions,omitempty"`
	Nonce               string                                     `json:"nonce,omitempty"`
	Pinned              bool                                       `json:"pinned,omitempty"`
	WebhookID           snowflake.ID                               `json:"webhook_id,omitempty"`
	Type                MessageType                                `json:"type,omitempty"`
	Activity            utils.Nullable[MessageActivity]            `json:"activity,omitempty"`
	Application         utils.Nullable[Application]                `json:"application,omitempty"`
	ApplicationID       snowflake.ID                               `json:"application_id,omitempty"`
	MessageReference    utils.Nullable[MessageReference]           `json:"message_reference,omitempty"`
	Flags               MessageFlag                                `json:"flags,omitempty"`
	InteractionMetadata utils.Nullable[MessageInteractionMetadata] `json:"interaction_metadata,omitempty"`
	Thread              utils.Nullable[Channel]                    `json:"thread,omitempty"`
	Components          ActionRows                                 `json:"components,omitempty"`
	Position            int                                        `json:"position,omitempty"`
	Resolved            ResolvedData                               `json:"resolved,omitempty"`
	Poll                utils.Nullable[Poll]                       `json:"poll,omitempty"`
}

type MessageChannelMention struct {
	ID      snowflake.ID `json:"id,omitempty"`
	GuildID snowflake.ID `json:"guild_id,omitempty"`
	Type    ChannelType  `json:"type,omitempty"`
	Name    string       `json:"name,omitempty"`
}

type MessageInteractionMetadata struct {
	ID                            snowflake.ID                `json:"id,omitempty"`
	Type                          InteractionType             `json:"type,omitempty"`
	User                          User                        `json:"user"`
	OriginalResponseMessageID     snowflake.ID                `json:"original_response_message_id,omitempty"`
	InteractedMessageID           snowflake.ID                `json:"interacted_message_id,omitempty"`
	TriggeringInteractionMetadata *MessageInteractionMetadata `json:"triggering_interaction_metadata,omitempty"`
	// AuthorizingIntegrationOwners []ApplicationIntegrationType
}

type MessageFlag = BitField

const (
	MessageFlagCrossPosted                      MessageFlag = 1 << 0
	MessageFlagIsCrossPost                      MessageFlag = 1 << 1
	MessageFlagSuppressEmbeds                   MessageFlag = 1 << 2
	MessageFlagSourceMessageDeleted             MessageFlag = 1 << 3
	MessageFlagUrgent                           MessageFlag = 1 << 4
	MessageFlagHasThread                        MessageFlag = 1 << 5
	MessageFlagEphemeral                        MessageFlag = 1 << 6
	MessageFlagFailedToMentionSomeRolesInThread MessageFlag = 1 << 7
	MessageFlagSuppressNotifications            MessageFlag = 1 << 12
	MessageFlagIsVoiceMessage                   MessageFlag = 1 << 13
)

type MessageReference struct {
	ChannelID       snowflake.ID `json:"channel_id,omitempty"`
	GuildID         snowflake.ID `json:"guild_id,omitempty"`
	MessageID       snowflake.ID `json:"message_id,omitempty"`
	FailIfNotExists bool         `json:"fail_if_not_exists"`
}

type MessageActivity struct {
	PartyID string              `json:"party_id,omitempty"`
	Type    MessageActivityType `json:"type,omitempty"`
}

type MessageActivityType uint8

const (
	MessageActivityTypeJoin MessageActivityType = iota + 1
	MessageActivityTypeSpectate
	MessageActivityTypeListen
	MessageActivityTypeJoinRequest MessageActivityType = iota + 2
)

type MessageType uint8

const (
	MessageTypeDefault                                 MessageType = 0
	MessageTypeRecipientAdd                            MessageType = 1
	MessageTypeRecipientRemove                         MessageType = 2
	MessageTypeCall                                    MessageType = 3
	MessageTypeChannelNameChange                       MessageType = 4
	MessageTypeChannelIconChange                       MessageType = 5
	MessageTypeChannelPinnedMessage                    MessageType = 6
	MessageTypeGuildMemberJoin                         MessageType = 7
	MessageTypeUserPremiumGuildSubscription            MessageType = 8
	MessageTypeUserPremiumGuildSubscriptionTier1       MessageType = 9
	MessageTypeUserPremiumGuildSubscriptionTier2       MessageType = 10
	MessageTypeUserPremiumGuildSubscriptionTier3       MessageType = 11
	MessageTypeChannelFollowAdd                        MessageType = 12
	MessageTypeGuildDiscoveryDisqualified              MessageType = 14
	MessageTypeGuildDiscoveryReQualified               MessageType = 15
	MessageTypeGuildDiscoveryGracePeriodInitialWarning MessageType = 16
	MessageTypeGuildDiscoveryGracePeriodFinalWarning   MessageType = 17
	MessageTypeThreadCreated                           MessageType = 18
	MessageTypeReply                                   MessageType = 19
	MessageTypeChatInputCommand                        MessageType = 20
	MessageTypeThreadStarterMessage                    MessageType = 21
	MessageTypeGuildInviteReminder                     MessageType = 22
	MessageTypeContextMenuCommand                      MessageType = 23
	MessageTypeAutoModerationAction                    MessageType = 24
	MessageTypeRoleSubscriptionPurchase                MessageType = 25
	MessageTypeInteractionPremiumUpsell                MessageType = 26
	MessageTypeStageStart                              MessageType = 27
	MessageTypeStageEnd                                MessageType = 28
	MessageTypeStageSpeaker                            MessageType = 29
	MessageTypeStageTopic                              MessageType = 31
	MessageTypeGuildApplicationPremiumSubscription     MessageType = 32
)

type Attachment struct {
	ID           snowflake.ID   `json:"id,omitempty"`
	FileName     string         `json:"file_name,omitempty"`
	Description  string         `json:"description,omitempty"`
	ContentType  string         `json:"content_type,omitempty"`
	Size         uint           `json:"size,omitempty"`
	Url          string         `json:"url,omitempty"`
	ProxyUrl     string         `json:"proxy_url,omitempty"`
	Height       uint           `json:"height,omitempty"`
	Width        uint           `json:"width,omitempty"`
	Ephemeral    bool           `json:"ephemeral,omitempty"`
	DurationSecs float64        `json:"duration_secs,omitempty"`
	Waveform     string         `json:"waveform,omitempty"`
	Flags        AttachmentFlag `json:"flags,omitempty"`
}

type AttachmentFlag uint8

const (
	AttachmentFlagRemix AttachmentFlag = 1 << 2
)
