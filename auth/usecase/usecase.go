package usecase

import (
	"github.com/imtanmoy/authn/auth"
	"github.com/imtanmoy/authn/user"
	"time"
)

type useCase struct {
	userRepo       user.Repository
	contextTimeout time.Duration
}

//func (uc *useCase) ExistsByEmail(ctx context.Context, email string) bool {
//	return uc.userRepo.ExistsByEmail(ctx, email)
//}
//
//func (uc *useCase) FindByEmail(ctx context.Context, email string) (*models.User, error) {
//	if !uc.ExistsByEmail(ctx, email) {
//		return nil, errorx.ErrorNotFound
//	}
//	return uc.userRepo.FindByEmail(ctx, email)
//}

func NewUseCase(userRepo user.Repository, contextTimeout time.Duration) auth.UseCase {
	return &useCase{userRepo: userRepo, contextTimeout: contextTimeout}
}
