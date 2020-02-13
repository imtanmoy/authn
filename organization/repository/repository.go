package repository

import (
	"context"
	"errors"
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authn/internal/errorx"
	"github.com/imtanmoy/authn/models"
	"github.com/imtanmoy/authn/organization"
	"github.com/imtanmoy/godbx"
)

type repository struct {
	db *pg.DB
}

func (repo *repository) FindAllByUserId(ctx context.Context, id int) ([]*models.Membership, error) {
	db := repo.db.WithContext(ctx)

	var orgs []*models.Membership
	_, err := db.Query(&orgs, `SELECT 
											"organization".id,
										    "organization".name,
											"organization".owner_id,
										    "organization".created_at,
										    "organization".updated_at,
										    "user_org".joined_at,
										    "user_org".created_by,
										    "user_org".updated_by,
										    "user_org".deleted_by,
										    "user_org".enabled
									FROM "organizations" "organization"
											 JOIN users_organizations "user_org" ON user_org.organization_id = "organization".id
											 JOIN users "u" on "user_org".user_id = "u".id
									WHERE "u".id = ?`, id)
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			orgs = make([]*models.Membership, 0)
			return orgs, nil
		} else {
			panic(err)
		}
	}
	return orgs, nil
}

func (repo *repository) FindAllUserOrganizationByOid(ctx context.Context, id int) ([]*models.UserOrganization, error) {
	db := repo.db.WithContext(ctx)
	var organizations []*models.UserOrganization
	err := db.Model(&organizations).Where("organization_id = ?", id).Select()
	return organizations, err
}

var _ organization.Repository = (*repository)(nil)

// NewRepository will create an object that represent the organization.Repository interface
func NewRepository(db *pg.DB) organization.Repository {
	return &repository{db}
}

func (repo *repository) FindAll(ctx context.Context) ([]*models.Organization, error) {
	db := repo.db.WithContext(ctx)
	var organizations []*models.Organization
	err := db.Model(&organizations).Select()
	err = godbx.ParsePgError(err)
	return organizations, err
}

func (repo *repository) Save(ctx context.Context, org *models.Organization) error {
	db := repo.db.WithContext(ctx)
	err := db.Insert(org)
	return err
}

func (repo *repository) SaveUserOrganization(ctx context.Context, orgUser *models.UserOrganization) error {
	db := repo.db.WithContext(ctx)
	err := db.Insert(orgUser)
	return err
}

func (repo *repository) Find(ctx context.Context, id int) (*models.Organization, error) {
	db := repo.db.WithContext(ctx)
	var org models.Organization
	err := db.Model(&org).Where("id = ?", id).Select()
	return &org, err
}

func (repo *repository) Exists(ctx context.Context, id int) bool {
	db := repo.db.WithContext(ctx)
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

func (repo *repository) Delete(ctx context.Context, org *models.Organization) error {
	db := repo.db.WithContext(ctx)
	err := db.Delete(org)
	return err
}

func (repo *repository) DeleteUserOrganization(ctx context.Context, orgs []*models.UserOrganization) error {
	db := repo.db.WithContext(ctx)
	for _, o := range orgs {
		err := db.Delete(o)
		return err
	}
	return nil
}

func (repo *repository) Update(ctx context.Context, org *models.Organization) error {
	db := repo.db.WithContext(ctx)
	err := db.Update(org)
	return err
}

func (repo *repository) GetMembershipById(ctx context.Context, id, uid int) (*models.Membership, error) {
	db := repo.db.WithContext(ctx)

	var orgs []*models.Membership
	//_, err := db.QueryOne(&org, `SELECT
	//											"organization".id,
	//										    "organization".name,
	//										    "organization".created_at,
	//										    "organization".updated_at,
	//										    "user_org".joined_at,
	//										    "user_org".created_by,
	//										    "user_org".updated_by,
	//										    "user_org".deleted_by,
	//										    "user_org".enabled
	//									FROM "organizations" "organization"
	//											 JOIN users_organizations "user_org" ON "user_org".organization_id = "organization".id
	//											 JOIN users "u" on "user_org".user_id = "u".id
	//									WHERE "organization".id = ?
	//									  AND "u".id = ?`, id, uid)
	_, err := db.QueryOne(&orgs, `SELECT 
											"organization".id,
										    "organization".name,
											"organization".owner_id,
										    "organization".created_at,
										    "organization".updated_at,
										    "user_org".joined_at,
										    "user_org".created_by,
										    "user_org".updated_by,
										    "user_org".deleted_by,
										    "user_org".enabled
									FROM "organizations" "organization"
											 JOIN users_organizations "user_org" ON user_org.organization_id = "organization".id
											 JOIN users "u" on "user_org".user_id = "u".id
									WHERE "u".id = ? AND "organization".id = ? LIMIT 1`, uid, id)
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, errorx.ErrorNotFound
		} else {
			panic(err)
		}
	}
	if len(orgs) == 0 {
		return nil, errorx.ErrorNotFound
	}
	return orgs[0], nil
}
