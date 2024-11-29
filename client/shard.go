package client

import (
	"sync"
	"time"

	"github.com/BOOMfinity/bfcord/utils"
	"github.com/BOOMfinity/bfcord/ws"
	"github.com/andersfylling/snowflake/v5"
)

type Shard interface {
	ws.Gateway

	Ping() int
	History() []int
	ID() uint16
	Unavailable() utils.SimpleMap[snowflake.ID, ws.UnavailableGuild]
}

type shardImpl struct {
	ws.Gateway

	ping        int
	history     []int
	unavailable utils.SimpleMap[snowflake.ID, ws.UnavailableGuild]

	mut sync.Mutex
}

func (s *shardImpl) Unavailable() utils.SimpleMap[snowflake.ID, ws.UnavailableGuild] {
	return s.unavailable
}

func (s *shardImpl) backgroundJob() {
	log := s.Log().Module("background")
	listener, cancel := s.Listen()
	defer cancel()
	for msg := range listener {
		switch data := msg.(type) {
		case ws.InternalHeartbeatEvent:
			s.mut.Lock()
			if len(s.history) >= 5 {
				s.history = append([]int{s.ping}, s.history[:4]...)
			}
			s.ping = int(time.Since(data.Start).Milliseconds())
			s.mut.Unlock()
			log.Debug().Param("ping", s.Ping()).Param("history", s.History()).Send("Received ACK event")
		}
	}
}

func (s *shardImpl) ID() uint16 {
	return s.Config().ID
}

func (s *shardImpl) Ping() int {
	s.mut.Lock()
	p := s.ping
	s.mut.Unlock()
	return p
}

func (s *shardImpl) History() []int {
	s.mut.Lock()
	h := s.history[:]
	s.mut.Unlock()
	return h
}
