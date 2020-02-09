package organization

import (
	"context"
	"github.com/imtanmoy/authn/models"
)

// UseCase represent the organization's use cases
type UseCase interface {
	FindAll(ctx context.Context) ([]*models.Organization, error)
	Store(ctx context.Context, org *models.Organization, u *models.User) error
	GetById(ctx context.Context, id int) (*models.Organization, error)
	Update(ctx context.Context, org *models.Organization) error
	Delete(ctx context.Context, org *models.Organization) error
	Exists(ctx context.Context, id int) bool
	FindAllByUserId(ctx context.Context, id int) ([]*models.Organization, error)
}
