package organization

import (
	"context"
	"github.com/imtanmoy/authn/models"
)

type Repository interface {
	Save(ctx context.Context, org *models.Organization) error
}
