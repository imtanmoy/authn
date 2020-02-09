package invitation

import (
	"context"
	"github.com/imtanmoy/authn/models"
)

type Repository interface {
	FindAll(ctx context.Context) ([]*models.Invitation, error)
	FindAllByOrganizationId(ctx context.Context, id int) ([]*models.Invitation, error)
	Save(ctx context.Context, u *models.Invitation) error
	Find(ctx context.Context, id int) (*models.Invitation, error)
	Update(ctx context.Context, u *models.Invitation) error
	FindByEmail(ctx context.Context, email string) (*models.Invitation, error)
	Exists(ctx context.Context, id int) bool
	ExistsByEmail(ctx context.Context, email string) bool
	Delete(ctx context.Context, u *models.Invitation) error
	FindByEmailAndOrganization(ctx context.Context, email string, oid int) (*models.Invitation, error)
	FindByToken(ctx context.Context, token string) (*models.Invitation, error)
}
