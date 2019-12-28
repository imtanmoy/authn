package usecase

import (
	"context"
	"github.com/imtanmoy/authn/models"
	"github.com/imtanmoy/authn/user"
	"time"
)

type useCase struct {
	userRepo       user.Repository
	contextTimeout time.Duration
}

var _ user.UseCase = (*useCase)(nil)

// NewUseCase will create new an useCase object representation of user.UseCase interface
func NewUseCase(g user.Repository, timeout time.Duration) user.UseCase {
	return &useCase{
		userRepo:       g,
		contextTimeout: timeout,
	}
}

func (uc *useCase) FindAll(ctx context.Context) ([]*models.User, error) {
	return uc.userRepo.FindAll(ctx)
}

func (uc *useCase) Store(ctx context.Context, u *models.User) error {
	return uc.userRepo.Save(ctx, u)
}

func (uc *useCase) GetById(ctx context.Context, id int) (*models.User, error) {
	panic("implement me")
}

func (uc *useCase) Update(ctx context.Context, u *models.User) error {
	panic("implement me")
}

func (uc *useCase) Delete(ctx context.Context, u *models.User) error {
	panic("implement me")
}

func (uc *useCase) Exists(ctx context.Context, id int) bool {
	panic("implement me")
}