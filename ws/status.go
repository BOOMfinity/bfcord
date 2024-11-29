package ws

import (
	"github.com/BOOMfinity/go-utils/rate"
	"github.com/BOOMfinity/golog/v2"
)

type Status string

const (
	StatusConnected    Status = "connected"
	StatusConnecting   Status = "connecting"
	StatusReconnecting Status = "reconnecting"
	StatusDisconnected Status = "disconnected"
)

type Config struct {
	URL           string
	ID            uint16
	ShardCount    uint16
	Logger        golog.Logger
	Compression   bool
	Token         string
	Intents       GatewayIntent
	GlobalLimiter *rate.Limiter
}
