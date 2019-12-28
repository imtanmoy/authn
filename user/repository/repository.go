package repository

import (
	"context"
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authn/models"
	"github.com/imtanmoy/authn/user"
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
	return users, err
}

func (repo *repository) FindAllByOrganizationId(ctx context.Context, id int) ([]*models.User, error) {
	panic("implement me")
}

func (repo *repository) Save(ctx context.Context, u *models.User) error {
	panic("implement me")
}

func (repo *repository) Find(ctx context.Context, id int) (*models.User, error) {
	panic("implement me")
}

func (repo *repository) Exists(ctx context.Context, id int) bool {
	panic("implement me")
}

func (repo *repository) Delete(ctx context.Context, u *models.User) error {
	panic("implement me")
}

func (repo *repository) Update(ctx context.Context, u *models.User) error {
	panic("implement me")
}
