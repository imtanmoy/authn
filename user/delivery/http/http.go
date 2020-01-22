package http

import (
	"context"
	"errors"
	"github.com/go-chi/chi"
	"github.com/imtanmoy/authn/models"
	"github.com/imtanmoy/authn/user"
	"github.com/imtanmoy/httpx"
	"gopkg.in/thedevsaddam/govalidator.v1"
	"net/http"
	"net/url"
)

type contextKey string

const (
	userKey contextKey = "user"
)

type userPayload struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	Designation    string `json:"designation"`
	OrganizationId int    `json:"organization_id"`
}

func (o *userPayload) validate() url.Values {
	rules := govalidator.MapData{
		"name":            []string{"required", "min:4", "max:100"},
		"email":           []string{"required", "min:4", "max:100", "email"},
		"password":        []string{"required", "min:8", "max:20"},
		"designation":     []string{"min:4", "max:100"},
		"organization_id": []string{"required", "numeric"},
	}
	opts := govalidator.Options{
		Data:  o,
		Rules: rules,
	}

	v := govalidator.New(opts)
	e := v.ValidateStruct()
	return e
}

// UserHandler  represent the http handler for user
type UserHandler struct {
	useCase user.UseCase
}

// NewHandler will initialize the user's resources endpoint
func NewHandler(r *chi.Mux, useCase user.UseCase) {
	handler := &UserHandler{
		useCase: useCase,
	}
	r.Route("/users", func(r chi.Router) {
		r.Get("/", handler.List)
		r.Post("/", handler.Create)
		r.Group(func(r chi.Router) {
			r.Use(handler.UserCtx)
			r.Get("/{id}", handler.Get)
			r.Put("/{id}", handler.Update)
			r.Delete("/{id}", handler.Delete)
		})
	})
}

func (uh *UserHandler) UserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}

func (uh *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	users, err := uh.useCase.FindAll(ctx)
	if err != nil {
		httpx.ResponseJSONError(w, r, http.StatusInternalServerError, err)
		return
	}

	httpx.ResponseJSON(w, http.StatusOK, models.NewUserListResponse(users))
	return
}

func (uh *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	data := &userPayload{}
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

	var u models.User
	u.Name = data.Name
	u.Email = data.Email
	u.Password = data.Password
	u.Designation = data.Designation
	u.OrganizationId = data.OrganizationId

	err := uh.useCase.Store(ctx, &u)
	if err != nil {
		httpx.ResponseJSONError(w, r, http.StatusInternalServerError, err)
		return
	}
	httpx.ResponseJSON(w, http.StatusCreated, models.NewUserResponse(&u))
	return
}

func (uh *UserHandler) Get(w http.ResponseWriter, r *http.Request) {

}

func (uh *UserHandler) Update(w http.ResponseWriter, r *http.Request) {

}

func (uh *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {

}
