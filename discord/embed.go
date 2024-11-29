package discord

import "github.com/BOOMfinity/bfcord/utils"

type MessageEmbed struct {
	Title       string                               `json:"title,omitempty"`
	Type        string                               `json:"type,omitempty"`
	Description string                               `json:"description,omitempty"`
	URL         string                               `json:"url,omitempty"`
	Timestamp   Timestamp                            `json:"timestamp,omitempty"`
	Color       int                                  `json:"color,omitempty"`
	Footer      utils.Nullable[MessageEmbedFooter]   `json:"footer,omitempty"`
	Fields      []MessageEmbedField                  `json:"fields,omitempty"`
	Image       utils.Nullable[MessageEmbedMedia]    `json:"image,omitempty"`
	Thumbnail   utils.Nullable[MessageEmbedMedia]    `json:"thumbnail,omitempty"`
	Video       utils.Nullable[MessageEmbedMedia]    `json:"video,omitempty"`
	Author      utils.Nullable[MessageEmbedAuthor]   `json:"author,omitempty"`
	Provider    utils.Nullable[MessageEmbedProvider] `json:"provider,omitempty"`
}

type MessageEmbedProvider struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

type MessageEmbedAuthor struct {
	Name         string `json:"name,omitempty"`
	URL          string `json:"url,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

type MessageEmbedFooter struct {
	Text         string `json:"text,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

type MessageEmbedField struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty"`
}

type MessageEmbedMedia struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Height   uint   `json:"height,omitempty"`
	Width    uint   `json:"width,omitempty"`
}
