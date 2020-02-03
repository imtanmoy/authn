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
	"github.com/imtanmoy/authn/invite"
	"github.com/imtanmoy/authn/models"
	"github.com/imtanmoy/authn/organization"
	"github.com/imtanmoy/httpx"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/thedevsaddam/govalidator.v1"
	"log"
	"net/http"
	"net/url"
)

type invitePayload struct {
	Email          string `json:"email"`
	OrganizationId int    `json:"organization_id"`
}

func (ip *invitePayload) validate() url.Values {
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

// GenerateToken returns a unique token based on the provided email string
func GenerateToken(email string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(email), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	hasher := md5.New()
	hasher.Write(hash)
	return hex.EncodeToString(hasher.Sum(nil))
}

// InviteHandler  represent the http handler for invitation
type InviteHandler struct {
	useCase    invite.UseCase
	orgUseCase organization.UseCase
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

func (handler *InviteHandler) SendInvite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	u, err := handler.GetCurrentUser(r)
	us, ok := u.(*models.User)
	if err != nil || !ok {
		panic(fmt.Sprintf("could not upgrade user to an authable user, type: %T", u))
	}

	data := &invitePayload{}
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
			var iv models.Invite
			iv.Email = data.Email
			iv.Status = "pending"
			iv.InvitedBy = us.ID
			iv.OrganizationId = data.OrganizationId
			iv.Token = GenerateToken(data.Email)

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

// NewHandler will initialize the invite's resources endpoint
func NewHandler(r *chi.Mux, useCase invite.UseCase, orgUseCase organization.UseCase, au *authx.Authx) {
	handler := &InviteHandler{
		useCase:    useCase,
		orgUseCase: orgUseCase,
		Authx:      au,
	}
	r.Route("/invites", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(handler.AuthMiddleware)
			r.Post("/", handler.SendInvite)
			r.Get("/", handler.List)
			r.Get("/{id}", handler.Get)
			r.Delete("/{id}", handler.Delete)
		})
	})
}
