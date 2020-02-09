package invitation

import (
	"context"
	"github.com/imtanmoy/authn/models"
)

// UseCase represent the invite's use cases
type UseCase interface {
	FindAll(ctx context.Context) ([]*models.Invitation, error)
	Store(ctx context.Context, u *models.Invitation) error
	GetById(ctx context.Context, id int) (*models.Invitation, error)
	Update(ctx context.Context, u *models.Invitation) error
	Delete(ctx context.Context, u *models.Invitation) error
	Exists(ctx context.Context, id int) bool
	FindByEmailAndOrganization(ctx context.Context, email string, oid int) (*models.Invitation, error)
	FindByToken(ctx context.Context, token string) (*models.Invitation, error)
}
