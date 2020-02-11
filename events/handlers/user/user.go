package user

import (
	"github.com/imtanmoy/authn/models"
	"github.com/imtanmoy/logx"
	"github.com/mustafaturan/bus"
)

func EventHandler(e *bus.Event, send func(fn func())) {
	switch e.Topic {
	case "user:created":
		fn := createEventHandler(e.Data)
		send(fn)
	default:
		logx.Errorf("whoops unexpected topic (%s)", e.Topic)
	}

}

func createEventHandler(data interface{}) func() {
	u, ok := data.(models.User)
	if !ok {
		logx.Errorf("could not get email, type: %T", u)
	}

	fn := func() {
		logx.Infof("new user registered: %s", u.Email)
	}
	return fn
}
