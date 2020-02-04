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

func (r *repository) FindAllUserOrganizationByOid(ctx context.Context, id int) ([]*models.UserOrganization, error) {
	db := r.db.WithContext(ctx)
	var organizations []*models.UserOrganization
	err := db.Model(&organizations).Where("organization_id = ?", id).Select()
	return organizations, err
}

var _ organization.Repository = (*repository)(nil)

// NewRepository will create an object that represent the organization.Repository interface
func NewRepository(db *pg.DB) organization.Repository {
	return &repository{db}
}

func (r *repository) FindAll(ctx context.Context) ([]*models.Organization, error) {
	db := r.db.WithContext(ctx)
	var organizations []*models.Organization
	err := db.Model(&organizations).Select()
	err = godbx.ParsePgError(err)
	return organizations, err
}

func (r *repository) Save(ctx context.Context, org *models.Organization) error {
	db := r.db.WithContext(ctx)
	err := db.Insert(org)
	return err
}

func (r *repository) SaveUserOrganization(ctx context.Context, orgUser *models.UserOrganization) error {
	db := r.db.WithContext(ctx)
	err := db.Insert(orgUser)
	return err
}

func (r *repository) Find(ctx context.Context, id int) (*models.Organization, error) {
	db := r.db.WithContext(ctx)
	var org models.Organization
	err := db.Model(&org).Where("id = ?", id).Select()
	return &org, err
}

func (r *repository) Exists(ctx context.Context, id int) bool {
	db := r.db.WithContext(ctx)
	o := new(models.Organization)
	err := db.Model(o).Where("id = ?", id).Select()
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return false
		} else {
			panic(err)
		}
	}
	return o.ID == id
}

func (r *repository) Delete(ctx context.Context, org *models.Organization) error {
	db := r.db.WithContext(ctx)
	err := db.Delete(org)
	return err
}

func (r *repository) DeleteUserOrganization(ctx context.Context, orgs []*models.UserOrganization) error {
	db := r.db.WithContext(ctx)
	for _, o := range orgs {
		err := db.Delete(o)
		return err
	}
	return nil
}

func (r *repository) Update(ctx context.Context, org *models.Organization) error {
	db := r.db.WithContext(ctx)
	err := db.Update(org)
	return err
}
