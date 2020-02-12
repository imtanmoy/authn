package events

import (
	"context"
	"fmt"
	_userEventHandler "github.com/imtanmoy/authn/events/handlers/user"
	"github.com/imtanmoy/authn/events/worker"
	"github.com/mustafaturan/bus"
	"github.com/mustafaturan/monoton"
	"github.com/mustafaturan/monoton/sequencer"
)

const (
	UserCreateEvent = "user:created"
	UserUpdateEvent = "user:updated"
)

type EventEmitter interface {
	Emit(ctx context.Context, eventName string, data interface{})
	EmitWithDelay(ctx context.Context, eventName string, data interface{})
}

type Event interface {
	Init()
	Close()
	EventEmitter
}

type EventData struct {
	Data    interface{}
	Delayed bool
}

type event struct {
	ctx           context.Context
	nonDelayedBus *bus.Bus
	delayedBus    *bus.Bus
	quit          chan bool
}

var _ Event = (*event)(nil)

func New() Event {
	nonDelayedBus := newBus()
	delayedBus := newBus()
	return &event{nonDelayedBus: nonDelayedBus, delayedBus: delayedBus, quit: make(chan bool)}
}

func (event *event) Init() {
	ctx, cancel := context.WithCancel(context.Background())
	event.ctx = ctx
	go func() {
		<-event.quit
		cancel()
	}()

	dispatcher := worker.NewDispatcher()
	dispatcher.Run(event.ctx)

	event.nonDelayedBus.RegisterTopics(UserCreateEvent, UserUpdateEvent)
	event.delayedBus.RegisterTopics(UserCreateEvent, UserUpdateEvent)
	event.nonDelayedBus.RegisterHandler("user_event_1", _userEventHandler.EventHandler(dispatcher.Send, false))
	event.delayedBus.RegisterHandler("user_event_2", _userEventHandler.EventHandler(dispatcher.Send, true))
}
func (event *event) Close() {
	event.quit <- true
}

func (event *event) Emit(ctx context.Context, eventName string, data interface{}) {
	_, err := event.nonDelayedBus.Emit(ctx, eventName, data)
	if err != nil {
		fmt.Println(err)
	}
}

func (event *event) EmitWithDelay(ctx context.Context, eventName string, data interface{}) {
	_, err := event.delayedBus.Emit(ctx, eventName, data)
	if err != nil {
		fmt.Println(err)
	}
}

func newBus() *bus.Bus {
	node := uint64(1)
	initialTime := uint64(1577865600000)
	m, err := monoton.New(sequencer.NewMillisecond(), node, initialTime)
	if err != nil {
		panic(err)
	}
	var idGenerator bus.Next = (*m).Next
	b, err := bus.NewBus(idGenerator)
	if err != nil {
		panic(err)
	}
	return b
}
