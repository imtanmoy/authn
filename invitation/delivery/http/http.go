package http

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/imtanmoy/authn/internal/authx"
	"github.com/imtanmoy/authn/internal/errorx"
	"github.com/imtanmoy/authn/invitation"
	"github.com/imtanmoy/authn/models"
	"github.com/imtanmoy/authn/organization"
	"github.com/imtanmoy/authn/user"
	"github.com/imtanmoy/httpx"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/thedevsaddam/govalidator.v1"
	"log"
	"net/http"
	"net/url"
)

type invitationPayload struct {
	Email          string `json:"email"`
	OrganizationId int    `json:"organization_id"`
}

func (ip *invitationPayload) validate() url.Values {
	rules := govalidator.MapData{
		"email":           []string{"required", "min:4", "max:100", "email"},
		"organization_id": []string{"required", "numeric"},
	}
	opts := govalidator.Options{
		Data:  ip,
		Rules: rules,
	}

	v := govalidator.New(opts)
	e := v.ValidateStruct()
	return e
}

type acceptInvitePayload struct {
	Name            string `json:"name"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Token           string `json:"token"`
}

func (sp *acceptInvitePayload) validate() url.Values {
	rules := govalidator.MapData{
		"name":             []string{"required", "min:4", "max:100"},
		"password":         []string{"required", "min:8", "max:20"},
		"confirm_password": []string{"required", "min:8", "max:20"},
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

// generateToken returns a unique token based on the provided email string
func generateToken(email string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(email), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	hasher := md5.New()
	hasher.Write(hash)
	return hex.EncodeToString(hasher.Sum(nil))
}

func (handler *InvitationHandler) List(w http.ResponseWriter, r *http.Request) {

}

func (handler *InvitationHandler) Get(w http.ResponseWriter, r *http.Request) {

}

func (handler *InvitationHandler) Accept(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	data := &acceptInvitePayload{}
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
	tokenUser, err := handler.useCase.FindByToken(ctx, data.Token)
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

	org, err := handler.orgUseCase.GetById(ctx, tokenUser.OrganizationId)

	if err != nil {
		httpx.ResponseJSONError(w, r, http.StatusBadRequest, "organization not found")
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

	var ou models.UserOrganization
	ou.UserId = u.ID
	ou.OrganizationId = org.ID
	ou.CreatedBy = tokenUser.InvitedBy

	err = handler.userUseCase.StoreWithOrg(ctx, &u, &ou)
	if err != nil {
		httpx.ResponseJSONError(w, r, http.StatusInternalServerError, err)
		return
	}

	tokenUser.Status = "successful"
	tokenUser.UserId = u.ID

	err = handler.useCase.Update(ctx, tokenUser)
	// TODO handle error
	httpx.ResponseJSON(w, http.StatusCreated, NewUserResponse(&u))
	return
}

func (handler *InvitationHandler) Delete(w http.ResponseWriter, r *http.Request) {

}

func (handler *InvitationHandler) SendInvite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	u, err := handler.GetCurrentUser(r)
	us, ok := u.(*models.User)
	if err != nil || !ok {
		panic(fmt.Sprintf("could not upgrade user to an authable user, type: %T", u))
	}

	data := &invitationPayload{}
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

	if !handler.orgUseCase.Exists(ctx, data.OrganizationId) {
		validationErrors.Add("organization_id", "organization does not exist")
		goto ReturnError
	}

	//if data.OrganizationId != 0 {
	//	if data.OrganizationId != us.Organizations[] {
	//		validationErrors.Add("organization_id", "invalid organization provided")
	//	}
	//}

ReturnError:
	if len(validationErrors) > 0 {
		httpx.ResponseJSONError(w, r, 400, "Invalid Request", validationErrors)
		return
	}

	cu, err := handler.useCase.FindByEmailAndOrganization(ctx, data.Email, data.OrganizationId)

	if cu != nil {
		if cu.Status == "pending" {
			httpx.ResponseJSONError(w, r, 400, "Invalid Request", "invitation already sent once")
			return
		}
		if cu.Status == "successful" {
			httpx.ResponseJSONError(w, r, 400, "Invalid Request", "user already joined")
			return
		}
	}

	if err != nil {
		if errors.Is(err, errorx.ErrorNotFound) {
			var iv models.Invitation
			iv.Email = data.Email
			iv.Status = "pending"
			iv.InvitedBy = us.ID
			iv.OrganizationId = data.OrganizationId
			iv.Token = generateToken(data.Email)

			err = handler.useCase.Store(ctx, &iv)
			if err != nil {
				httpx.ResponseJSONError(w, r, http.StatusInternalServerError, err)
				return
			}
			res := struct {
				Id             int    `json:"id"`
				Email          string `json:"email"`
				Status         string `json:"status"`
				InvitedBy      int    `json:"invited_by"`
				OrganizationId int    `json:"organization_id"`
			}{
				Id:             iv.ID,
				Email:          iv.Email,
				Status:         iv.Status,
				InvitedBy:      iv.InvitedBy,
				OrganizationId: iv.OrganizationId,
			}
			httpx.ResponseJSON(w, http.StatusCreated, &res)
			return
		}
	}

}

// InvitationHandler  represent the http handler for invitation
type InvitationHandler struct {
	useCase     invitation.UseCase
	userUseCase user.UseCase
	orgUseCase  organization.UseCase
	*authx.Authx
}

// NewHandler will initialize the Invitation's resources endpoint
func NewHandler(r *chi.Mux, useCase invitation.UseCase, userUseCase user.UseCase, orgUseCase organization.UseCase, au *authx.Authx) {
	handler := &InvitationHandler{
		useCase:     useCase,
		userUseCase: userUseCase,
		orgUseCase:  orgUseCase,
		Authx:       au,
	}
	r.Route("/invitations", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(handler.AuthMiddleware)
			r.Post("/", handler.SendInvite)
			r.Get("/", handler.List)
			r.Get("/{id}", handler.Get)
			r.Delete("/{id}", handler.Delete)
		})
		r.Group(func(r chi.Router) {
			r.Post("/accept", handler.Accept)
		})
	})
}
