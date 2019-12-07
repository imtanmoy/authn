package organization

import (
	"context"
	"github.com/imtanmoy/authy/models"
)

// UseCase represent the groups's use cases
type UseCase interface {
	FindAll(ctx context.Context) ([]*models.Organization, error)
	Store(ctx context.Context, org *models.Organization) error
	GetById(ctx context.Context, id int32) (*models.Organization, error)
	Update(ctx context.Context, org *models.Organization) error
	Delete(ctx context.Context, org *models.Organization) error
	Exists(ctx context.Context, Id int32) bool
}
