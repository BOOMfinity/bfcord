package ws

import (
	"time"
)

type InternalEventAllocator interface {
	reference()
	Dereference()
}

type InternalStatusChangeEvent Status

// InternalConnectionClosed is sent when Gateway Status changes to StatusDisconnected
type InternalConnectionClosed struct{}

type InternalHeartbeatEvent struct {
	Start time.Time
	End   time.Time
}

type InternalDispatchEvent = *Event

type InternalReadyEvent = *ReadyEvent

type InternalMaxReconnectionLimitReached struct{}
