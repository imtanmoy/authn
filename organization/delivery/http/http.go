package http

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/imtanmoy/authy/entities"
	"github.com/imtanmoy/authy/organization"
	"github.com/imtanmoy/authy/organization/presenter"
	"github.com/imtanmoy/authy/utils/httputil"
	param "github.com/oceanicdev/chi-param"
	"net/http"
)

type contextKey string

const (
	orgKey contextKey = "organization"
)

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
			_ = render.Render(w, r, httputil.NewAPIError(500, "Something went wrong"))
			return
		}
		id, err := param.Int32(r, "id")
		if err != nil {
			_ = render.Render(w, r, httputil.NewAPIError(400, "Invalid request parameter", err))
			return
		}
		org, err := oh.useCase.GetById(ctx, id)
		if err != nil {
			_ = render.Render(w, r, httputil.NewAPIError(404, "organization not found", err))
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
		_ = render.Render(w, r, httputil.NewAPIError(err))
		return
	}

	if err := render.RenderList(w, r, presenter.NewOrganizationListResponse(organizations)); err != nil {
		_ = render.Render(w, r, httputil.NewAPIError(err))
		return
	}
}

func (oh *OrganizationHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	data := &OrganizationPayload{}
	if err := render.Bind(r, data); err != nil {
		_ = render.Render(w, r, httputil.NewAPIError(err))
		return
	}

	validationErrors := data.validate()

	if len(validationErrors) > 0 {
		_ = render.Render(w, r, httputil.NewAPIError(400, "Invalid Request", validationErrors))
		return
	}

	var org entities.Organization
	org.Name = data.Name

	err := oh.useCase.Store(ctx, &org)
	if err != nil {
		_ = render.Render(w, r, httputil.NewAPIError(err))
		return
	}

	render.Status(r, http.StatusCreated)
	_ = render.Render(w, r, presenter.NewOrganizationResponse(&org))
	return
}

func (oh *OrganizationHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	org, ok := ctx.Value(orgKey).(*entities.Organization)
	if !ok {
		_ = render.Render(w, r, httputil.NewAPIError(500, "Something went wrong"))
		return
	}
	if err := render.Render(w, r, presenter.NewOrganizationResponse(org)); err != nil {
		_ = render.Render(w, r, httputil.NewAPIError(err))
		return
	}
}

func (oh *OrganizationHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	org, ok := ctx.Value(orgKey).(*entities.Organization)
	if !ok {
		_ = render.Render(w, r, httputil.NewAPIError(500, "Something went wrong"))
		return
	}
	data := &OrganizationPayload{}
	if err := render.Bind(r, data); err != nil {
		_ = render.Render(w, r, httputil.NewAPIError(err))
		return
	}

	validationErrors := data.validate()

	if len(validationErrors) > 0 {
		_ = render.Render(w, r, httputil.NewAPIError(400, "Invalid Request", validationErrors))
		return
	}
	// update organization's data
	org.Name = data.Name
	err := oh.useCase.Update(ctx, org)
	if err != nil {
		_ = render.Render(w, r, httputil.NewAPIError(err))
		return
	}
	if err := render.Render(w, r, presenter.NewOrganizationResponse(org)); err != nil {
		_ = render.Render(w, r, httputil.NewAPIError(err))
		return
	}
}

func (oh *OrganizationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	org, ok := ctx.Value(orgKey).(*entities.Organization)
	if !ok {
		_ = render.Render(w, r, httputil.NewAPIError(500, "Something went wrong"))
		return
	}
	err := oh.useCase.Delete(ctx, org)
	if err != nil {
		_ = render.Render(w, r, httputil.NewAPIError(err))
		return
	}
	render.NoContent(w, r)
}
