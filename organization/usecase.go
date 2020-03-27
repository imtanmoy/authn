package organization

import (
	"context"
	"github.com/imtanmoy/authn/models"
)

type UseCase interface {
	Save(ctx context.Context, org *models.Organization) error
}
