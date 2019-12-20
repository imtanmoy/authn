package presenter

import (
	"github.com/imtanmoy/authy/entities"
)

type OrganizationResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

func NewOrganizationResponse(organization *entities.Organization) *OrganizationResponse {
	resp := &OrganizationResponse{
		ID:   organization.ID,
		Name: organization.Name,
	}
	return resp
}

func NewOrganizationListResponse(organizations []*entities.Organization) []*OrganizationResponse {
	var list []*OrganizationResponse
	if len(organizations) == 0 {
		list = make([]*OrganizationResponse, 0)
	}
	for _, organization := range organizations {
		list = append(list, NewOrganizationResponse(organization))
	}
	return list
}
