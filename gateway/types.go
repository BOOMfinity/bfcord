package gateway

import (
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/bfcord/gateway/events"
	"github.com/segmentio/encoding/json"
)

type OpCode uint8

const (
	DispatchOp OpCode = iota
	HeartbeatOp
	IdentifyOp
	PresenceUpdateOp
	VoiceStateUpdateOp
	ResumeOp OpCode = iota + 1
	ReconnectOp
	RequestGuildMembersOp
	InvalidSessionOp
	HelloOp
	HeartbeatAckOp
)

type PresenceUpdate struct {
	Since      *uint64                `json:"since"`
	Status     discord.PresenceStatus `json:"status"`
	Activities []discord.Activity     `json:"activities"`
	AFK        bool                   `json:"afk"`
}

type Payload struct {
	Event events.Event    `json:"t"`
	Data  json.RawMessage `json:"d"`
	Seq   uint64          `json:"s"`
	Shard uint16          `json:"-"`
	OP    OpCode          `json:"op"`
}

type helloData struct {
	Interval int `json:"heartbeat_interval"`
}
