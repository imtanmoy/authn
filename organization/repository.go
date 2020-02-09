package organization

import (
	"context"
	"github.com/imtanmoy/authn/models"
)

type Repository interface {
	FindAll(ctx context.Context) ([]*models.Organization, error)
	Save(ctx context.Context, org *models.Organization) error
	SaveUserOrganization(ctx context.Context, orgUser *models.UserOrganization) error
	Find(ctx context.Context, id int) (*models.Organization, error)
	Exists(ctx context.Context, id int) bool
	Delete(ctx context.Context, org *models.Organization) error
	DeleteUserOrganization(ctx context.Context, orgs []*models.UserOrganization) error
	Update(ctx context.Context, org *models.Organization) error

	FindAllUserOrganizationByOid(ctx context.Context, id int) ([]*models.UserOrganization, error)
	FindAllByUserId(ctx context.Context, id int) ([]*models.Organization, error)
}
