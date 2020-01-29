package repository

import (
	"context"
	"errors"
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authn/internal/errorx"
	"github.com/imtanmoy/authn/invite"
	"github.com/imtanmoy/authn/models"
	"github.com/imtanmoy/godbx"
)

type repository struct {
	db *pg.DB
}

func (repo *repository) FindByEmailAndOrganization(ctx context.Context, email string, oid int) (*models.Invite, error) {
	db := repo.db.WithContext(ctx)
	u := new(models.Invite)
	err := db.Model(u).Where("email = ?", email).Where("organization_id = ?", oid).Select()
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, errorx.ErrorNotFound
		} else {
			panic(err)
		}
	}
	return u, nil
}

func (repo *repository) Update(ctx context.Context, u *models.Invite) error {
	db := repo.db.WithContext(ctx)
	err := db.Update(u)
	err = godbx.ParsePgError(err)
	return err
}

func (repo *repository) FindByEmail(ctx context.Context, email string) (*models.Invite, error) {
	db := repo.db.WithContext(ctx)
	var iv models.Invite
	err := db.Model(&iv).Where("email = ?", email).Select()
	err = godbx.ParsePgError(err)
	return &iv, err
}

func (repo *repository) Delete(ctx context.Context, u *models.Invite) error {
	db := repo.db.WithContext(ctx)
	err := db.Delete(u)
	err = godbx.ParsePgError(err)
	return err
}

func (repo *repository) FindAll(ctx context.Context) ([]*models.Invite, error) {
	db := repo.db.WithContext(ctx)
	var invites []*models.Invite
	err := db.Model(&invites).Select()
	if err != nil {
		_, ok := err.(pg.Error)
		if ok {
			return nil, errorx.ErrInternalDB
		} else {
			return nil, errorx.ErrInternalServer
		}
	}
	return invites, err
}

func (repo *repository) FindAllByOrganizationId(ctx context.Context, id int) ([]*models.Invite, error) {
	db := repo.db.WithContext(ctx)
	var invites []*models.Invite
	err := db.Model(&invites).Where("organization_id = ?", id).Select()
	err = godbx.ParsePgError(err)
	return invites, err
}

func (repo *repository) Save(ctx context.Context, u *models.Invite) error {
	db := repo.db.WithContext(ctx)
	err := db.Insert(u)
	return err
}

func (repo *repository) Find(ctx context.Context, id int) (*models.Invite, error) {
	db := repo.db.WithContext(ctx)
	var i models.Invite
	err := db.Model(&i).Where("id = ?", id).Select()
	err = godbx.ParsePgError(err)
	return &i, err
}

func (repo *repository) Exists(ctx context.Context, id int) bool {
	db := repo.db.WithContext(ctx)
	u := new(models.Invite)
	err := db.Model(u).Where("id = ?", id).Select()
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return false
		} else {
			panic(err)
		}
	}
	return u.ID == id
}

func (repo *repository) ExistsByEmail(ctx context.Context, email string) bool {
	db := repo.db.WithContext(ctx)
	u := new(models.Invite)
	err := db.Model(u).Where("email = ?", email).Select()
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return false
		} else {
			panic(err)
		}
	}
	return u.Email == email
}

var _ invite.Repository = (*repository)(nil)

// NewRepository will create an object that represent the invite.Repository interface
func NewRepository(db *pg.DB) invite.Repository {
	return &repository{db}
}
