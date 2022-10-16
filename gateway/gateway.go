// This package allows communication with Discord through its gateway.
//
// We use pre-allocated buffers to prevent future allocations. You can adjust them using options.
//
// Each gateway is limited to ONE data handler
//
// IMPORTANT: Payload pointer that is passed as first data handler argument is SHARED.
// This means that the gateway uses it for every message that comes through the websocket, and DO NOT SHARE this pointer to another thread (goroutine).
// Payload processing from A to Z must be synchronous for gateway to wait when payload is free to be reused.

package gateway

import (
	"bytes"
	"context"
	"fmt"
	"github.com/BOOMfinity/go-utils/broadcaster"
	"io"
	"sync"
	"time"

	"github.com/BOOMfinity/bfcord/api"
	"github.com/BOOMfinity/bfcord/gateway/events"
	"github.com/BOOMfinity/bfcord/gateway/intents"
	"github.com/BOOMfinity/golog"
	"github.com/BOOMfinity/wshelper"
	"github.com/klauspost/compress/zlib"
	"github.com/segmentio/encoding/json"
	"github.com/unxcepted/websocket"
	"go.uber.org/atomic"
)

const (
	preallocateBufferSize = 1 << 20 // 1MiB
)

var panicStatusCodes = map[websocket.StatusCode]string{
	4004: "Authentication failed",
	4010: "Invalid shard",
	4011: "Sharding required",
	4012: "Invalid API version",
	4013: "Invalid intents",
	4014: "Disallowed intents",
}

var codesToReconnect = map[websocket.StatusCode]bool{
	4000: true,
	4001: true,
	4002: true,
	4003: true,
	4005: true,
	4008: true,
	4009: true,
}

type Gateway struct {
	heartbeatTime    time.Time
	sharedZlibReader io.ReadCloser
	onData           func(data *Payload)
	ws               *wshelper.Connection
	sharedBuffer     *bytes.Buffer
	sharedDecoder    *json.Decoder
	sharedPayload    *Payload
	seq              *atomic.Uint64
	channel          *broadcaster.Group[OpCode]
	eventChannel     *broadcaster.Group[events.Event]
	api              *api.Client
	membersChannel   map[string]chan<- RequestMembersResponse
	Logger           golog.Logger
	sessionID        string
	resumeGatewayURL string
	options          Options
	mut              sync.RWMutex
	resume           bool
}

func (g *Gateway) Presence() PresenceSet {
	return &presenceSet{g: g}
}

func (g *Gateway) EventChannel() *broadcaster.Group[events.Event] {
	return g.eventChannel
}

func (g *Gateway) OnData(fn func(data *Payload)) {
	g.onData = fn
}

func (g *Gateway) releaseBuffer() {
	if g.sharedBuffer.Cap() == g.options.bufferSize {
		g.sharedBuffer.Reset()
	} else {
		g.sharedBuffer = bytes.NewBuffer(make([]byte, 0, g.options.bufferSize))
		g.sharedDecoder = json.NewDecoder(g.sharedBuffer)
	}
}

func (g *Gateway) _onClose(_ *wshelper.Connection, code websocket.StatusCode, reason string) {
	if msg, ok := panicStatusCodes[code]; ok {
		panic(msg)
	}
	g.eventChannel.SendAsync(events.ShardDisconnected)
	if g.onData != nil {
		g.onData(staticDisconnectedPayload)
	}
	g.Logger.Warn().Send("Connection has been closed with code %v and reason %v", code, reason)
	if g.sessionID != "" && codesToReconnect[code] {
		_ = g.Reconnect(context.Background(), true)
	} else {
		_ = g.Reconnect(context.Background(), false)
	}
}

var (
	staticDisconnectedPayload      = &Payload{OP: DispatchOp, Event: events.ShardDisconnected, Data: nil}
	staticPrefetchCompletedPayload = &Payload{OP: DispatchOp, Event: events.ShardPrefetchCompleted, Data: nil}
)

func (g *Gateway) Disconnect(code websocket.StatusCode) {
	g.channel.SendAsync(111)
	g.ws.OnClose(nil)
	g.ws.OnError(nil)
	g.ws.OnMessage(nil)
	err := g.ws.Close(code, "-")
	if err == nil {
		if g.onData != nil {
			g.onData(staticDisconnectedPayload)
		}
		g.Logger.Warn().Send("Connection with Discord gateway has been closed")
	}
}

func (g *Gateway) Reconnect(ctx context.Context, resume bool) (err error) {
	g.resume = resume
	if resume {
		g.Disconnect(websocket.StatusServiceRestart)
	} else {
		g.Disconnect(websocket.StatusNormalClosure)
		g.sessionID = ""
	}
	g.Logger.Debug().Send("Waiting 5 seconds before reconnecting...")
	time.Sleep(5 * time.Second)
	g.Logger.Info().Send("Reconnecting...")
	return g.dial(ctx)
}

func (g *Gateway) dial(ctx context.Context) (err error) {
	var url string
	if g.resumeGatewayURL != "" {
		url = g.resumeGatewayURL
	} else {
		url, err = g.api.GatewayURL()
	}
	if err != nil {
		return
	}
	g.ws, err = wshelper.Dial(ctx, url+"?v=10&encoding=json", nil)
	if err != nil {
		return
	}
	g.ws.OnClose(g._onClose)
	g.ws.WS().SetReadLimit(819200000000)
	g.ws.OnMessageReader(g._onMessage)
	g.ws.OnError(func(_ *wshelper.Connection, err error) {
		panic(err)
	})
	return nil
}

func (g *Gateway) Connect(ctx context.Context) (err error) {
	g.seq = atomic.NewUint64(0)
	g.Logger.Info().Send("Connecting shard #%v", g.options.identify.Shard[0])
	return g.dial(ctx)
}

func (g *Gateway) _onMessage(_ *wshelper.Connection, t websocket.MessageType, data io.Reader) {
	var err error
	defer g.releaseBuffer()
	if t == websocket.MessageBinary {
		if g.sharedZlibReader == nil {
			g.sharedZlibReader, err = zlib.NewReader(data)
		} else {
			err = g.sharedZlibReader.(zlib.Resetter).Reset(data, nil)
		}
		if err != nil {
			g.Logger.Error().Send("failed create new zlib reader: %v", err.Error())
			return
		}
		_, err = g.sharedBuffer.ReadFrom(g.sharedZlibReader)
		if err != nil {
			g.Logger.Error().Send("failed read from zlib reader: %v", err.Error())
			return
		}
		_ = g.sharedZlibReader.Close()
		//g.Logger.Debug().Send("Received %v bytes of zlib-compressed data (%v bytes after decompression)", len(data), g.sharedBuffer.Len())
		err = g.sharedDecoder.Decode(g.sharedPayload)
		if err != nil {
			g.Logger.Error().Send("failed decoding json: %v", err.Error())
			return
		}
	} else {
		_, err = g.sharedBuffer.ReadFrom(data)
		if err != nil {
			g.Logger.Error().Send("failed reading from ws reader: %v", err.Error())
			_ = g.ws.Close(websocket.StatusAbnormalClosure, "-")
			return
		}
		err = g.sharedDecoder.Decode(g.sharedPayload)
		if err != nil {
			g.Logger.Error().Send("failed unmarshalling non-binary payload: %v", err.Error())
			_ = g.ws.Close(websocket.StatusAbnormalClosure, "-")
			return
		}
		//g.Logger.Debug().Send("Received %v bytes of uncompressed data", len(data))
	}
	g.channel.SendAsync(g.sharedPayload.OP)
	err = g.handleOp()
	if err != nil {
		g.Logger.Error().Add("OP %v", g.sharedPayload.OP).Send("Failed handling OP: %v", err.Error())
	}
}

func New(token string, shard uint16, totalShards uint16, options ...Option) *Gateway {
	def := Options{
		bufferSize: preallocateBufferSize,
		identify: Identify{
			Compress: true,
			Token:    token,
			Shard:    []uint16{shard, totalShards},
			Properties: IdentifyProperties{
				OS:      "linux",
				Browser: "bfcord",
				Device:  "v0.0.1",
			},
			Intents: intents.Default,
		},
	}
	for i := range options {
		options[i](&def)
	}
	if def.logger == nil {
		def.logger = golog.New("gateway")
	}
	g := new(Gateway)
	g.options = def
	g.channel = broadcaster.NewGroup[OpCode]()
	g.eventChannel = broadcaster.NewGroup[events.Event]()
	g.seq = atomic.NewUint64(0)
	g.Logger = def.logger.Module(fmt.Sprintf("#%v", shard))
	g.sharedPayload = &Payload{Shard: shard}
	g.sharedBuffer = bytes.NewBuffer(make([]byte, 0, def.bufferSize))
	g.sharedDecoder = json.NewDecoder(g.sharedBuffer)
	g.membersChannel = map[string]chan<- RequestMembersResponse{}
	if def.apiClient == nil {
		def.apiClient = api.NewClient(token, api.WithLogger(g.Logger.Module("api")))
	}
	g.api = def.apiClient
	return g
}
