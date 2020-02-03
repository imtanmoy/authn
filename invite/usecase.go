package invite

import (
	"context"
	"github.com/imtanmoy/authn/models"
)

// UseCase represent the invite's use cases
type UseCase interface {
	FindAll(ctx context.Context) ([]*models.Invite, error)
	Store(ctx context.Context, u *models.Invite) error
	GetById(ctx context.Context, id int) (*models.Invite, error)
	Update(ctx context.Context, u *models.Invite) error
	Delete(ctx context.Context, u *models.Invite) error
	Exists(ctx context.Context, id int) bool
	FindByEmailAndOrganization(ctx context.Context, email string, oid int) (*models.Invite, error)
	FindByToken(ctx context.Context, token string) (*models.Invite, error)
}
