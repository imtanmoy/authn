package usecase

import (
	"context"
	"github.com/imtanmoy/authn/models"
	"github.com/imtanmoy/authn/organization"
	"time"
)

type useCase struct {
	repo           organization.Repository
	contextTimeout time.Duration
}

func (u *useCase) FindByID(ctx context.Context, id int) (*models.Organization, error) {
	return u.repo.FindByID(ctx, id)
}

func (u *useCase) Save(ctx context.Context, org *models.Organization) error {
	return u.repo.Save(ctx, org)
}

var _ organization.UseCase = (*useCase)(nil)

// NewUseCase will create new an useCase object representation of user.UseCase interface
func NewUseCase(repo organization.Repository, timeout time.Duration) organization.UseCase {
	return &useCase{
		repo:           repo,
		contextTimeout: timeout,
	}
}
