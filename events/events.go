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

type Event interface {
	Init()
	Close()
	Emit(ctx context.Context, eventName string, data interface{})
	EmitWithDelay(ctx context.Context, eventName string, data interface{})
}

type EventData struct {
	Data    interface{}
	Delayed bool
}

type event struct {
	ctx  context.Context
	b1   *bus.Bus
	b2   *bus.Bus
	quit chan bool
}

var _ Event = (*event)(nil)

func New() Event {
	b1 := newBus()
	b2 := newBus()
	return &event{b1: b1, b2: b2, quit: make(chan bool)}
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

	event.b1.RegisterTopics(UserCreateEvent, UserUpdateEvent)
	event.b2.RegisterTopics(UserCreateEvent, UserUpdateEvent)
	event.b1.RegisterHandler("user_event_1", _userEventHandler.EventHandler(dispatcher.Send, false))
	event.b2.RegisterHandler("user_event_2", _userEventHandler.EventHandler(dispatcher.Send, true))
}
func (event *event) Close() {
	event.quit <- true
}

func (event *event) Emit(ctx context.Context, eventName string, data interface{}) {
	_, err := event.b1.Emit(ctx, eventName, data)
	if err != nil {
		fmt.Println(err)
	}
}

func (event *event) EmitWithDelay(ctx context.Context, eventName string, data interface{}) {
	_, err := event.b2.Emit(ctx, eventName, data)
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
