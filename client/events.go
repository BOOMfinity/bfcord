package client

import (
	"fmt"
	"time"

	"github.com/BOOMfinity/go-utils/gpool"
	"github.com/BOOMfinity/golog/v2"
	"github.com/segmentio/encoding/json"

	"github.com/BOOMfinity/bfcord/ws"
)

func (s *sessionImpl) handleEvents(shard Shard) {
	pool := golog.NewPool(s.log.Module("event-handler"))
	listener, cancel := shard.Listen()
	defer cancel()
	for msg := range listener {
		switch data := msg.(type) {
		case ws.InternalDispatchEvent:
			handler, _ := s.handlers.Get(data.Event)
			if handler != nil {
				go func() {
					bench := golog.AcquireBenchmarkContext()
					defer golog.ReleaseBenchmarkContext(bench)
					log := pool.Get()
					defer pool.Put(log)
					log.Param("shard", shard.ID()).Param("event", data.Event)
					defer log.Trace().Duration(bench.Elapsed())
					bench.Update()
					if err := handler(log, s, shard, data); err != nil {
						log.Error().Throw(fmt.Errorf("failed to execute event handler: %w", err))
						return
					}
					log.Trace().Duration(bench.Elapsed()).Send("Event processed")
					s.metrics.events.Add(1)
					s.metrics.totalTime.Add(uint64(bench.Total().Nanoseconds()))
				}()
			}
		}
	}
}

func (s *sessionImpl) metricsService() {
	log := s.log.Module("metrics")

	for {
		time.Sleep(30 * time.Second)
		events := s.metrics.events.Swap(0)
		totalTime := s.metrics.totalTime.Swap(0)
		if totalTime == 0 || events == 0 {
			continue
		}
		log.Debug().Send("Processed %d events with average execution time of %s (%s total)", events, time.Duration(totalTime/events).String(), time.Duration(totalTime).String())
	}
}

type handleDispatchFn func(log golog.Logger, sess Session, shard Shard, data ws.InternalDispatchEvent) error
type handleEventFn[T any] func(log golog.Logger, sess Session, _ ws.InternalDispatchEvent, _ Shard, data *T)

func handle[T any](fn handleEventFn[T]) handleDispatchFn {
	pool := gpool.New[T]()
	return func(log golog.Logger, sess Session, shard Shard, v ws.InternalDispatchEvent) error {
		defer v.Dereference()
		obj := pool.Get()
		defer pool.Put(obj)
		if err := json.Unmarshal(v.Data, obj); err != nil {
			return fmt.Errorf("failed to unmarshal %T for %s event: %w", v, v.Event, err)
		}
		fn(log, sess, v, shard, obj)
		return nil
	}
}
