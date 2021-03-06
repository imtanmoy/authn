package user

import (
	"github.com/imtanmoy/authn/confirmation"
	"github.com/imtanmoy/authn/models"
	"github.com/imtanmoy/logx"
	"github.com/mustafaturan/bus"
	"time"
)

func EventHandler(send func(fn func()), delayed bool) *bus.Handler {
	userHandler := bus.Handler{Handle: func(e *bus.Event) {
		var fn func()
		switch e.Topic {
		case "user:created":
			fn = createEventHandler(e.Data)
		default:
			logx.Errorf("whoops unexpected topic (%s)", e.Topic)
		}
		if fn != nil {
			if delayed {
				send(fn)
			} else {
				fn()
			}
		}
	}, Matcher: "^user.(created|updated)$"}
	return &userHandler
}

func createEventHandler(data interface{}) func() {
	u, ok := data.(models.User)
	if !ok {
		logx.Errorf("could not get email, type: %T", u)
	}

	fn := func() {
		time.Sleep(15 * time.Second)
		sendConfirmation(&u)
		logx.Infof("new user registered: %s", u.Email)

	}
	return fn
}

func sendConfirmation(u *models.User) {
	token := confirmation.GenerateConfirmationToken()
	logx.Infof("sending confirmation to %s with token %s", u.Email, token)
}
