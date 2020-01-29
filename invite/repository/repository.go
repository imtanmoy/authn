package repository

import (
	"context"
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authn/invite"
	"github.com/imtanmoy/authn/models"
)

type repository struct {
	db *pg.DB
}

func (r *repository) Delete(ctx context.Context, u *models.Invite) error {
	panic("implement me")
}

func (r *repository) FindAll(ctx context.Context) ([]*models.Invite, error) {
	panic("implement me")
}

func (r *repository) FindAllByOrganizationId(ctx context.Context, id int) ([]*models.Invite, error) {
	panic("implement me")
}

func (r *repository) Save(ctx context.Context, u *models.Invite) error {
	panic("implement me")
}

func (r *repository) Find(ctx context.Context, id int) (*models.Invite, error) {
	panic("implement me")
}

func (r *repository) Exists(ctx context.Context, id int) bool {
	panic("implement me")
}

func (r *repository) ExistsByEmail(ctx context.Context, email string) bool {
	panic("implement me")
}

var _ invite.Repository = (*repository)(nil)

// NewRepository will create an object that represent the invite.Repository interface
func NewRepository(db *pg.DB) invite.Repository {
	return &repository{db}
}
