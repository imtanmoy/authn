package usecase

import (
	"context"
	"github.com/imtanmoy/authn/internal/errorx"
	"github.com/imtanmoy/authn/models"
	"github.com/imtanmoy/authn/organization"
	"time"
)

type useCase struct {
	orgRepo        organization.Repository
	contextTimeout time.Duration
}

var _ organization.UseCase = (*useCase)(nil)

// NewUseCase will create new an useCase object representation of organization.UseCase interface
func NewUseCase(g organization.Repository, timeout time.Duration) organization.UseCase {
	return &useCase{
		orgRepo:        g,
		contextTimeout: timeout,
	}
}

func (u *useCase) FindAll(ctx context.Context) ([]*models.Organization, error) {
	return u.orgRepo.FindAll(ctx)
}

func (u *useCase) FindAllByUserId(ctx context.Context, id int) ([]*models.Membership, error) {
	return u.orgRepo.FindAllByUserId(ctx, id)
}

func (u *useCase) Store(ctx context.Context, org *models.Organization, user *models.User) error {
	err := u.orgRepo.Save(ctx, org)
	if err != nil {
		return err
	}
	var ou models.UserOrganization
	ou.UserId = user.ID
	ou.OrganizationId = org.ID
	ou.JoinedAt = time.Now()
	err = u.orgRepo.SaveUserOrganization(ctx, &ou)
	return err
}

func (u *useCase) GetById(ctx context.Context, id int) (*models.Organization, error) {
	if !u.Exists(ctx, id) {
		return nil, errorx.ErrorNotFound
	}
	return u.orgRepo.Find(ctx, id)
}

func (u *useCase) Update(ctx context.Context, org *models.Organization) error {
	return u.orgRepo.Update(ctx, org)
}

func (u *useCase) Delete(ctx context.Context, org *models.Organization) error {
	orgs, err := u.orgRepo.FindAllUserOrganizationByOid(ctx, org.ID)
	if err != nil {
		return err
	}
	err = u.orgRepo.DeleteUserOrganization(ctx, orgs)
	if err != nil {
		return err
	}
	return u.orgRepo.Delete(ctx, org)
}

func (u *useCase) Exists(ctx context.Context, id int) bool {
	return u.orgRepo.Exists(ctx, id)
}
