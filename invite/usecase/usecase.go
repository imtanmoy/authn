package usecase

import (
	"context"
	"github.com/imtanmoy/authn/invite"
	"github.com/imtanmoy/authn/models"
	"time"
)

type useCase struct {
	inviteRepo     invite.Repository
	contextTimeout time.Duration
}

func (uc *useCase) FindByToken(ctx context.Context, token string) (*models.Invite, error) {
	return uc.inviteRepo.FindByToken(ctx, token)
}

func (uc *useCase) FindByEmailAndOrganization(ctx context.Context, email string, oid int) (*models.Invite, error) {
	return uc.inviteRepo.FindByEmailAndOrganization(ctx, email, oid)
}

func (uc *useCase) FindAll(ctx context.Context) ([]*models.Invite, error) {
	panic("implement me")
}

func (uc *useCase) Store(ctx context.Context, u *models.Invite) error {
	return uc.inviteRepo.Save(ctx, u)
}

func (uc *useCase) GetById(ctx context.Context, id int) (*models.Invite, error) {
	panic("implement me")
}

func (uc *useCase) Update(ctx context.Context, u *models.Invite) error {
	return uc.inviteRepo.Update(ctx, u)
}

func (uc *useCase) Delete(ctx context.Context, u *models.Invite) error {
	panic("implement me")
}

func (uc *useCase) Exists(ctx context.Context, id int) bool {
	panic("implement me")
}

func (uc *useCase) ExistsByEmail(ctx context.Context, email string) bool {
	panic("implement me")
}

var _ invite.UseCase = (*useCase)(nil)

// NewUseCase will create new an useCase object representation of invite.UseCase interface
func NewUseCase(g invite.Repository, timeout time.Duration) invite.UseCase {
	return &useCase{
		inviteRepo:     g,
		contextTimeout: timeout,
	}
}
