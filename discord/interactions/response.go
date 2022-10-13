package interactions

import (
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/bfcord/slash"
)

type CallbackType uint8

const (
	PongCallback                     CallbackType = iota + 1
	ChannelMessageWithSourceCallback CallbackType = iota + 3
	DeferredChannelMessageWithSourceCallback
	DeferredUpdateMessageCallback
	UpdateMessageCallback
	AutocompleteResultCallback
	ModalCallback
)

type InteractionCallbackFlags uint8

const (
	EphemeralFlag InteractionCallbackFlags = 1 << 6
)

type ResponseData struct {
	discord.MessageCreate
	//AllowedMentions discord.MessageAllowedMentions `json:"allowed_mentions,omitempty"`
	Choices *[]slash.Choice          `json:"choices,omitempty"`
	Flags   InteractionCallbackFlags `json:"flags,omitempty"`
}

type InteractionResponse struct {
	Data *ResponseData `json:"data,omitempty"`
	Type CallbackType  `json:"type"`
}
