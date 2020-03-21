package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/imtanmoy/authn/internal/authx"
	"github.com/imtanmoy/authn/internal/errorx"
	"github.com/imtanmoy/authn/models"
	"github.com/imtanmoy/authn/organization"
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

func (u *userUpdatePayload) validate() url.Values {
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

type UserResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewUserResponse(u *models.User) *UserResponse {
	resp := &UserResponse{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
	return resp
}

func NewUserListResponse(users []*models.User) []*UserResponse {
	var list []*UserResponse
	if len(users) == 0 {
		list = make([]*UserResponse, 0)
	}
	for _, u := range users {
		list = append(list, NewUserResponse(u))
	}
	return list
}

// UserHandler  represent the http handler for user
type UserHandler struct {
	useCase    user.UseCase
	orgUseCase organization.UseCase
	*authx.Authx
}

// NewHandler will initialize the user's resources endpoint
func NewHandler(r *chi.Mux, useCase user.UseCase, orgUseCase organization.UseCase, au *authx.Authx) {
	handler := &UserHandler{
		useCase:    useCase,
		orgUseCase: orgUseCase,
		Authx:      au,
	}
	r.Route("/users", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(handler.AuthMiddleware)
			r.Get("/", handler.List)
			r.Post("/", handler.Create)
			r.Group(func(r chi.Router) {
				r.Use(handler.UserCtx)
				r.Get("/{id}", handler.Get)
				r.Put("/{id}", handler.Update)
				r.Delete("/{id}", handler.Delete)
			})
		})
	})
}

func (handler *UserHandler) UserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if ctx == nil {
			ctx = context.Background()
		}
		id, err := param.Int(r, "id")
		if err != nil {
			httpx.ResponseJSONError(w, r, http.StatusBadRequest, "invalid request parameter", err)
			return
		}
		u, err := handler.useCase.GetByID(ctx, id)
		if err != nil {
			if errors.Is(err, errorx.ErrorNotFound) {
				httpx.ResponseJSONError(w, r, http.StatusNotFound, "user not found", err)
			} else {
				httpx.ResponseJSONError(w, r, http.StatusInternalServerError, err)
			}
			return
		}
		ctx = context.WithValue(r.Context(), userKey, u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (handler *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	users, err := handler.useCase.FindAll(ctx)
	if err != nil {
		httpx.ResponseJSONError(w, r, http.StatusInternalServerError, "could not fetch user's list", err)
		return
	}

	httpx.ResponseJSON(w, http.StatusOK, NewUserListResponse(users))
	return
}

func (handler *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	cu, err := handler.GetCurrentUser(r)
	currentUser, ok := cu.(*models.User)
	if err != nil || !ok {
		panic(fmt.Sprintf("could not upgrade user to an authable user, type: %T", cu))
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

	validationErrors := data.validate(ctx, handler.useCase)

	org, err := handler.orgUseCase.GetById(ctx, data.OrganizationId)
	if err != nil {
		validationErrors.Add("organization_id", "organization not found")
	}

	found := false
	for _, b := range currentUser.Organizations {
		if b.ID == data.OrganizationId {
			found = true
			break
		}
	}
	if !found {
		validationErrors.Add("organization_id", "organization not found")
	}

	if len(validationErrors) > 0 {
		httpx.ResponseJSONError(w, r, 400, "invalid request", validationErrors)
		return
	}

	hashedPassword, err := handler.HashPassword(data.Password)
	if err != nil {
		httpx.ResponseJSONError(w, r, http.StatusInternalServerError, "could not create user, try again")
		return
	}

	var u models.User
	u.Name = data.Name
	u.Email = data.Email
	u.Password = hashedPassword

	err = handler.useCase.StoreWithOrg(ctx, &u, org)
	if err != nil {
		httpx.ResponseJSONError(w, r, http.StatusInternalServerError, err)
		return
	}
	httpx.ResponseJSON(w, http.StatusCreated, NewUserResponse(&u))
	return
}

func (handler *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u, ok := ctx.Value(userKey).(*models.User)
	if !ok {
		httpx.ResponseJSONError(w, r, http.StatusInternalServerError, httpx.ErrInternalServerError)
		return
	}
	httpx.ResponseJSON(w, http.StatusOK, NewUserResponse(u))
	return
}

func (handler *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	validationErrors := data.validate()

	if len(validationErrors) > 0 {
		httpx.ResponseJSONError(w, r, 400, "invalid request", validationErrors)
		return
	}

	u.Name = data.Name
	if data.Password != "" {
		hashedPassword, err := handler.HashPassword(data.Password)
		if err != nil {
			httpx.ResponseJSONError(w, r, http.StatusInternalServerError, "could not update user, try again")
			return
		}
		u.Password = hashedPassword
	}

	err := handler.useCase.Update(ctx, u)
	if err != nil {
		httpx.ResponseJSONError(w, r, http.StatusInternalServerError, "could not update user, try again", err)
		return
	}
	httpx.ResponseJSON(w, http.StatusCreated, NewUserResponse(u))
	return
}

func (handler *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u, ok := ctx.Value(userKey).(*models.User)
	if !ok {
		httpx.ResponseJSONError(w, r, http.StatusInternalServerError, httpx.ErrInternalServerError)
		return
	}
	err := handler.useCase.Delete(ctx, u)
	if err != nil {
		httpx.ResponseJSONError(w, r, http.StatusInternalServerError, "could not delete user, try again", err)
		return
	}
	httpx.NoContent(w)
}
