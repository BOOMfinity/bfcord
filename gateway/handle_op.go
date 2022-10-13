package gateway

import (
	"context"
	"github.com/BOOMfinity/go-utils/broadcaster"
	"time"

	"github.com/BOOMfinity/bfcord/gateway/events"
	"github.com/BOOMfinity/bfcord/gateway/intents"
	"github.com/BOOMfinity/bfcord/internal/bitfield"
	"github.com/segmentio/encoding/json"
)

func (g *Gateway) handleOp() error {
	switch g.sharedPayload.OP {
	case HelloOp:
		data, err := PayloadTo[helloData](g.sharedPayload)
		if err != nil {
			return err
		}
		if g.resume {
			g.resume = false
			_ = g.ws.WriteJSON(context.Background(), map[string]any{
				"op": ResumeOp,
				"d": map[string]any{
					"token":      g.options.identify.Token,
					"session_id": g.sessionID,
					"seq":        g.seq.Load(),
				},
			})
		} else {
			_ = g.ws.WriteJSON(context.Background(), map[string]any{
				"op": IdentifyOp,
				"d":  g.options.identify,
			})
		}
		go g.handleHeartbeat(time.Duration(data.Interval) * time.Millisecond)
	case InvalidSessionOp:
		resumable, err := PayloadTo[bool](g.sharedPayload)
		if err != nil {
			return err
		}
		g.Logger.Warn().Add("Resumable: %v", resumable).Send("Received invalid session")
		if resumable {
			_ = g.Reconnect(context.Background(), true)
		} else {
			_ = g.Reconnect(context.Background(), false)
		}
	case HeartbeatOp:
		g.Logger.Debug().Send("Discord requested the heartbeat")
		g.sendHeartbeat()
	case HeartbeatAckOp:
		g.Logger.Debug().Send("Heartbeat ACK (%v)", time.Since(g.heartbeatTime).String())
	case ReconnectOp:
		g.Logger.Debug().Send("Discord sent the reconnect OP")
		go g.Reconnect(context.Background(), true)
	case DispatchOp:
		g.seq.Store(g.sharedPayload.Seq)
		g.eventChannel.SendAsync(g.sharedPayload.Event)
		switch g.sharedPayload.Event {
		case events.ShardReady:
			data, err := PayloadTo[ReadyEvent](g.sharedPayload)
			if err != nil {
				return err
			}
			g.Logger.Debug().Any(data.Shard).Add("%v guilds", len(data.Guilds)).Send("Ready as %v", data.User.Tag())
			g.sessionID = data.SessionID
			g.resumeGatewayURL = data.ResumeGatewayURL
			if len(data.Guilds) > 0 && bitfield.Has(g.options.identify.Intents, intents.Guilds) {
				go g.handlePrefetching()
			} else {
				if g.onData != nil {
					g.onData(staticPrefetchCompletedPayload)
				}
				g.eventChannel.SendAsync(events.ShardPrefetchCompleted)
			}
		case events.ShardResumed:
			g.Logger.Warn().Send("Session resumed successfully")
		case events.GuildMembersChunk:
			data, err := PayloadTo[RequestMembersResponse](g.sharedPayload)
			if err != nil {
				return err
			}
			g.mut.RLock()
			ch, ok := g.membersChannel[data.Nonce]
			if ok {
				ch <- data
			}
			g.mut.RUnlock()
		}
		if g.onData != nil {
			g.onData(g.sharedPayload)
		}
	}
	return nil
}

func (g *Gateway) handlePrefetching() {
	member := g.eventChannel.Join()
	member.WithFilter(func(msg broadcaster.Message[events.Event]) bool {
		return msg.Data() == events.GuildCreate
	})
	opMember := g.channel.Join()
	defer member.Close()
	defer opMember.Close()
	timer := time.NewTimer(4 * time.Second)
	defer timer.Stop()
	for {
		select {
		case _, more := <-member.Out:
			if !more {
				return
			}
			timer.Reset(2 * time.Second)
		case msg, more := <-opMember.Out:
			if !more {
				return
			}
			if msg.Data() == 111 {
				return
			}
		case <-timer.C:
			if g.onData != nil {
				g.onData(staticPrefetchCompletedPayload)
			}
			g.eventChannel.SendAsync(events.ShardPrefetchCompleted)
			return
		}
	}
}

func (g *Gateway) sendHeartbeat() {
	g.heartbeatTime = time.Now()
	_ = g.ws.WriteJSON(context.Background(), map[string]any{
		"op": HeartbeatOp,
		"d":  g.seq.Load(),
	})
}

func (g *Gateway) handleHeartbeat(dur time.Duration) {
	g.Logger.Debug().Send("Starting heartbeat with interval: %v", dur.String())
	g.sendHeartbeat()
	member := g.channel.Join()
	timer := time.NewTicker(dur)
	for {
		select {
		case <-timer.C:
			g.sendHeartbeat()
		case msg := <-member.Out:
			if msg.Data() == 111 {
				g.Logger.Debug().Send("Exiting heartbeat loop")
				member.Close()
				return
			}
		}
	}
}

func PayloadTo[V any](p *Payload) (x V, err error) {
	err = json.Unmarshal(p.Data, &x)
	return
}
