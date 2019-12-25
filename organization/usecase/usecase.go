package usecase

import (
	"context"
	"github.com/imtanmoy/authN/entities"
	"github.com/imtanmoy/authN/organization"
	"time"
)

type useCase struct {
	orgRepo        organization.Repository
	contextTimeout time.Duration
}

var _ organization.UseCase = (*useCase)(nil)

// NewUseCase will create new an usecase object representation of organization.Usecase interface
func NewUseCase(g organization.Repository, timeout time.Duration) organization.UseCase {
	return &useCase{
		orgRepo:        g,
		contextTimeout: timeout,
	}
}

func (u *useCase) FindAll(ctx context.Context) ([]*entities.Organization, error) {
	return u.orgRepo.FindAll(ctx)
}

func (u *useCase) Store(ctx context.Context, org *entities.Organization) error {
	org1, err := u.orgRepo.Save(ctx, org)
	if err != nil {
		return err
	}
	org = org1
	return nil
}

func (u *useCase) GetById(ctx context.Context, id int32) (*entities.Organization, error) {
	return u.orgRepo.Find(ctx, id)
}

func (u *useCase) Update(ctx context.Context, org *entities.Organization) error {
	return u.orgRepo.Update(ctx, org)
}

func (u *useCase) Delete(ctx context.Context, org *entities.Organization) error {
	return u.orgRepo.Delete(ctx, org)
}

func (u *useCase) Exists(ctx context.Context, id int32) bool {
	return u.orgRepo.Exists(ctx, id)
}
