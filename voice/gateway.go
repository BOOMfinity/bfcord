package voice

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/BOOMfinity/wshelper"
	"github.com/unxcepted/websocket"
	"io"
	"strings"
	"time"
)

const (
	maxReconnectAttempts = 3
	gatewayVersion       = 4
)

// openGateway method opens connection to voice gateway and identifies itself. Further connection steps are handled automatically.
func (s *Session) openGateway(ctx context.Context) (err error) {
	if s.state.SessionID == "" || s.server.Token == "" || s.server.Endpoint == "" {
		return GatewayMissingValuesError
	}

	s.ws, err = wshelper.Dial(ctx, fmt.Sprintf("wss://%v?v=%v", strings.TrimSuffix(s.server.Endpoint, ":80"), gatewayVersion), nil)
	if err != nil {
		return err
	}

	s.ws.OnMessageReader(s.onMessage)
	s.ws.OnClose(s.onClose)
	s.ws.OnError(func(_ *wshelper.Connection, err error) {
		panic(err)
	})

	s.logger.Debug().Send("Gateway opened")
	return s.identify()
}

func (s *Session) onMessage(_ *wshelper.Connection, _ websocket.MessageType, data io.Reader) {
	var v Payload
	err := json.NewDecoder(data).Decode(&v)
	if err != nil {
		s.logger.Error().Send("Error while decoding ws payload: %v", err.Error())
	}

	s.logger.Debug().Send("Received OPCode %v", v.Op)

	switch v.Op {
	case readyOP:
		ready, err := payloadTo[ReadyPayload](&v)
		if err != nil {
			s.logger.Error().Send("Error while decoding ready payload: %v. Can not continue.", err.Error())
			_ = s.ws.Close(websocket.StatusProtocolError, "malformed ready payload")
			return
		}
		s.ready = &ready

		s.udp, err = newUDP(context.Background(), fmt.Sprintf("%v:%v", ready.IP, ready.Port), ready.SSRC, s.logger)
		if err != nil {
			s.logger.Error().Send("Error while creating UDP connection: %v.", err.Error())
			_ = s.ws.Close(BfcordVoiceReconnect, "internal library error")
		} else {
			err := s.selectProtocol()
			if err != nil {
				s.logger.Error().Send("Error while sending select protocol: %v.", err.Error())
				_ = s.ws.Close(BfcordVoiceReconnect, "internal library error")
			}
		}
	case helloOP:
		hello, err := payloadTo[helloPayload](&v)
		if err != nil {
			s.logger.Error().Send("Error while decoding hello payload: %v. Can not continue.", err.Error())
			_ = s.ws.Close(websocket.StatusProtocolError, "malformed hello payload")
			return
		}
		go s.heartbeatLoop(time.Duration(hello.HeartbeatInterval) * time.Millisecond)
	case sessionDescriptionOP:
		desc, err := payloadTo[sessionDescriptionPayload](&v)
		if err != nil {
			s.logger.Error().Send("Error while decoding SessionDescription payload: %v. Can not continue.", err.Error())
			_ = s.ws.Close(websocket.StatusProtocolError, "malformed sessiondescription payload")
			return
		}
		s.udp.PutSecretKey(desc.SecretKey)
		s.reconnectAttempt = 0
		s.eventChannel.SendAsync(101)
	case resumedOP:
		s.logger.Info().Send("Connection resumed")
		s.reconnectAttempt = 0
		s.eventChannel.SendAsync(101)
		_ = s.SendSpeaking(Microphone)
	}
}

func (s *Session) onClose(_ *wshelper.Connection, code websocket.StatusCode, reason string) {
	s.reconnectAttempt++
	s.eventChannel.SendAsync(111)
	if s.reconnectAttempt > maxReconnectAttempts {
		s.logger.Warn().Send("Gateway closed: max reconnect attempts. Destroying connection.")
		s.destroy(true)
		return
	}

	switch code {
	case 1001, 1006, 4000, ServerCrashed, 4969: // resume
		s.logger.Warn().Send("Gateway closed with code %v. Trying to resume. (#%v/%v)", code, s.reconnectAttempt, maxReconnectAttempts)
		_ = s.openGateway(context.Background())
	default: // probably intended disconnection or unrecoverable (abnormal) error - destroy connection
		s.logger.Warn().Send("Gateway closed with code %v. Destroying connection.", code)
		s.destroy(false)
	}
}

func (s *Session) heartbeatLoop(interval time.Duration) {
	s.logger.Debug().Send("Starting heartbeat with interval: %v", interval.String())
	member := s.eventChannel.Join()
	timer := time.NewTicker(interval)
	defer func() {
		timer.Stop()
		member.Close()
	}()
	for {
		select {
		case <-timer.C:
			s.sendHeartbeat()
		case msg := <-member.Out:
			if msg.Data() == 111 {
				return
			}
		}
	}
}

func (s *Session) sendHeartbeat() {
	err := s.ws.WriteJSON(context.Background(), map[string]any{
		"op": heartbeatOP,
		"d":  time.Now().Unix(),
	})

	if err != nil {
		s.logger.Error().Send("Error while sending heartbeat: %v", err.Error())
	}
}

func (s *Session) identify() error {
	if s.ready != nil { // has already identified once - resume. If session is invalid WS will probably close with code 4006 and fully reconnect
		s.logger.Debug().Send("Resuming connection")
		return s.ws.WriteJSON(context.Background(), map[string]any{
			"op": resumeOP,
			"d": resumePayload{
				ServerID:  s.state.GuildID,
				SessionID: s.state.SessionID,
				Token:     s.server.Token,
			},
		})
	}

	s.logger.Debug().Send("Identifying connection")
	return s.ws.WriteJSON(context.Background(), map[string]any{
		"op": identifyOP,
		"d": IdentifyPayload{
			SessionID: s.state.SessionID,
			Token:     s.server.Token,
			ServerID:  s.state.GuildID,
			UserID:    s.state.UserID,
		},
	})
}

func (s *Session) selectProtocol() error {
	return s.ws.WriteJSON(context.Background(), map[string]any{
		"op": selectProtocolOP,
		"d": selectProtocolPayload{
			Protocol: "udp",
			Data: selectProtocolData{
				Address:        s.udp.OwnIP,
				Port:           s.udp.OwnPort,
				EncryptionMode: "xsalsa20_poly1305",
			},
		},
	})
}

// SendSpeaking method sends speaking flags to discord.
//
// It's automatically invoked during (re)connection.
//
// Will throw an error if connection is not yet finished.
func (s *Session) SendSpeaking(flags SpeakingFlag) error {
	if s.ready == nil {
		return NotConnected
	}
	return s.ws.WriteJSON(context.Background(), map[string]any{
		"op": speakingOP,
		"d": speakingPayload{
			Speaking: flags,
			SSRC:     s.ready.SSRC,
		},
	})
}
