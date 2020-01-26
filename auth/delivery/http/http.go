package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/imtanmoy/authn/auth"
	"github.com/imtanmoy/authn/internal/authx"
	"github.com/imtanmoy/authn/internal/errorx"
	"github.com/imtanmoy/authn/models"
	"github.com/imtanmoy/httpx"
	"gopkg.in/thedevsaddam/govalidator.v1"
	"net/http"
	"net/url"
)

type loginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (o *loginPayload) validate() url.Values {
	rules := govalidator.MapData{
		"email":    []string{"required", "min:4", "max:100", "email"},
		"password": []string{"required", "min:8", "max:20"},
	}
	opts := govalidator.Options{
		Data:  o,
		Rules: rules,
	}

	v := govalidator.New(opts)
	e := v.ValidateStruct()
	return e
}

// AuthHandler  represent the http handler for auth
type AuthHandler struct {
	useCase auth.UseCase
	*authx.Authx
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	data := &loginPayload{}

	if err := httpx.DecodeJSON(r, data); err != nil {
		var mr *httpx.MalformedRequest
		if errors.As(err, &mr) {
			httpx.ResponseJSONError(w, r, mr.Status, mr.Status, mr.Msg)
		} else {
			httpx.ResponseJSONError(w, r, http.StatusInternalServerError, err)
		}
		return
	}

	validationErrors := data.validate()

	if len(validationErrors) > 0 {
		httpx.ResponseJSONError(w, r, 400, "invalid request", validationErrors)
		return
	}

	u, err := ah.useCase.FindByEmail(ctx, data.Email)
	if err != nil {
		if errors.Is(err, errorx.ErrorNotFound) {
			httpx.ResponseJSONError(w, r, http.StatusBadRequest, "invalid credentials", err)
		} else {
			httpx.ResponseJSONError(w, r, http.StatusInternalServerError, err)
		}
		return
	}
	if !ah.VerifyPassword(u, data.Password) {
		httpx.ResponseJSONError(w, r, http.StatusBadRequest, "invalid credentials", err)
		return
	}

	token, err := ah.GenerateToken(u.Email)
	if err != nil {
		httpx.ResponseJSONError(w, r, http.StatusInternalServerError, err)
		return
	}
	httpx.ResponseJSON(w, http.StatusCreated, struct {
		Token string `json:"token"`
	}{Token: token})
	return
}

func (ah *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	httpx.NoContent(w)
}

func (ah *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	u, err := ah.GetCurrentUser(r)
	us, ok := u.(*models.User)
	if err != nil || !ok {
		panic(fmt.Sprintf("could not upgrade user to an authable user, type: %T", u))
	}
	httpx.ResponseJSON(w, http.StatusOK, models.NewUserResponse(us))
	return
}

// NewHandler will initialize the user's resources endpoint
func NewHandler(r *chi.Mux, useCase auth.UseCase, aux *authx.Authx) {
	handler := &AuthHandler{
		useCase: useCase,
		Authx:   aux,
	}
	r.Route("/", func(r chi.Router) {
		r.Post("/login", handler.Login)
		r.Group(func(r chi.Router) {
			r.Use(handler.AuthMiddleware)
			r.Post("/logout", handler.Logout)
			r.Get("/me", handler.GetMe)
		})
	})
}
