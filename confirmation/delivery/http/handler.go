package http

import (
	"github.com/go-chi/chi"
	"github.com/imtanmoy/authn/confirmation"
	"net/http"
)

// confirmationHandler  represent the http handler for Confirmation
type confirmationHandler struct {
	useCase confirmation.UseCase
}

func (handler *confirmationHandler) Confirm(w http.ResponseWriter, r *http.Request) {

}

func NewHandler(r *chi.Mux, useCase confirmation.UseCase) {
	handler := &confirmationHandler{
		useCase: useCase,
	}

	r.Route("/", func(r chi.Router) {
		r.Post("/confirm", handler.Confirm)
	})
}
