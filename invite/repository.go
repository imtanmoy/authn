package invite

import (
	"context"
	"github.com/imtanmoy/authn/models"
)

type Repository interface {
	FindAll(ctx context.Context) ([]*models.Invite, error)
	FindAllByOrganizationId(ctx context.Context, id int) ([]*models.Invite, error)
	Save(ctx context.Context, u *models.Invite) error
	Find(ctx context.Context, id int) (*models.Invite, error)
	Update(ctx context.Context, u *models.Invite) error
	FindByEmail(ctx context.Context, email string) (*models.Invite, error)
	Exists(ctx context.Context, id int) bool
	ExistsByEmail(ctx context.Context, email string) bool
	Delete(ctx context.Context, u *models.Invite) error
	FindByEmailAndOrganization(ctx context.Context, email string, oid int) (*models.Invite, error)
}
