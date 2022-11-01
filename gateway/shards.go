package gateway

import (
	"github.com/BOOMfinity/go-utils/broadcaster"
	"sync"
	"time"
)

type ShardStatus uint8

const (
	ShardStatusDisconnected ShardStatus = iota + 1
	ShardStatusConnected
)

func NewShard(gateway *Gateway) *Shard {
	shard := new(Shard)
	shard.gtw = gateway
	shard.latency = make([]uint16, 0, 3)
	shard.status = ShardStatusDisconnected
	shard.m = new(sync.RWMutex)
	go shard.listenHeartbeat()
	return shard
}

type Shard struct {
	gtw     *Gateway
	m       *sync.RWMutex
	latency []uint16
	status  ShardStatus
}

func (v *Shard) SetStatus(status ShardStatus) {
	v.m.Lock()
	v.status = status
	v.m.Unlock()
	v.gtw.Logger.Debug().Send("Status changed to %T", status)
}

func (v *Shard) ID() uint16 {
	return v.gtw.options.identify.Shard[0]
}

func (v *Shard) listenHeartbeat() {
	member := v.gtw.channel.Join()
	defer member.Close()
	member.WithFilter(func(msg broadcaster.Message[OpCode]) bool {
		return msg.Data() == HeartbeatAckOp
	})
	for {
		_, more := member.Recv()
		if !more {
			return
		}
		t := time.Since(v.gtw.heartbeatTime).Milliseconds()
		v.m.Lock()
		if len(v.latency) == 0 {
			v.latency = append(v.latency, uint16(t))
		} else {
			v.latency = append([]uint16{uint16(t)}, v.latency[1:]...)
		}
		v.m.Unlock()
	}
}

func (v *Shard) Gateway() *Gateway {
	return v.gtw
}

func (v *Shard) Status() ShardStatus {
	return v.status
}

func (v *Shard) Latency() uint16 {
	v.m.RLock()
	defer v.m.RUnlock()
	if len(v.latency) == 0 {
		return 0
	}
	return v.latency[0]
}

func (v *Shard) LatencyHistory() []uint16 {
	v.m.RLock()
	defer v.m.RUnlock()
	return v.latency
}
