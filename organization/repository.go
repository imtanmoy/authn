package organization

import (
	"context"
	"github.com/imtanmoy/authy/entities"
)

type Repository interface {
	FindAll(ctx context.Context) ([]*entities.Organization, error)
	Save(ctx context.Context, org *entities.Organization) (*entities.Organization, error)
	Find(ctx context.Context, ID int32) (*entities.Organization, error)
	Exists(ctx context.Context, ID int32) bool
	Delete(ctx context.Context, org *entities.Organization) error
	Update(ctx context.Context, org *entities.Organization) error
}
