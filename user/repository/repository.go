package repository

import (
	"context"
	"errors"
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authn/internal/errorx"
	"github.com/imtanmoy/authn/models"
	"github.com/imtanmoy/authn/user"
	"github.com/imtanmoy/godbx"
)

type repository struct {
	db *pg.DB
}

var _ user.Repository = (*repository)(nil)

// NewRepository will create an object that represent the user.Repository interface
func NewRepository(db *pg.DB) user.Repository {
	return &repository{db}
}

func (repo *repository) FindAll(ctx context.Context) ([]*models.User, error) {
	db := repo.db.WithContext(ctx)
	var users []*models.User
	err := db.Model(&users).Select()
	if err != nil {
		_, ok := err.(pg.Error)
		if ok {
			return nil, errorx.ErrInternalDB
		} else {
			return nil, errorx.ErrInternalServer
		}
	}
	return users, err
}

func (repo *repository) FindAllByOrganizationId(ctx context.Context, id int) ([]*models.User, error) {
	db := repo.db.WithContext(ctx)
	var users []*models.User
	err := db.Model(&users).Where("organization_id = ?", id).Select()
	err = godbx.ParsePgError(err)
	return users, err
}

func (repo *repository) Save(ctx context.Context, u *models.User) error {
	db := repo.db.WithContext(ctx)
	err := db.Insert(u)
	err = godbx.ParsePgError(err)
	return err
}

func (repo *repository) Find(ctx context.Context, id int) (*models.User, error) {
	db := repo.db.WithContext(ctx)
	var u models.User
	err := db.Model(&u).Where("id = ?", id).Select()
	err = godbx.ParsePgError(err)
	return &u, err
}

func (repo *repository) Exists(ctx context.Context, id int) bool {
	db := repo.db.WithContext(ctx)
	var u *models.User
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
	var u *models.User
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

func (repo *repository) Delete(ctx context.Context, u *models.User) error {
	db := repo.db.WithContext(ctx)
	err := db.Delete(u)
	err = godbx.ParsePgError(err)
	return err
}

func (repo *repository) Update(ctx context.Context, u *models.User) error {
	db := repo.db.WithContext(ctx)
	err := db.Update(u)
	err = godbx.ParsePgError(err)
	return err
}
