package discord

import (
	"github.com/BOOMfinity/bfcord/utils"
	"github.com/andersfylling/snowflake/v5"
)

type Poll struct {
	Question         PollMedia      `json:"question"`
	Answers          []PollAnswer   `json:"answers"`
	Expiry           Timestamp      `json:"expiry"`
	AllowMultiselect bool           `json:"allow_multiselect,omitempty"`
	LayoutType       PollLayoutType `json:"layout_type,omitempty"`
	Results          PollResult     `json:"results"`
}

type PollResult struct {
	IsFinalized  bool            `json:"is_finalized,omitempty"`
	AnswerCounts PollAnswerCount `json:"answer_counts"`
}

type PollAnswerCount struct {
	ID      snowflake.ID `json:"id,omitempty"`
	Count   uint         `json:"count,omitempty"`
	MeVoted bool         `json:"me_voted,omitempty"`
}

type PollLayoutType uint

const (
	PollLayoutTypeDefault = 1
)

type PollAnswer struct {
	AnswerID  snowflake.ID `json:"answer_id,omitempty"`
	PollMedia PollMedia    `json:"poll_media"`
}

type PollMedia struct {
	Text  string                `json:"text,omitempty"`
	Emoji utils.Nullable[Emoji] `json:"emoji,omitempty"`
}
