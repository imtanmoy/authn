package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/imtanmoy/authn/auth"
	"github.com/imtanmoy/authn/internal/authx"
	"github.com/imtanmoy/authn/internal/errorx"
	"github.com/imtanmoy/authn/invite"
	"github.com/imtanmoy/authn/models"
	"github.com/imtanmoy/authn/user"
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

type signupPayload struct {
	Name            string `json:"name"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Designation     string `json:"designation"`
	Token           string `json:"token"`
}

func (sp *signupPayload) validate() url.Values {
	rules := govalidator.MapData{
		"name":             []string{"required", "min:4", "max:100"},
		"password":         []string{"required", "min:8", "max:20"},
		"confirm_password": []string{"required", "min:8", "max:20"},
		"designation":      []string{"min:2", "max:100"},
		"token":            []string{"required", "min:32", "max:32"},
	}
	opts := govalidator.Options{
		Data:  sp,
		Rules: rules,
	}

	v := govalidator.New(opts)
	e := v.ValidateStruct()
	if sp.Password != "" && sp.ConfirmPassword != "" {
		if sp.Password != sp.ConfirmPassword {
			e.Add("password", "password and confirmation password do not match")
			e.Add("confirm_password", "password and confirmation password do not match")
		}
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

// AuthHandler  represent the http handler for auth
type AuthHandler struct {
	useCase       auth.UseCase
	userUseCase   user.UseCase
	inviteUseCase invite.UseCase
	*authx.Authx
}

func (handler *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
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

	u, err := handler.useCase.FindByEmail(ctx, data.Email)
	if err != nil {
		if errors.Is(err, errorx.ErrorNotFound) {
			httpx.ResponseJSONError(w, r, http.StatusBadRequest, "invalid credentials", err)
		} else {
			httpx.ResponseJSONError(w, r, http.StatusInternalServerError, err)
		}
		return
	}
	if !handler.VerifyPassword(u, data.Password) {
		httpx.ResponseJSONError(w, r, http.StatusBadRequest, "invalid credentials", err)
		return
	}

	token, err := handler.GenerateToken(u.Email)
	if err != nil {
		httpx.ResponseJSONError(w, r, http.StatusInternalServerError, err)
		return
	}
	httpx.ResponseJSON(w, http.StatusCreated, struct {
		Token string `json:"token"`
	}{Token: token})
	return
}

func (handler *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	httpx.NoContent(w)
}

func (handler *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	u, err := handler.GetCurrentUser(r)
	us, ok := u.(*models.User)
	if err != nil || !ok {
		panic(fmt.Sprintf("could not upgrade user to an authable user, type: %T", u))
	}
	httpx.ResponseJSON(w, http.StatusOK, NewUserResponse(us))
	return
}

func (handler *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	data := &signupPayload{}
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
	tokenUser, err := handler.inviteUseCase.FindByToken(ctx, data.Token)
	if err != nil {
		if errors.Is(err, errorx.ErrorNotFound) {
			validationErrors.Add("token", "token not found")
			httpx.ResponseJSONError(w, r, 404, "invalid request", validationErrors)
		} else {
			httpx.ResponseJSONError(w, r, http.StatusInternalServerError, err)
		}
		return
	}
	if tokenUser.Status != "pending" {
		validationErrors.Add("token", "user already signed up with this token")
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
	u.Email = tokenUser.Email
	u.Password = hashedPassword

	err = handler.userUseCase.Store(ctx, &u)
	if err != nil {
		httpx.ResponseJSONError(w, r, http.StatusInternalServerError, err)
		return
	}

	tokenUser.Status = "successful"
	tokenUser.UserId = u.ID

	err = handler.inviteUseCase.Update(ctx, tokenUser)
	// TODO handle error
	httpx.ResponseJSON(w, http.StatusCreated, NewUserResponse(&u))
	return
}

// NewHandler will initialize the user's resources endpoint
func NewHandler(r *chi.Mux, useCase auth.UseCase, userUseCase user.UseCase, inviteUseCase invite.UseCase, aux *authx.Authx) {
	handler := &AuthHandler{
		useCase:       useCase,
		userUseCase:   userUseCase,
		inviteUseCase: inviteUseCase,
		Authx:         aux,
	}
	r.Route("/", func(r chi.Router) {
		r.Post("/login", handler.Login)
		r.Post("/signup", handler.SignUp)
		r.Group(func(r chi.Router) {
			r.Use(handler.AuthMiddleware)
			r.Post("/logout", handler.Logout)
			r.Get("/me", handler.GetMe)
		})
	})
}
