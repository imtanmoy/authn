package tests

import (
	"context"
	"github.com/imtanmoy/authn/events"
)

type event struct{}

func (e *event) Emit(ctx context.Context, eventName string, data interface{}) {
	//Do nothing
}

func (e *event) EmitWithDelay(ctx context.Context, eventName string, data interface{}) {
	//Do nothing
}

func NewMockEventEmitter() events.EventEmitter {
	return &event{}
}
