package discord

import (
	"github.com/andersfylling/snowflake/v5"
)

type User struct {
	ID            snowflake.ID `json:"id,omitempty"`
	Username      string       `json:"username,omitempty"`
	Discriminator string       `json:"discriminator,omitempty"`
	GlobalName    string       `json:"global_name,omitempty"`
	Avatar        string       `json:"avatar,omitempty"`
	Bot           bool         `json:"bot,omitempty"`
	System        bool         `json:"system,omitempty"`
	MFAEnabled    bool         `json:"mfa_enabled,omitempty"`
	Banner        string       `json:"banner,omitempty"`
	AccentColor   uint         `json:"accent_color,omitempty"`
	Locale        string       `json:"locale,omitempty"`
	Verified      bool         `json:"verified,omitempty"`
}

func (u User) Partial() bool {
	return u.Username == ""
}
