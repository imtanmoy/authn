package repository

import (
	"context"
	"errors"
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authy/entities"
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

func (r repository) FindAll(ctx context.Context) ([]*entities.Organization, error) {
	db := r.db.WithContext(ctx)
	var organizations []*entities.Organization
	err := db.Model(&organizations).Select()
	return organizations, err
}

func (r repository) Save(ctx context.Context, org *entities.Organization) (*entities.Organization, error) {
	db := r.db.WithContext(ctx)
	err := db.Insert(org)
	return org, err
}

func (r repository) Find(ctx context.Context, id int32) (*entities.Organization, error) {
	db := r.db.WithContext(ctx)
	if !r.Exists(ctx, id) {
		return nil, errors.New("organization does not exist")
	}
	var org entities.Organization
	err := db.Model(&org).Where("id = ?", id).Select()
	return &org, err
}

func (r repository) Exists(ctx context.Context, id int32) bool {
	db := r.db.WithContext(ctx)
	var num int32
	_, err := db.Query(pg.Scan(&num), "SELECT id from organizations where id = ?", id)
	if err != nil {
		panic(err)
	}
	return num == id
}

func (r repository) Delete(ctx context.Context, org *entities.Organization) error {
	db := r.db.WithContext(ctx)
	err := db.Delete(org)
	return err
}

func (r repository) Update(ctx context.Context, org *entities.Organization) error {
	db := r.db.WithContext(ctx)
	err := db.Update(org)
	return err
}
