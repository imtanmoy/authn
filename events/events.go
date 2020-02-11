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
}

type event struct {
	ctx  context.Context
	b    *bus.Bus
	quit chan bool
}

var _ Event = (*event)(nil)

func New() Event {
	b := newBus()
	return &event{b: b, quit: make(chan bool)}
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

	event.b.RegisterTopics(UserCreateEvent, UserUpdateEvent)
	userCreateHandler := bus.Handler{Handle: func(e *bus.Event) {
		_userEventHandler.EventHandler(e, dispatcher.Send)
	}, Matcher: "^user.(created|updated)$"}
	event.b.RegisterHandler(UserCreateEvent, &userCreateHandler)
}
func (event *event) Close() {
	event.quit <- true
}

func (event *event) Emit(ctx context.Context, eventName string, data interface{}) {
	_, err := event.b.Emit(ctx, eventName, data)
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
