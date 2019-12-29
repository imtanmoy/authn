package repository

import (
	"context"
	"errors"
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authn/models"
	"github.com/imtanmoy/authn/organization"
	"github.com/imtanmoy/godbx"
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
	err = godbx.ParsePgError(err)
	return organizations, err
}

func (r repository) Save(ctx context.Context, org *models.Organization) (*models.Organization, error) {
	db := r.db.WithContext(ctx)
	err := db.Insert(org)
	return org, err
}

func (r repository) Find(ctx context.Context, id int) (*models.Organization, error) {
	db := r.db.WithContext(ctx)
	if !r.Exists(ctx, id) {
		return nil, errors.New("organization does not exist")
	}
	var org models.Organization
	err := db.Model(&org).Where("id = ?", id).Select()
	return &org, err
}

func (r repository) Exists(ctx context.Context, id int) bool {
	db := r.db.WithContext(ctx)
	var num int
	_, err := db.Query(pg.Scan(&num), "SELECT id from organizations where id = ?", id)
	if err != nil {
		panic(err)
	}
	return num == id
}

func (r repository) Delete(ctx context.Context, org *models.Organization) error {
	db := r.db.WithContext(ctx)
	err := db.Delete(org)
	return err
}

func (r repository) Update(ctx context.Context, org *models.Organization) error {
	db := r.db.WithContext(ctx)
	err := db.Update(org)
	return err
}
