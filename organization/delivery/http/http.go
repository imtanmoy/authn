package http

import (
	"context"
	"errors"
	"github.com/go-chi/chi"
	"github.com/imtanmoy/authy/entities"
	"github.com/imtanmoy/authy/organization"
	"github.com/imtanmoy/authy/organization/presenter"
	"github.com/imtanmoy/httpx"
	param "github.com/oceanicdev/chi-param"
	"gopkg.in/thedevsaddam/govalidator.v1"
	"net/http"
	"net/url"
)

type contextKey string

const (
	orgKey contextKey = "organization"
)

type OrganizationPayload struct {
	Name string `json:"name"`
}

func (o *OrganizationPayload) validate() url.Values {
	rules := govalidator.MapData{
		"name": []string{"required", "min:4", "max:20"},
	}
	opts := govalidator.Options{
		Data:  o,
		Rules: rules,
	}

	v := govalidator.New(opts)
	e := v.ValidateStruct()
	return e
}

// OrganizationHandler  represent the http handler for organization
type OrganizationHandler struct {
	useCase organization.UseCase
}

// NewHandler will initialize the organization's resources endpoint
func NewHandler(r *chi.Mux, useCase organization.UseCase) {
	handler := &OrganizationHandler{
		useCase: useCase,
	}
	r.Route("/organizations", func(r chi.Router) {
		r.Get("/", handler.List)
		r.Post("/", handler.Create)
		r.Group(func(r chi.Router) {
			r.Use(handler.OrganizationCtx)
			r.Get("/{id}", handler.Get)
			r.Put("/{id}", handler.Update)
			r.Delete("/{id}", handler.Delete)
		})
	})
}

func (oh *OrganizationHandler) OrganizationCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if ctx == nil {
			httpx.ResponseJSONError(w, r, http.StatusInternalServerError, httpx.ErrInternalServerError)
			return
		}
		id, err := param.Int32(r, "id")
		if err != nil {
			httpx.ResponseJSONError(w, r, http.StatusBadRequest, "Invalid request parameter", err)
			return
		}
		org, err := oh.useCase.GetById(ctx, id)
		if err != nil {
			httpx.ResponseJSONError(w, r, http.StatusNotFound, "organization not found", err)
			return
		}
		ctx = context.WithValue(r.Context(), orgKey, org)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (oh *OrganizationHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	organizations, err := oh.useCase.FindAll(ctx)
	if err != nil {
		httpx.ResponseJSONError(w, r, http.StatusInternalServerError, err)
		return
	}

	httpx.ResponseJSON(w, http.StatusOK, presenter.NewOrganizationListResponse(organizations))
	return
}

func (oh *OrganizationHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	data := &OrganizationPayload{}
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

	var org entities.Organization
	org.Name = data.Name

	err := oh.useCase.Store(ctx, &org)
	if err != nil {
		httpx.ResponseJSONError(w, r, http.StatusInternalServerError, err)
		return
	}
	httpx.ResponseJSON(w, http.StatusCreated, presenter.NewOrganizationResponse(&org))
	return
}

func (oh *OrganizationHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	org, ok := ctx.Value(orgKey).(*entities.Organization)
	if !ok {
		httpx.ResponseJSONError(w, r, http.StatusInternalServerError, httpx.ErrInternalServerError)
		return
	}
	httpx.ResponseJSON(w, http.StatusOK, presenter.NewOrganizationResponse(org))
	return
}

func (oh *OrganizationHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	org, ok := ctx.Value(orgKey).(*entities.Organization)
	if !ok {
		httpx.ResponseJSONError(w, r, http.StatusInternalServerError, httpx.ErrInternalServerError)
		return
	}
	data := &OrganizationPayload{}
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
	err := oh.useCase.Update(ctx, org)
	if err != nil {
		httpx.ResponseJSONError(w, r, http.StatusInternalServerError, err)
		return
	}
	httpx.ResponseJSON(w, http.StatusOK, presenter.NewOrganizationResponse(org))
	return
}

func (oh *OrganizationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	org, ok := ctx.Value(orgKey).(*entities.Organization)
	if !ok {
		httpx.ResponseJSONError(w, r, http.StatusInternalServerError, httpx.ErrInternalServerError)
		return
	}
	err := oh.useCase.Delete(ctx, org)
	if err != nil {
		httpx.ResponseJSONError(w, r, http.StatusInternalServerError, err)
		return
	}
	httpx.NoContent(w)
}
