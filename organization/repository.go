package organization

import (
	"context"
	"github.com/imtanmoy/authn/models"
)

type Repository interface {
	FindAll(ctx context.Context) ([]*models.Organization, error)
	Save(ctx context.Context, org *models.Organization) (*models.Organization, error)
	Find(ctx context.Context, id int) (*models.Organization, error)
	Exists(ctx context.Context, id int) bool
	Delete(ctx context.Context, org *models.Organization) error
	Update(ctx context.Context, org *models.Organization) error
}
