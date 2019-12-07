package http

import (
	"github.com/go-chi/render"
	"github.com/imtanmoy/authy/models"
	"net/http"
)

type OrganizationResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

func (o *OrganizationResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewOrganizationResponse(organization *models.Organization) *OrganizationResponse {
	resp := &OrganizationResponse{
		ID:   organization.ID,
		Name: organization.Name,
	}
	return resp
}

func NewOrganizationListResponse(organizations []*models.Organization) []render.Renderer {
	var list []render.Renderer
	for _, organization := range organizations {
		list = append(list, NewOrganizationResponse(organization))
	}
	return list
}
