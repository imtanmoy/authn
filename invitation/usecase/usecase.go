package usecase

import (
	"context"
	"github.com/imtanmoy/authn/invitation"
	"github.com/imtanmoy/authn/models"
	"time"
)

type useCase struct {
	inviteRepo     invitation.Repository
	contextTimeout time.Duration
}

func (uc *useCase) FindByToken(ctx context.Context, token string) (*models.Invitation, error) {
	return uc.inviteRepo.FindByToken(ctx, token)
}

func (uc *useCase) FindByEmailAndOrganization(ctx context.Context, email string, oid int) (*models.Invitation, error) {
	return uc.inviteRepo.FindByEmailAndOrganization(ctx, email, oid)
}

func (uc *useCase) FindAll(ctx context.Context) ([]*models.Invitation, error) {
	panic("implement me")
}

func (uc *useCase) Store(ctx context.Context, u *models.Invitation) error {
	return uc.inviteRepo.Save(ctx, u)
}

func (uc *useCase) GetById(ctx context.Context, id int) (*models.Invitation, error) {
	panic("implement me")
}

func (uc *useCase) Update(ctx context.Context, u *models.Invitation) error {
	return uc.inviteRepo.Update(ctx, u)
}

func (uc *useCase) Delete(ctx context.Context, u *models.Invitation) error {
	panic("implement me")
}

func (uc *useCase) Exists(ctx context.Context, id int) bool {
	panic("implement me")
}

func (uc *useCase) ExistsByEmail(ctx context.Context, email string) bool {
	panic("implement me")
}

var _ invitation.UseCase = (*useCase)(nil)

// NewUseCase will create new an useCase object representation of invite.UseCase interface
func NewUseCase(g invitation.Repository, timeout time.Duration) invitation.UseCase {
	return &useCase{
		inviteRepo:     g,
		contextTimeout: timeout,
	}
}
