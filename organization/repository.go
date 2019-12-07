package organization

import (
	"context"
	"github.com/imtanmoy/authy/models"
)

type Repository interface {
	FindAll(ctx context.Context) ([]*models.Organization, error)
	Store(ctx context.Context, org *models.Organization) (*models.Organization, error)
	Find(ctx context.Context, ID int32) (*models.Organization, error)
	Exists(ctx context.Context, ID int32) bool
	Delete(ctx context.Context, org *models.Organization) error
	Update(ctx context.Context, org *models.Organization) error
	FindAllByIdIn(ctx context.Context, ids []int32) []*models.Organization
}
