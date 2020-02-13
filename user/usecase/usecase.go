package usecase

import (
	"context"
	"github.com/imtanmoy/authn/internal/errorx"
	"github.com/imtanmoy/authn/models"
	"github.com/imtanmoy/authn/user"
	"time"
)

type useCase struct {
	userRepo       user.Repository
	contextTimeout time.Duration
}

func (uc *useCase) StoreWithOrg(ctx context.Context, u *models.User, ou *models.UserOrganization) error {
	err := uc.userRepo.Save(ctx, u)
	if err != nil {
		return err
	}

	ou.UserId = u.ID
	err = uc.userRepo.SaveUserOrganization(ctx, ou)
	return err
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
	if !uc.Exists(ctx, id) {
		return nil, errorx.ErrorNotFound
	}
	return uc.userRepo.Find(ctx, id)
}

func (uc *useCase) Update(ctx context.Context, u *models.User) error {
	return uc.userRepo.Update(ctx, u)
}

func (uc *useCase) Delete(ctx context.Context, u *models.User) error {
	return uc.userRepo.Delete(ctx, u)
}

func (uc *useCase) Exists(ctx context.Context, id int) bool {
	return uc.userRepo.Exists(ctx, id)
}

func (uc *useCase) ExistsByEmail(ctx context.Context, email string) bool {
	return uc.userRepo.ExistsByEmail(ctx, email)
}
