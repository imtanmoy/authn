package http

import (
	"context"
	"errors"
	"github.com/go-chi/chi"
	"github.com/imtanmoy/authn/internal/errorx"
	"github.com/imtanmoy/authn/models"
	"github.com/imtanmoy/authn/user"
	"github.com/imtanmoy/httpx"
	param "github.com/oceanicdev/chi-param"
	"gopkg.in/thedevsaddam/govalidator.v1"
	"net/http"
	"net/url"
)

type contextKey string

const (
	userKey contextKey = "user"
)

type userPayload struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Designation     string `json:"designation"`
	OrganizationId  int    `json:"organization_id"`
}

func (o *userPayload) validate(ctx context.Context, useCase user.UseCase) url.Values {
	rules := govalidator.MapData{
		"name":             []string{"required", "min:4", "max:100"},
		"email":            []string{"required", "min:4", "max:100", "email"},
		"password":         []string{"required", "min:8", "max:20"},
		"confirm_password": []string{"required", "min:8", "max:20"},
		"designation":      []string{"min:2", "max:100"},
		"organization_id":  []string{"required", "numeric"},
	}
	opts := govalidator.Options{
		Data:  o,
		Rules: rules,
	}

	v := govalidator.New(opts)
	e := v.ValidateStruct()
	if o.Password != o.ConfirmPassword {
		e.Add("password", "password and confirmation password do not match")
		e.Add("confirm_password", "password and confirmation password do not match")
	}
	if useCase.ExistsByEmail(ctx, o.Email) {
		e.Add("email", "user with this email already exists")
	}
	return e
}

type userUpdatePayload struct {
	Name            string `json:"name"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Designation     string `json:"designation"`
}

func (u *userUpdatePayload) validate(ctx context.Context, useCase user.UseCase) url.Values {
	rules := govalidator.MapData{
		"name":             []string{"required", "min:4", "max:100"},
		"password":         []string{"min:8", "max:20"},
		"confirm_password": []string{"min:8", "max:20"},
		"designation":      []string{"min:2", "max:100"},
	}
	opts := govalidator.Options{
		Data:  u,
		Rules: rules,
	}

	v := govalidator.New(opts)
	e := v.ValidateStruct()
	if u.Password != u.ConfirmPassword {
		e.Add("password", "password and confirmation password do not match")
		e.Add("confirm_password", "password and confirmation password do not match")
	}
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
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if ctx == nil {
			httpx.ResponseJSONError(w, r, http.StatusInternalServerError, httpx.ErrInternalServerError)
			return
		}
		id, err := param.Int(r, "id")
		if err != nil {
			httpx.ResponseJSONError(w, r, http.StatusBadRequest, "invalid request parameter", err)
			return
		}
		org, err := uh.useCase.GetById(ctx, id)
		if err != nil {
			if errors.Is(err, errorx.ErrorNotFound) {
				httpx.ResponseJSONError(w, r, http.StatusNotFound, "user not found", err)
			} else {
				httpx.ResponseJSONError(w, r, http.StatusInternalServerError, err)
			}
			return
		}
		ctx = context.WithValue(r.Context(), userKey, org)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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

	validationErrors := data.validate(ctx, uh.useCase)

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
	ctx := r.Context()
	u, ok := ctx.Value(userKey).(*models.User)
	if !ok {
		httpx.ResponseJSONError(w, r, http.StatusInternalServerError, httpx.ErrInternalServerError)
		return
	}
	httpx.ResponseJSON(w, http.StatusOK, models.NewUserResponse(u))
	return
}

func (uh *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u, ok := ctx.Value(userKey).(*models.User)
	if !ok {
		httpx.ResponseJSONError(w, r, http.StatusInternalServerError, httpx.ErrInternalServerError)
		return
	}
	data := &userUpdatePayload{}
	if err := httpx.DecodeJSON(r, data); err != nil {
		var mr *httpx.MalformedRequest
		if errors.As(err, &mr) {
			httpx.ResponseJSONError(w, r, mr.Status, mr.Status, mr.Msg)
		} else {
			httpx.ResponseJSONError(w, r, http.StatusInternalServerError, err)
		}
		return
	}

	validationErrors := data.validate(ctx, uh.useCase)

	if len(validationErrors) > 0 {
		httpx.ResponseJSONError(w, r, 400, "invalid request", validationErrors)
		return
	}

	u.Name = data.Name
	u.Designation = data.Designation
	if data.Password != "" {
		u.Password = data.Password
	}

	err := uh.useCase.Update(ctx, u)
	if err != nil {
		httpx.ResponseJSONError(w, r, http.StatusInternalServerError, err)
		return
	}
	httpx.ResponseJSON(w, http.StatusCreated, models.NewUserResponse(u))
	return
}

func (uh *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u, ok := ctx.Value(userKey).(*models.User)
	if !ok {
		httpx.ResponseJSONError(w, r, http.StatusInternalServerError, httpx.ErrInternalServerError)
		return
	}
	err := uh.useCase.Delete(ctx, u)
	if err != nil {
		httpx.ResponseJSONError(w, r, http.StatusInternalServerError, err)
		return
	}
	httpx.NoContent(w)
}
