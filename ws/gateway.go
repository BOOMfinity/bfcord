package ws

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/BOOMfinity/go-utils/gpool"
	"github.com/BOOMfinity/go-utils/inlineif"
	"github.com/BOOMfinity/golog/v2"
	"github.com/gorilla/websocket"
	"github.com/segmentio/encoding/json"

	"github.com/BOOMfinity/bfcord"
	"github.com/BOOMfinity/bfcord/discord"
)

const maxReconnections = 3

var eventPool = gpool.New[Event]()
var reconnectCodes = []int{4000, 4001, 4002, 4003, 4005, 4007, 4008, 4009}

type Gateway interface {
	Listen() (events <-chan any, cancel func())
	Disconnect()
	Connect(ctx context.Context) error
	Config() Config
	Status() Status
	Log() golog.Logger
	FetchMembers(ctx context.Context, params RequestGuildMembersParams) ([]discord.MemberWithUser, []discord.Presence, error)
}

type gatewayImpl struct {
	conn          *websocket.Conn
	cfg           Config
	session       string
	resumeURL     string
	seq           *atomic.Uint64
	events        []chan<- any
	mut           sync.RWMutex
	log           golog.Logger
	status        Status
	buff          *bytes.Buffer
	reconnections *atomic.Uint64
}

func (g *gatewayImpl) Log() golog.Logger {
	return g.log
}

func (g *gatewayImpl) Status() Status {
	g.mut.RLock()
	status := g.status
	g.mut.RUnlock()
	return status
}

func (g *gatewayImpl) sendEvent(ev any) {
	g.mut.RLock()
	for _, ch := range g.events {
		if allocator, ok := ev.(InternalEventAllocator); ok {
			allocator.reference()
		}
		go func() {
			ch <- ev
		}()
	}
	g.mut.RUnlock()
}

func (g *gatewayImpl) Listen() (events <-chan any, cancel func()) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		g.log.Trace().Send("Registering new event listener (%s:%d)", file, line)
	}
	ch := make(chan any, 1)
	idx := -1

	g.mut.Lock()
	for i, v := range g.events {
		if v == nil {
			idx = i
			break
		}
	}

	if idx == -1 {
		idx = len(g.events)
		g.events = append(g.events, ch)
	} else {
		g.events[idx] = ch
	}
	g.mut.Unlock()

	return ch, func() {
		if idx == -1 {
			return
		}
		g.mut.Lock()
		defer g.mut.Unlock()
		if len(g.events)-1 < idx {
			return
		}
		if obj := g.events[idx]; obj != nil {
			g.log.Trace().Send("Unregistering event listener (%s:%d)", file, line)
			g.events[idx] = nil
		}
	}
}

func (g *gatewayImpl) changeStatus(status Status) {
	g.sendEvent(status)
	g.mut.Lock()
	g.status = status
	g.mut.Unlock()
	g.log.Trace().Send("Status changed to '%s'", status)
}

func (g *gatewayImpl) Disconnect() {
	g.log.Warn().Send("Disconnect command received, closing the connection in FORCE mode")
	g.disconnect(true, false)
}

func (g *gatewayImpl) disconnect(reset bool, reconnect bool) {
	if g.status == StatusDisconnected {
		return
	}
	g.sendEvent(InternalConnectionClosed{})
	g.log.Trace().Param("can-resume", !reset).Param("reconnect", reconnect).Send("Closing connection")
	g.changeStatus(StatusDisconnected)
	if g.conn != nil {
		g.log.Trace().Send("Sending close frame as connection is not nil")
		_ = g.conn.Close()
	}
	g.conn = nil
	if reset {
		g.reset()
	}
	g.buff.Reset()

	if reconnect {
		if g.reconnections.Load() >= maxReconnections {
			g.log.Warn().Send("Max reconnections reached, next reconnection will be in 5 minutes")
			g.sendEvent(InternalMaxReconnectionLimitReached{})
			time.Sleep(5 * time.Minute)
		}
		go func() {
			g.reconnections.Add(1)
			if err := g.Connect(context.Background()); err != nil {
				g.log.Error().Throw(fmt.Errorf("could not reconnect: %w", err))
			}
		}()
	}
}

func (g *gatewayImpl) reset() {
	g.log.Trace().Send("Resetting session data")
	g.resumeURL = ""
	g.session = ""
	g.seq = &atomic.Uint64{}
}

func (g *gatewayImpl) read() (*Event, error) {
	g.buff.Reset()
	messageType, reader, err := g.conn.NextReader()
	if err != nil {
		var closed *websocket.CloseError
		if errors.As(err, &closed) {
			g.log.Warn().
				Param("code", closed.Code).
				Param("reason", closed.Text).
				Send("Connection closed")
			defer g.disconnect(false, true)
			return nil, fmt.Errorf("connection closed: %w", err)
		}
		return nil, fmt.Errorf("unexpected error: %w", err)
	}
	ev := eventPool.Get()
	if messageType == websocket.BinaryMessage {
		panic("not implemented")
	} else {
		_, err = g.buff.ReadFrom(reader)
		if err != nil {
			return nil, fmt.Errorf("error reading from connection reader: %w", err)
		}
		if err = json.Unmarshal(g.buff.Bytes(), ev); err != nil {
			return nil, fmt.Errorf("error unmarshalling event: %w", err)
		}
	}
	{
		msg := g.log.Trace().
			Param("op", ev.OpCode)
		if ev.Event != "" {
			msg.Param("event", ev.Event)
		}
		if len(ev.Data) > 0 {
			if len(ev.Data) > 1024*1024 {
				msg.Param("size", fmt.Sprintf("%dMiB", len(ev.Data)/1024/1024))
			} else if len(ev.Data) > 1024 {
				msg.Param("size", fmt.Sprintf("%dKiB", len(ev.Data)/1024))
			} else {
				msg.Param("size", fmt.Sprintf("%dB", len(ev.Data)))
			}
		}
		msg.Param("seq", ev.Seq).Send("Got message from Discord")
	}
	g.seq.Store(ev.Seq)
	return ev, nil
}

func (g *gatewayImpl) startHeartbeat(dur time.Duration) {
	log := g.log.Module("heartbeat")
	log.Trace().Send("Initializing")
	listener, cancel := g.Listen()
	defer cancel()
	timer := time.NewTimer(dur)
	bench := golog.CreateBenchmarkContext()
	for {
		select {
		case raw, ok := <-listener:
			if !ok {
				log.Trace().Send("Closing heartbeat, listener closed")
				return
			}
			switch data := raw.(type) {
			case InternalConnectionClosed:
				log.Trace().Send("Closing heartbeat, connection closed")
				return
			case InternalDispatchEvent:
				if data.OpCode == 11 {
					log.Scope("ACK").Debug().Duration(bench.Elapsed()).Send("Received heartbeat response")
				}
			}
		case <-timer.C:
			timer.Reset(dur)
			if err := g.conn.WriteJSON(sendEvent[uint64]{
				OpCode: 1,
				Data:   g.seq.Load(),
			}); err != nil {
				g.log.Error().Throw(fmt.Errorf("failed to send message to Discord: %w", err))
				break
			}
			bench.Update()
			log.Trace().Send("Sent, ACK should be in milliseconds...")
		}
	}
}

func (g *gatewayImpl) handshake(resume bool) error {
	var (
		hello HelloOp
		ready ReadyEvent
	)
	{
		ev, err := g.read()
		if err != nil {
			return fmt.Errorf("could not read hello op: %w", err)
		}
		if ev.OpCode != 10 {
			return fmt.Errorf("unexpected op code %d (wanted 10)", ev.OpCode)
		}
		if err = json.Unmarshal(ev.Data, &hello); err != nil {
			return fmt.Errorf("could not unmarshal hello op: %w", err)
		}
		g.log.Trace().Param("heartbeat_interval", hello.HeartbeatInterval).Send("First message received, Hello OP (10)")
	}
	if !resume {
		g.log.Trace().Send("That's new connection, sending Identify OP (2)")
		_ = g.conn.WriteJSON(sendEvent[Identify]{
			OpCode: 2,
			Data: Identify{
				Token: g.cfg.Token,
				Properties: IdentifyProperties{
					OS:      runtime.GOOS,
					Browser: "bfcord v0.0.1",
					Device:  "bfcord v0.0.1",
				},
				Compress: false,
				Shard:    []uint16{g.cfg.ID, g.cfg.ShardCount},
				Intents:  g.cfg.Intents,
			},
		})
		{
			ev, err := g.read()
			if err != nil {
				return fmt.Errorf("could not read ready event: %w", err)
			}
			if ev.OpCode != 0 {
				return fmt.Errorf("unexpected op code %d (wanted 0)", ev.OpCode)
			}
			if ev.Event != "READY" {
				return fmt.Errorf("unexpected event %s (wanted READY)", ev.Event)
			}
			if err = json.Unmarshal(ev.Data, &ready); err != nil {
				return fmt.Errorf("could not unmarshal ready event: %w", err)
			}
			g.log.Trace().
				Param("user", ready.User.Username).
				Param("user_id", ready.User.ID).
				Param("guilds", len(ready.Guilds)).
				Send("The last handshake message received - READY event - ready to read events")
			g.mut.Lock()
			g.resumeURL = ready.ResumeGatewayURL
			g.session = ready.SessionID
			g.mut.Unlock()
			g.sendEvent((InternalDispatchEvent)(ev))
		}
	} else {
		g.log.Trace().Send("Trying to resume the session")
		_ = g.conn.WriteJSON(sendEvent[resumeEvent]{
			OpCode: 6,
			Data: resumeEvent{
				Token:     g.cfg.Token,
				Seq:       g.seq.Load(),
				SessionID: g.session,
			},
		})
		ev, err := g.read()
		if err != nil {
			return fmt.Errorf("could not read resume event: %w", err)
		}
		if ev.OpCode != 7 {
			return fmt.Errorf("unexpected op code %d (wanted 7)", ev.OpCode)
		}
		g.log.Trace().Send("Received Resume OP (7), ready to read events")
	}
	g.log.Trace().Send("Starting event loop and heartbeat goroutines")
	go g.readEventsForever()
	go g.startHeartbeat(time.Duration(hello.HeartbeatInterval) * time.Millisecond)
	return nil
}

func (g *gatewayImpl) readEventsForever() {
	for {
		ev, err := g.read()
		if err != nil {
			if g.Status() == StatusDisconnected {
				return
			}
			g.log.Error().Throw(fmt.Errorf("could not read event: %w", err))
			g.disconnect(false, true)
			return
		}
		g.sendEvent((InternalDispatchEvent)(ev))
	}
}

func (g *gatewayImpl) Connect(ctx context.Context) error {
	g.log.Trace().Send("Preparing to connect to the '%s'", g.Config().URL)
	if g.conn != nil {
		g.log.Trace().Send("Connection was not closed, closing before connecting")
		g.disconnect(true, false)
		_ = g.conn.Close()
	}
	if ctx == nil {
		ctx = context.Background()
	}
	url := inlineif.IfElse(g.resumeURL == "", g.Config().URL, g.resumeURL)
	url += "?v=" + strings.TrimLeft(bfcord.APIVersion, "v")
	if g.cfg.Compression {
		url += "&compress=zlib-stream"
	}
	if g.cfg.GlobalLimiter != nil && g.resumeURL == "" {
		g.log.Trace().Send("Using identify global limiter")
		if err := g.cfg.GlobalLimiter.Wait(ctx); err != nil {
			return fmt.Errorf("failed to wait for identify rate limiter: %w", err)
		}
	}
	g.log.Trace().Send("Connecting to the '%s'", url)
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, url, nil)
	if err != nil {
		return fmt.Errorf("could not connect to the %s: %w", g.Config().URL, err)
	}
	g.conn = conn
	g.log.Trace().Send("Connection successfully created, handshaking with Discord Gateway")
	g.changeStatus(StatusConnecting)
	if err = g.handshake(g.resumeURL != ""); err != nil {
		g.disconnect(true, true)
		return fmt.Errorf("error while handshaking: %w", err)
	}
	g.changeStatus(StatusConnected)
	g.reconnections.Store(0)
	return nil
}

func (g *gatewayImpl) Config() Config {
	return g.cfg
}

func NewGateway(cfg Config) Gateway {
	if cfg.Logger == nil {
		cfg.Logger = golog.New("gateway")
	}
	if cfg.Intents == 0 {
		cfg.Intents = GatewayIntentDefault
	}
	gtw := new(gatewayImpl)
	gtw.log = cfg.Logger.Scope(fmt.Sprint(cfg.ID))
	gtw.cfg = cfg
	gtw.buff = bytes.NewBuffer(make([]byte, 0, 1024*1024))
	gtw.reset()
	gtw.status = StatusDisconnected
	gtw.reconnections = &atomic.Uint64{}
	gtw.log.Trace().Send("Gateway instance created")
	return gtw
}
