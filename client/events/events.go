package events

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/BOOMfinity/golog/v2"
)

type DispatcherError struct {
	ListenerID        uint64
	Event             string
	ListenerCreatedAt time.Time
	DeclaredAt        string
}

func (e DispatcherError) Error() string {
	return fmt.Sprintf("%s listener with ID %d declared at %s", e.Event, e.ListenerID, e.DeclaredAt)
}

type Listener[T SessionEvents] struct {
	ID         uint64
	CreatedAt  time.Time
	DeclaredAt string
	Nonce      bool
	handler    T
	cancel     ListenerCancelFn
}

type ListenerCancelFn func()
type DispatcherSendFn[T SessionEvents] func(handler T) error

type SessionEvents interface {
	ReadyEvent | GuildCreateEvent | GuildDeleteEvent | ChannelCreateEvent | ChannelUpdateEvent | ChannelDeleteEvent | MessageCreateEvent | MessageUpdateEvent | MessageDeleteEvent | ChannelPinsUpdateEvent | GuildUpdateEvent | ThreadCreateEvent | ThreadUpdateEvent | ThreadDeleteEvent | ThreadListSyncEvent | ThreadMembersUpdateEvent | GuildRoleAddEvent | GuildRoleUpdateEvent | GuildRoleDeleteEvent | GuildScheduledCreateEvent | GuildScheduledUpdateEvent | GuildScheduledDeleteEvent | GuildScheduledUserAddEvent | GuildScheduledUserRemoveEvent | GuildMemberAddEvent | GuildMemberUpdateEvent | GuildMemberRemoveEvent | InviteCreateEvent | InviteDeleteEvent | GuildBanAddEvent | GuildBanRemoveEvent | InteractionCreateEvent | VoiceServerUpdateEvent | VoiceStateUpdateEvent
}

type SessionDispatcher interface {
	Ready() Dispatcher[ReadyEvent]
	GuildCreate() Dispatcher[GuildCreateEvent]
	GuildDelete() Dispatcher[GuildDeleteEvent]
	ChannelCreate() Dispatcher[ChannelCreateEvent]
	ChannelUpdate() Dispatcher[ChannelUpdateEvent]
	ChannelPinsUpdate() Dispatcher[ChannelPinsUpdateEvent]
	ChannelDelete() Dispatcher[ChannelDeleteEvent]
	MessageCreate() Dispatcher[MessageCreateEvent]
	MessageUpdate() Dispatcher[MessageUpdateEvent]
	MessageDelete() Dispatcher[MessageDeleteEvent]
	GuildUpdate() Dispatcher[GuildUpdateEvent]
	ThreadCreate() Dispatcher[ThreadCreateEvent]
	ThreadUpdate() Dispatcher[ThreadUpdateEvent]
	ThreadDelete() Dispatcher[ThreadDeleteEvent]
	ThreadListSync() Dispatcher[ThreadListSyncEvent]
	ThreadMembersUpdate() Dispatcher[ThreadMembersUpdateEvent]
	GuildRoleAdd() Dispatcher[GuildRoleAddEvent]
	GuildRoleUpdate() Dispatcher[GuildRoleUpdateEvent]
	GuildRoleDelete() Dispatcher[GuildRoleDeleteEvent]
	GuildScheduledCreate() Dispatcher[GuildScheduledCreateEvent]
	GuildScheduledUpdate() Dispatcher[GuildScheduledUpdateEvent]
	GuildScheduledDelete() Dispatcher[GuildScheduledDeleteEvent]
	GuildScheduledUserAdd() Dispatcher[GuildScheduledUserAddEvent]
	GuildScheduledUserRemove() Dispatcher[GuildScheduledUserRemoveEvent]
	GuildMemberAdd() Dispatcher[GuildMemberAddEvent]
	GuildMemberRemove() Dispatcher[GuildMemberRemoveEvent]
	GuildMemberUpdate() Dispatcher[GuildMemberUpdateEvent]
	InviteCreate() Dispatcher[InviteCreateEvent]
	InviteDelete() Dispatcher[InviteDeleteEvent]
	GuildBanAdd() Dispatcher[GuildBanAddEvent]
	GuildBanRemove() Dispatcher[GuildBanRemoveEvent]
	InteractionCreate() Dispatcher[InteractionCreateEvent]
	VoiceStateUpdate() Dispatcher[VoiceStateUpdateEvent]
	VoiceServerUpdate() Dispatcher[VoiceServerUpdateEvent]
}

type Dispatcher[T SessionEvents] interface {
	Listen(fn T) ListenerCancelFn
	Nonce(fn T) ListenerCancelFn
	Sender(fn func(handler T))
}

type dispatcher[T SessionEvents] struct {
	log       golog.Logger
	listeners []*Listener[T]
	nils      []int
	mut       sync.RWMutex
	id        atomic.Uint64
}

func (d *dispatcher[T]) createListener() *Listener[T] {
	id := d.id.Add(1)
	listener := &Listener[T]{
		CreatedAt: time.Now(),
		ID:        id,
	}
	_, file, number, _ := runtime.Caller(2)
	listener.DeclaredAt = fmt.Sprintf("%s:%d", file, number)
	var index int
	if len(d.nils) > 0 {
		index, d.nils = d.nils[len(d.nils)-1], d.nils[:len(d.nils)-1]
		d.listeners[index] = listener
	} else {
		index = len(d.listeners)
		d.listeners = append(d.listeners, listener)
	}
	listener.cancel = func() {
		d.mut.Lock()
		d.listeners[index] = nil
		d.nils = append(d.nils, index)
		d.mut.Unlock()
	}
	return listener
}

func (d *dispatcher[T]) Listen(fn T) ListenerCancelFn {
	d.mut.Lock()
	defer d.mut.Unlock()

	listener := d.createListener()
	listener.handler = fn

	return listener.cancel
}

func (d *dispatcher[T]) Nonce(fn T) ListenerCancelFn {
	d.mut.Lock()
	defer d.mut.Unlock()

	listener := d.createListener()
	listener.handler = fn
	listener.Nonce = true

	return listener.cancel
}

func (d *dispatcher[T]) recoverPanic(l *Listener[T], fn func(handler T)) {
	defer func() {
		if err := recover(); err != nil {
			disErr := &DispatcherError{
				ListenerID:        l.ID,
				ListenerCreatedAt: l.CreatedAt,
				Event:             fmt.Sprintf("%T", l.handler),
				DeclaredAt:        l.DeclaredAt,
			}

			switch v := err.(type) {
			case error:
				d.log.Error().Throw(errors.Join(disErr, v))
			default:
				d.log.Error().Throw(errors.Join(disErr, fmt.Errorf("%v", v)))
			}
		}
	}()

	fn(l.handler)
}

func (d *dispatcher[T]) Sender(fn func(handler T)) {
	d.mut.RLock()
	listeners := d.listeners[:]
	d.mut.RUnlock()
	var wg sync.WaitGroup
	wg.Add(len(listeners))
	for _, listener := range listeners {
		go func() {
			defer wg.Done()
			d.recoverPanic(listener, fn)
		}()
		if listener.Nonce {
			listener.cancel()
		}
	}
	wg.Wait()
}

func NewDispatcher[T SessionEvents](log golog.Logger) Dispatcher[T] {
	return &dispatcher[T]{log: log}
}
