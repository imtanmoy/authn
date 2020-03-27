package http

import (
	"github.com/go-chi/chi"
	"github.com/imtanmoy/authn/events"
	"github.com/imtanmoy/authn/internal/authx"
	"github.com/imtanmoy/authn/organization"
	"net/http"
)

// orgHandler  represent the http handler for org
type orgHandler struct {
	useCase organization.UseCase
	*authx.Authx
	event events.EventEmitter
}

func (handler *orgHandler) Create(w http.ResponseWriter, r *http.Request) {

}

// NewHandler will initialize the org's resources endpoint
func NewHandler(
	r *chi.Mux,
	aux *authx.Authx,
	useCase organization.UseCase,
	event events.EventEmitter,
) {
	handler := &orgHandler{
		useCase: useCase,
		Authx:   aux,
		event:   event,
	}
	r.Route("/organizations", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(handler.AuthMiddleware)
			r.Post("/", handler.Create)
		})
	})
}
