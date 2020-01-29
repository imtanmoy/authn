package http

import (
	"github.com/go-chi/chi"
	"github.com/imtanmoy/authn/internal/authx"
	"github.com/imtanmoy/authn/invite"
	"net/http"
)

// InviteHandler  represent the http handler for invitation
type InviteHandler struct {
	useCase invite.UseCase
	*authx.Authx
}

func (handler *InviteHandler) List(w http.ResponseWriter, r *http.Request) {

}

func (handler *InviteHandler) Get(w http.ResponseWriter, r *http.Request) {

}

func (handler *InviteHandler) Accept(w http.ResponseWriter, r *http.Request) {

}

func (handler *InviteHandler) Delete(w http.ResponseWriter, r *http.Request) {

}

// NewHandler will initialize the invite's resources endpoint
func NewHandler(r *chi.Mux, useCase invite.UseCase, au *authx.Authx) {
	handler := &InviteHandler{
		useCase: useCase,
		Authx:   au,
	}
	r.Route("/invites", func(r chi.Router) {
		r.Get("/{token}", handler.Accept)
		r.Group(func(r chi.Router) {
			r.Use(handler.AuthMiddleware)
			r.Get("/", handler.List)
			r.Get("/{id}", handler.Get)
			r.Delete("/{id}", handler.Delete)
		})
	})
}
