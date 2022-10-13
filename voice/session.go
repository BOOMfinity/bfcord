package voice

import (
	"context"
	gateway2 "github.com/BOOMfinity/bfcord/gateway"
	"github.com/BOOMfinity/go-utils/broadcaster"
	"github.com/BOOMfinity/golog"
	"github.com/BOOMfinity/wshelper"
	"time"
)

type Session struct {
	config           *ConnectOptions
	state            *gateway2.VoiceStateUpdateEvent
	server           *gateway2.VoiceServerUpdateEvent
	ws               *wshelper.Connection
	logger           golog.Logger
	eventChannel     *broadcaster.Group[OPCode] // code 111: disconnection, 101: connection done
	udp              *udpConnection
	ready            *ReadyPayload
	reconnectAttempt uint8
}

// NewSession function creates new session and issues first connection. Blocks until connection is finished (or until timeout occurs)
func NewSession(ctx context.Context, opts ConnectOptions) (*Session, error) {
	conn := Session{logger: golog.New("voice").Module(opts.GuildID.String()), eventChannel: broadcaster.NewGroup[OPCode](), config: &opts}
	if opts.Debug {
		conn.logger.SetLevel(golog.LevelDebug)
	}
	conn.logger.Info().Send("Connection requested")

	return &conn, conn.connect(ctx, &opts)
}

// connect method requests connection from main gateway, opens voice gateway connection and blocks until further connection steps are finished. Throws an error after 10s timeout.
func (s *Session) connect(ctx context.Context, opts *ConnectOptions) error {
	start := time.Now()
	events := s.eventChannel.Join() // upon successful connection 101 code is sent to this channel
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer func() {
		cancel()
		events.Close()
	}()

	s.server = opts.VoiceServer
	s.state = opts.VoiceState

	err := s.openGateway(ctx)
	if err != nil {
		return err
	}

wait:
	for {
		select {
		case op := <-events.Out:
			if op.Data() == 101 {
				break wait
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	s.logger.Info().Send("Connected in %vms", time.Now().Sub(start).Milliseconds())
	return s.SendSpeaking(Microphone)
}

// Destroy method disconnects session from discord. Using SendOpusFrame after destroying will throw errors.
func (s *Session) Destroy() {
	_ = s.ws.Close(1000, "bfcord-voice: destroy connection")
}

// internal destroy
func (s *Session) destroy(errored bool) {
	s.udp.Close()

	if errored {
		s.logger.Warn().Send("Connection closed due to error")
	} else {
		s.logger.Info().Send("Connection closed")
	}

	s.config.OnClose()
}

// SendOpusFrame method writes raw audio frame to discord.
//
// Frame MUST be 20ms in length. We may allow changing the value later if needed.
//
// This method should not be used concurrently, as it will cause audio overlapping.
//
// Will throw ConnectionClosedError if connection is closed permanently.
func (s *Session) SendOpusFrame(frame []byte) error {
	if s.udp.isClosed.Load() {
		return ConnectionClosedError
	}

	return s.udp.WriteOpusFrame(frame)
}

// IsClosed will return true if session is fully closed.
func (s *Session) IsClosed() bool {
	return s.udp.isClosed.Load()
}
