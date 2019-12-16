package presenter

import (
	"github.com/go-chi/render"
	"github.com/imtanmoy/authy/entities"
	"net/http"
)

type OrganizationResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

func (o *OrganizationResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewOrganizationResponse(organization *entities.Organization) *OrganizationResponse {
	resp := &OrganizationResponse{
		ID:   organization.ID,
		Name: organization.Name,
	}
	return resp
}

func NewOrganizationListResponse(organizations []*entities.Organization) []render.Renderer {
	var list []render.Renderer
	if len(organizations) == 0 {
		list = make([]render.Renderer, 0)
	}
	for _, organization := range organizations {
		list = append(list, NewOrganizationResponse(organization))
	}
	return list
}
