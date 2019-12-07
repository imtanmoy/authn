package usecase

import (
	"context"
	"github.com/imtanmoy/authy/models"
	"github.com/imtanmoy/authy/organization"
	"time"
)

type usecase struct {
	orgRepo        organization.Repository
	contextTimeout time.Duration
}

var _ organization.UseCase = (*usecase)(nil)

// NewUsecase will create new an usecase object representation of organization.Usecase interface
func NewUsecase(g organization.Repository, timeout time.Duration) organization.UseCase {
	return &usecase{
		orgRepo:        g,
		contextTimeout: timeout,
	}
}

func (u usecase) FindAll(ctx context.Context) ([]*models.Organization, error) {
	return u.orgRepo.FindAll(ctx)
}

func (u usecase) Store(ctx context.Context, org *models.Organization) error {
	org1, err := u.orgRepo.Store(ctx, org)
	if err != nil {
		return err
	}
	org = org1
	return nil
}

func (u usecase) GetById(ctx context.Context, id int32) (*models.Organization, error) {
	panic("implement me")
}

func (u usecase) Update(ctx context.Context, org *models.Organization) error {
	panic("implement me")
}

func (u usecase) Delete(ctx context.Context, org *models.Organization) error {
	panic("implement me")
}

func (u usecase) Exists(ctx context.Context, Id int32) bool {
	panic("implement me")
}
