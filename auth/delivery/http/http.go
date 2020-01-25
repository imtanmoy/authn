package http

import (
	"context"
	"errors"
	"github.com/go-chi/chi"
	"github.com/imtanmoy/authn/auth"
	"github.com/imtanmoy/authn/internal/authlib"
	"github.com/imtanmoy/authn/internal/errorx"
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
	if !authlib.ComparePasswords(data.Password, u.Password) {
		httpx.ResponseJSONError(w, r, http.StatusBadRequest, "invalid credentials", err)
		return
	}

	token, err := authlib.GenerateToken(u.Email)
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

// NewHandler will initialize the user's resources endpoint
func NewHandler(r *chi.Mux, useCase auth.UseCase) {
	handler := &AuthHandler{
		useCase: useCase,
	}
	r.Route("/", func(r chi.Router) {
		r.Post("/login", handler.Login)
		r.Group(func(r chi.Router) {
			r.Use(authlib.AuthMiddleware)
			r.Post("/logout", handler.Logout)
		})
	})
}
