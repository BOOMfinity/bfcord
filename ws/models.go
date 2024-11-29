package ws

import (
	"github.com/segmentio/encoding/json"
	"sync/atomic"

	"github.com/andersfylling/snowflake/v5"
)

type Event struct {
	OpCode     uint            `json:"op"`
	Data       json.RawMessage `json:"d"`
	Seq        uint64          `json:"s"`
	Event      string          `json:"t"`
	references atomic.Int64
}

func (ev *Event) free() {
	ev.Data = nil
	eventPool.Put(ev)
}

func (e *Event) reference() {
	e.references.Add(1)
}

func (e *Event) Dereference() {
	if e.references.Add(-1) <= 0 {
		e.free()
	}
}

type HelloOp struct {
	HeartbeatInterval uint `json:"heartbeat_interval"`
}

type resumeEvent struct {
	Token     string `json:"token"`
	SessionID string `json:"session_id"`
	Seq       uint64 `json:"seq"`
}

type sendEvent[T any] struct {
	OpCode uint `json:"op"`
	Data   T    `json:"d"`
}

type Identify struct {
	Token          string             `json:"token"`
	Properties     IdentifyProperties `json:"properties"`
	Compress       bool               `json:"compress,omitempty"`
	LargeThreshold int                `json:"large_threshold,omitempty"`
	Shard          []uint16           `json:"shard"`
	Intents        GatewayIntent      `json:"intents"`
}

type IdentifyProperties struct {
	OS      string `json:"os"`
	Browser string `json:"browser"`
	Device  string `json:"device"`
}

type UnavailableGuild struct {
	Name        string       `json:"name,omitempty"`
	ID          snowflake.ID `json:"id,omitempty"`
	Unavailable bool         `json:"unavailable,omitempty"`
}

type GatewayIntent uint

const (
	GatewayIntentGuilds                      GatewayIntent = 1 << 0
	GatewayIntentGuildMembers                GatewayIntent = 1 << 1
	GatewayIntentGuildModeration             GatewayIntent = 1 << 2
	GatewayIntentGuildEmojisAndStickers      GatewayIntent = 1 << 3
	GatewayIntentGuildIntegrations           GatewayIntent = 1 << 4
	GatewayIntentGuildWebhooks               GatewayIntent = 1 << 5
	GatewayIntentGuildInvites                GatewayIntent = 1 << 6
	GatewayIntentGuildVoiceStates            GatewayIntent = 1 << 7
	GatewayIntentGuildPresences              GatewayIntent = 1 << 8
	GatewayIntentGuildMessages               GatewayIntent = 1 << 9
	GatewayIntentGuildMessageReactions       GatewayIntent = 1 << 10
	GatewayIntentGuildMessageTyping          GatewayIntent = 1 << 11
	GatewayIntentGuildDirectMessages         GatewayIntent = 1 << 12
	GatewayIntentDirectMessageReactions      GatewayIntent = 1 << 13
	GatewayIntentDirectMessageTyping         GatewayIntent = 1 << 14
	GatewayIntentMessageContent              GatewayIntent = 1 << 15
	GatewayIntentGuildScheduledEvents        GatewayIntent = 1 << 16
	GatewayIntentAutoModerationConfiguration GatewayIntent = 1 << 20
	GatewayIntentAutoModerationExecution     GatewayIntent = 1 << 21
	GatewayIntentGuildMessagePolls           GatewayIntent = 1 << 24
	GatewayIntentDirectMessagePolls          GatewayIntent = 1 << 25

	GatewayIntentDefault = GatewayIntentGuilds | GatewayIntentGuildMembers | GatewayIntentGuildModeration | GatewayIntentGuildEmojisAndStickers | GatewayIntentGuildIntegrations | GatewayIntentGuildWebhooks | GatewayIntentGuildInvites | GatewayIntentGuildVoiceStates | GatewayIntentGuildMessages | GatewayIntentGuildMessageReactions | GatewayIntentGuildDirectMessages | GatewayIntentDirectMessageReactions | GatewayIntentMessageContent | GatewayIntentGuildScheduledEvents | GatewayIntentAutoModerationConfiguration | GatewayIntentAutoModerationExecution | GatewayIntentGuildMessagePolls | GatewayIntentDirectMessagePolls
)
