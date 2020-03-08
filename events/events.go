package events

import (
	"context"
	"fmt"
	"github.com/gammazero/workerpool"
	_userEventHandler "github.com/imtanmoy/authn/events/handlers/user"
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

type EventBus interface {
	Init()
	Close()
	Run(ctx context.Context)
	EventEmitter
}

type EventData struct {
	Data    interface{}
	Delayed bool
}

type event struct {
	nonDelayedBus *bus.Bus
	delayedBus    *bus.Bus
	wp            *workerpool.WorkerPool
}

func (event *event) Run(ctx context.Context) {
	event.Init()
	<- ctx.Done()
	event.Close()
}

var _ EventBus = (*event)(nil)

func New() EventBus {
	nonDelayedBus := newBus()
	delayedBus := newBus()
	return &event{nonDelayedBus: nonDelayedBus, delayedBus: delayedBus}
}

func (event *event) Init() {
	event.wp = workerpool.New(2)
	event.nonDelayedBus.RegisterTopics(UserCreateEvent, UserUpdateEvent)
	event.delayedBus.RegisterTopics(UserCreateEvent, UserUpdateEvent)
	event.nonDelayedBus.RegisterHandler("user_event_non_delayed", _userEventHandler.EventHandler(event.wp.Submit, false))
	event.delayedBus.RegisterHandler("user_event_delayed", _userEventHandler.EventHandler(event.wp.Submit, true))
}
func (event *event) Close() {
	event.wp.StopWait()
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
