package repository

import (
	"context"
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authy/models"
	"github.com/imtanmoy/authy/organization"
)

type repository struct {
	db *pg.DB
}

var _ organization.Repository = (*repository)(nil)

// NewRepository will create an object that represent the organization.Repository interface
func NewRepository(db *pg.DB) organization.Repository {
	return &repository{db}
}

func (r repository) FindAll(ctx context.Context) ([]*models.Organization, error) {
	db := r.db.WithContext(ctx)
	var organizations []*models.Organization
	err := db.Model(&organizations).Select()
	return organizations, err
}

func (r repository) Store(ctx context.Context, org *models.Organization) (*models.Organization, error) {
	db := r.db.WithContext(ctx)
	err := db.Insert(org)
	return org, err
}

func (r repository) Find(ctx context.Context, ID int32) (*models.Organization, error) {
	panic("implement me")
}

func (r repository) Exists(ctx context.Context, ID int32) bool {
	panic("implement me")
}

func (r repository) Delete(ctx context.Context, org *models.Organization) error {
	panic("implement me")
}

func (r repository) Update(ctx context.Context, org *models.Organization) error {
	panic("implement me")
}

func (r repository) FindAllByIdIn(ctx context.Context, ids []int32) []*models.Organization {
	panic("implement me")
}
