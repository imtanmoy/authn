package user

import (
	"context"
	"github.com/imtanmoy/authn/internal/authx"
	"github.com/imtanmoy/authn/models"
)

type Repository interface {
	FindAll(ctx context.Context) ([]*models.User, error)
	//FindAllByOrganizationId(ctx context.Context, id int) ([]*models.User, error)
	Save(ctx context.Context, u *models.User) error
	//SaveUserOrganization(ctx context.Context, orgUser *models.UserOrganization) error
	ExistsByID(ctx context.Context, id int) bool
	ExistsByEmail(ctx context.Context, email string) bool
	////Delete(ctx context.Context, u *models.User) error
	////Update(ctx context.Context, u *models.User) error
	FindByID(ctx context.Context, id int) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	GetByEmail(ctx context.Context, identity string) (authx.AuthUser, error)
}
