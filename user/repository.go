package user

import (
	"context"
	"github.com/imtanmoy/authn/models"
)

type Repository interface {
	FindAll(ctx context.Context) ([]*models.User, error)
	FindAllByOrganizationId(ctx context.Context, id int) ([]*models.User, error)
	Save(ctx context.Context, u *models.User) error
	Find(ctx context.Context, id int) (*models.User, error)
	Exists(ctx context.Context, id int) bool
	Delete(ctx context.Context, u *models.User) error
	Update(ctx context.Context, u *models.User) error
}
