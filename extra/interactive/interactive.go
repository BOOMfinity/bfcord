// Package interactive contains basic API for creating menus with buttons and selects.
//
/*
	menu := NewInteractive(manager, interaction)
	yes := func(i *interactions.Interaction) bool {
		i.ComponentMessageUpdate().Content("Message successfully deleted!").Execute()
		return false
	}
	no := func(i *interactions.Interaction) bool {
		i.ComponentMessageUpdate().Content("Message deletion has been canceled!").Execute()
		return false
	}
	interaction.SendMessageReply().
		Content("Are you sure you want to delete this message?").
		ActionRow(
			components.NewButton(components.ButtonStyleSuccess).SetCustomID("yes").SetLabel("Yes!"),
			components.NewButton(components.ButtonStyleDanger).SetCustomID("nope").SetLabel("Nope!")).
		Execute()
	menu.SetFilter(func(i *interactions.Interaction) bool {
		return i.UserID() == interaction.UserID()
	})
	menu.Start(interactive.Actions{
		"yes": yes,
		"nope": no,
	})

The boolean output variable in action handler means if you want to listen for another interaction (if false, stops listening and exits a loop).
*/
package interactive

import (
	"reflect"
	"sync"
	"time"

	"github.com/BOOMfinity/bfcord/client"

	"github.com/BOOMfinity/bfcord/discord/interactions"
	"github.com/BOOMfinity/bfcord/gateway"
	"github.com/andersfylling/snowflake/v5"
)

type Manager struct {
	listeners *sync.Map
	discord   client.Client
}

func (m *Manager) init() *Manager {
	m.discord.Sub().Interaction(m.onInteractionCreate)
	return m
}

func (m *Manager) listen(msgID snowflake.ID, listener chan<- *interactions.Interaction) {
	m.listeners.Store(msgID, listener)
}

func (m *Manager) stopListen(msgID snowflake.ID) {
	m.listeners.Delete(msgID)
}

func (m *Manager) onInteractionCreate(_ client.Client, _ *gateway.Shard, ev *interactions.Interaction) {
	if !ev.InGuild() || !ev.IsMessageComponent() {
		return
	}
	val, ok := m.listeners.LoadAndDelete(ev.Message.ID)
	if ok {
		val.(chan<- *interactions.Interaction) <- ev
	}
}

// Return false if you no longer want to listen for interactions.
type Handler func(i *interactions.Interaction) bool

type Actions map[string]Handler

type Interactive struct {
	attachedTo    *interactions.Interaction
	timer         *time.Timer
	onTimeout     func(i *Interactive)
	filterHandler func(i *interactions.Interaction) bool
	manager       *Manager
	quit          chan bool
	timeout       time.Duration
}

func (i *Interactive) Close() {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	reflect.ValueOf(i.quit).TrySend(reflect.ValueOf(true))
}

func (i *Interactive) Start(actions Actions) error {
	c := make(chan *interactions.Interaction)
	if i.attachedTo.Message.ID.IsZero() {
		if err := i.attachedTo.FetchOriginalResponse(); err != nil {
			return err
		}
	}
	i.manager.listen(i.attachedTo.Message.ID, c)
	if i.timer == nil {
		i.timer = time.NewTimer(i.timeout)
	} else {
		i.timer.Reset(i.timeout)
	}
	for {
		select {
		case <-i.quit:
			return nil
		case trigger := <-c:
			action, ok := actions[trigger.Data.CustomID]
			if !ok {
				break
			}
			if i.filterHandler != nil {
				if !i.filterHandler(trigger) {
					break
				}
			}
			i.timer.Stop()
			if !action(trigger) {
				return nil
			}
			i.ResetTimeout()
			i.manager.listen(i.attachedTo.Message.ID, c)
		case <-i.timer.C:
			if i.onTimeout != nil {
				i.onTimeout(i)
			}
			return nil
		}
	}
}

func (i *Interactive) ResetTimeout() {
	i.timer.Reset(i.timeout)
}

func (i *Interactive) WithTimeout(time time.Duration) *Interactive {
	i.timeout = time
	return i
}

func (i *Interactive) OnTimeout(handler func(i *Interactive)) {
	i.onTimeout = handler
}

func (i *Interactive) SetFilter(handler func(i *interactions.Interaction) bool) {
	i.filterHandler = handler
}

func NewInteractive(manager *Manager, i *interactions.Interaction) *Interactive {
	return &Interactive{
		manager:    manager,
		attachedTo: i,
		quit:       make(chan bool),
	}
}

func NewManager(client client.Client) *Manager {
	x := &Manager{
		listeners: new(sync.Map),
		discord:   client,
	}
	return x.init()
}
