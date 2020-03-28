package organization

import (
	"context"
	"github.com/imtanmoy/authn/models"
)

type UseCase interface {
	Save(ctx context.Context, org *models.Organization) error
	FindByID(ctx context.Context, id int) (*models.Organization, error)
}
