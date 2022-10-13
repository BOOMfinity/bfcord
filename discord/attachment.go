package discord

import "github.com/andersfylling/snowflake/v5"

type Attachment struct {
	ProxyUrl    string       `json:"proxy_url,omitempty"`
	Filename    string       `json:"filename,omitempty"`
	Description string       `json:"description,omitempty"`
	ContentType string       `json:"content_type,omitempty"`
	Url         string       `json:"url,omitempty"`
	Size        int          `json:"size,omitempty"`
	ID          snowflake.ID `json:"id"`
	Height      int          `json:"height,omitempty"`
	Width       int          `json:"width,omitempty"`
	Ephemeral   bool         `json:"ephemeral,omitempty"`
}
