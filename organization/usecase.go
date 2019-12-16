package organization

import (
	"context"
	"github.com/imtanmoy/authy/entities"
)

// UseCase represent the groups's use cases
type UseCase interface {
	FindAll(ctx context.Context) ([]*entities.Organization, error)
	Store(ctx context.Context, org *entities.Organization) error
	GetById(ctx context.Context, id int32) (*entities.Organization, error)
	Update(ctx context.Context, org *entities.Organization) error
	Delete(ctx context.Context, org *entities.Organization) error
	Exists(ctx context.Context, Id int32) bool
}
