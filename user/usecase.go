package user

import (
	"context"
	"github.com/imtanmoy/authn/models"
)

// UseCase represent the user's use cases
type UseCase interface {
	FindAll(ctx context.Context) ([]*models.User, error)
	Store(ctx context.Context, u *models.User) error
	StoreWithOrg(ctx context.Context, u *models.User, org *models.UserOrganization) error
	GetById(ctx context.Context, id int) (*models.User, error)
	Update(ctx context.Context, u *models.User) error
	Delete(ctx context.Context, u *models.User) error
	Exists(ctx context.Context, id int) bool
	ExistsByEmail(ctx context.Context, email string) bool
}
