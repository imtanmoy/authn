package http

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/imtanmoy/authy/models"
	"github.com/imtanmoy/authy/organization"
	"github.com/imtanmoy/authy/utils/httputil"
	"net/http"
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
		r.Get("/", handler.FindAll)
		r.Post("/", handler.Create)
	})
}

func (oh *OrganizationHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		_ = render.Render(w, r, httputil.NewAPIError(500, "Something went wrong"))
		return
	}
	organizations, err := oh.useCase.FindAll(ctx)
	if err != nil {
		_ = render.Render(w, r, httputil.NewAPIError(err))
		return
	}

	if err := render.RenderList(w, r, NewOrganizationListResponse(organizations)); err != nil {
		_ = render.Render(w, r, httputil.NewAPIError(err))
		return
	}
}

func (oh *OrganizationHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
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

	var org models.Organization
	org.Name = data.Name

	err := oh.useCase.Store(ctx, &org)
	if err != nil {
		_ = render.Render(w, r, httputil.NewAPIError(err))
		return
	}

	render.Status(r, http.StatusCreated)
	_ = render.Render(w, r, NewOrganizationResponse(&org))
	return
}
