package voice

import (
	"github.com/BOOMfinity/bfcord/gateway"
	"github.com/andersfylling/snowflake/v5"
	"github.com/segmentio/encoding/json"
	"github.com/unxcepted/websocket"
)

/*
Own eventChannel codes:
101 - connection finished - sent when session description OR resumed is received
111 - shutdown - shutdowns heartbeat loop and websocket connection
*/

type ConnectOptions struct {
	VoiceState  *gateway.VoiceStateUpdateEvent
	VoiceServer *gateway.VoiceServerUpdateEvent
	ChannelID   snowflake.ID
	GuildID     snowflake.ID
	Debug       bool
	OnClose     func()
}

type SpeakingFlag uint64

const (
	NotSpeaking SpeakingFlag = 0
	Microphone  SpeakingFlag = 1 << iota
	SoundShare
	Priority
)

type OPCode uint8

const (
	identifyOP OPCode = iota
	selectProtocolOP
	readyOP
	heartbeatOP
	sessionDescriptionOP
	speakingOP
	heartbeatAckOP
	resumeOP
	helloOP
	resumedOP
	clientDisconnectOP
)

const ( // known close codes
	// Discord normal
	SessionNoLongerValid websocket.StatusCode = 4006
	SessionTimeout       websocket.StatusCode = 4009
	Disconnected         websocket.StatusCode = 4014
	ServerCrashed        websocket.StatusCode = 4015

	// Discord abnormal
	UnknownOPCode         websocket.StatusCode = 4001
	FailedToDecodePayload websocket.StatusCode = 4002
	NotAuthenticated      websocket.StatusCode = 4003
	AuthenticationFailed  websocket.StatusCode = 4004
	AlreadyAuthenticated  websocket.StatusCode = 4005
	ServerNotFound        websocket.StatusCode = 4011
	UnknownProtocol       websocket.StatusCode = 4012
	UnknownEncryption     websocket.StatusCode = 4016

	BfcordVoiceReconnect websocket.StatusCode = 4969
)

type Payload struct {
	Data json.RawMessage `json:"d"`
	Op   OPCode          `json:"op"`
}

func payloadTo[V any](p *Payload) (res V, err error) {
	err = json.Unmarshal(p.Data, &res)
	return
}

type IdentifyPayload struct {
	SessionID string       `json:"session_id"`
	Token     string       `json:"token"`
	ServerID  snowflake.ID `json:"server_id"`
	UserID    snowflake.ID `json:"user_id"`
}

type ReadyPayload struct {
	IP    string   `json:"IP"`
	Modes []string `json:"modes"`
	SSRC  uint32   `json:"ssrc"`
	Port  uint16   `json:"port"`
	// heartbeat_interval here is an erroneous field and should be ignored. The correct heartbeat_interval value comes from the Hello payload.
}

type helloPayload struct {
	HeartbeatInterval float32 `json:"heartbeat_interval"`
}

type selectProtocolPayload struct {
	Protocol string             `json:"protocol"` // always UDP
	Data     selectProtocolData `json:"data"`
}

type selectProtocolData struct {
	Address        string `json:"address"`
	EncryptionMode string `json:"mode"`
	Port           uint16 `json:"port"`
}

type sessionDescriptionPayload struct {
	EncryptionMode string   `json:"mode"`
	SecretKey      [32]byte `json:"secret_key"`
}

type speakingPayload struct {
	Speaking SpeakingFlag `json:"speaking"`
	Delay    int          `json:"delay"`
	SSRC     uint32       `json:"ssrc"`
}

type resumePayload struct {
	SessionID string       `json:"session_id"`
	Token     string       `json:"token"`
	ServerID  snowflake.ID `json:"server_id"`
}
