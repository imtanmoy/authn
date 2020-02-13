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
	"github.com/imtanmoy/httpx"
	param "github.com/oceanicdev/chi-param"
	"gopkg.in/thedevsaddam/govalidator.v1"
	"net/http"
	"net/url"
	"time"
)

type contextKey string

const (
	orgKey contextKey = "organization"
)

type organizationPayload struct {
	Name string `json:"name"`
}

func (o *organizationPayload) validate() url.Values {
	rules := govalidator.MapData{
		"name": []string{"required", "min:4", "max:100"},
	}
	opts := govalidator.Options{
		Data:  o,
		Rules: rules,
	}

	v := govalidator.New(opts)
	e := v.ValidateStruct()
	return e
}

type OrganizationResponse struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	JoinedAt time.Time `json:"joined_at"`
}

func NewOrganizationResponse(organization *models.Organization) *OrganizationResponse {
	resp := &OrganizationResponse{
		ID:   organization.ID,
		Name: organization.Name,
	}
	return resp
}

// OrganizationHandler  represent the http handler for organization
type OrganizationHandler struct {
	useCase organization.UseCase
	*authx.Authx
}

// NewHandler will initialize the organization's resources endpoint
func NewHandler(r *chi.Mux, useCase organization.UseCase, aux *authx.Authx) {
	handler := &OrganizationHandler{
		useCase: useCase,
		Authx:   aux,
	}
	r.Route("/organizations", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(handler.AuthMiddleware)
			r.Post("/", handler.Create)
			r.Get("/", handler.List)
			r.Group(func(r chi.Router) {
				r.Route("/{id}", func(r chi.Router) {
					r.Use(handler.OrganizationCtx)
					r.Get("/", handler.Get)
					r.Put("/", handler.Update)
					r.Delete("/", handler.Delete)
				})
			})
		})
	})
}

func (handler *OrganizationHandler) OrganizationCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if ctx == nil {
			httpx.ResponseJSONError(w, r, http.StatusInternalServerError, httpx.ErrInternalServerError)
			return
		}
		id, err := param.Int(r, "id")
		if err != nil {
			httpx.ResponseJSONError(w, r, http.StatusBadRequest, "Invalid request parameter", err)
			return
		}
		org, err := handler.useCase.GetById(ctx, id)
		if err != nil {
			if errors.Is(err, errorx.ErrorNotFound) {
				httpx.ResponseJSONError(w, r, http.StatusNotFound, "organization not found", err)
			} else {
				httpx.ResponseJSONError(w, r, http.StatusInternalServerError, err)
			}
			return
		}
		u, err := handler.GetCurrentUser(r)
		currentUser, ok := u.(*models.User)
		if err != nil || !ok {
			panic(fmt.Sprintf("could not upgrade user to an authable user, type: %T", u))
		}
		found := false
		for _, b := range currentUser.Organizations {
			if b.ID == org.ID {
				found = true
				break
			}
		}
		if !found {
			httpx.ResponseJSONError(w, r, http.StatusBadRequest, "organization not found")
			return
		}
		ctx = context.WithValue(r.Context(), orgKey, org)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (handler *OrganizationHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	u, err := handler.GetCurrentUser(r)
	currentUser, ok := u.(*models.User)
	if err != nil || !ok {
		panic(fmt.Sprintf("could not upgrade user to an authable user, type: %T", u))
	}
	organizations, err := handler.useCase.FindAllByUserId(ctx, currentUser.ID)
	if err != nil {
		httpx.ResponseJSONError(w, r, http.StatusInternalServerError, err)
		return
	}

	httpx.ResponseJSON(w, http.StatusOK, organizations)
	return
}

func (handler *OrganizationHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	u, err := handler.GetCurrentUser(r)
	currentUser, ok := u.(*models.User)
	if err != nil || !ok {
		panic(fmt.Sprintf("could not upgrade user to an authable user, type: %T", u))
	}
	data := &organizationPayload{}
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
		httpx.ResponseJSONError(w, r, 400, "Invalid Request", validationErrors)
		return
	}

	var org models.Organization
	org.Name = data.Name
	org.OwnerId = currentUser.ID

	err = handler.useCase.Store(ctx, &org, currentUser)
	if err != nil {
		httpx.ResponseJSONError(w, r, http.StatusInternalServerError, err)
		return
	}
	httpx.ResponseJSON(w, http.StatusCreated, NewOrganizationResponse(&org))
	return
}

func (handler *OrganizationHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	org, ok := ctx.Value(orgKey).(*models.Organization)
	if !ok {
		panic(fmt.Sprintf("could not get organization, type: %T", org))
	}
	httpx.ResponseJSON(w, http.StatusOK, NewOrganizationResponse(org))
	return
}

func (handler *OrganizationHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	org, ok := ctx.Value(orgKey).(*models.Organization)
	if !ok {
		panic(fmt.Sprintf("could not get organization, type: %T", org))
	}
	data := &organizationPayload{}
	if err := httpx.DecodeJSON(r, data); err != nil {
		httpx.ResponseJSONError(w, r, http.StatusBadRequest, err)
		return
	}

	validationErrors := data.validate()

	if len(validationErrors) > 0 {
		httpx.ResponseJSONError(w, r, http.StatusBadRequest, "Invalid Request", validationErrors)
		return
	}
	// update organization's data
	org.Name = data.Name
	err := handler.useCase.Update(ctx, org)
	if err != nil {
		httpx.ResponseJSONError(w, r, http.StatusInternalServerError, err)
		return
	}
	httpx.ResponseJSON(w, http.StatusOK, NewOrganizationResponse(org))
	return
}

func (handler *OrganizationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	org, ok := ctx.Value(orgKey).(*models.Organization)
	if !ok {
		panic(fmt.Sprintf("could not get organization, type: %T", org))
	}
	err := handler.useCase.Delete(ctx, org)
	if err != nil {
		httpx.ResponseJSONError(w, r, http.StatusInternalServerError, err)
		return
	}
	httpx.NoContent(w)
}
