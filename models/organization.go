package models

import (
	"context"
	"github.com/go-pg/pg/v9/orm"
	"time"
)

// Organization represent organizations table
type Organization struct {
	ID        int    `pg:"id,notnull,unique,pk"`
	Name      string `pg:"name,notnull"`
	OwnerId   int    `pg:"owner_id,notnull"`
	Owner     *User
	CreatedAt time.Time `pg:"created_at,notnull,default:now()"`
	UpdatedAt time.Time `pg:"updated_at,notnull,default:now()"`
	Users     []*User   `pg:"many2many:users_organizations"`
}

var _ orm.BeforeInsertHook = (*Organization)(nil)
var _ orm.BeforeUpdateHook = (*Organization)(nil)

//BeforeInsert hooks
func (o *Organization) BeforeInsert(ctx context.Context) (context.Context, error) {
	now := time.Now()
	if o.CreatedAt.IsZero() {
		o.CreatedAt = now
	}
	if o.UpdatedAt.IsZero() {
		o.UpdatedAt = now
	}
	return ctx, nil
}

//BeforeUpdate hooks
func (o *Organization) BeforeUpdate(ctx context.Context) (context.Context, error) {
	o.UpdatedAt = time.Now()
	return ctx, nil
}

type Membership struct {
	tableName struct{}   `pg:",discard_unknown_columns"`
	ID        int        `pg:"id" json:"id"`
	Name      string     `pg:"name" json:"name"`
	OwnerId   int        `pg:"owner_id,notnull" json:"owner_id"`
	CreatedAt *time.Time `pg:"created_at" json:"created_at"`
	UpdatedAt *time.Time `pg:"updated_at" json:"updated_at"`
	JoinedAt  *time.Time `pg:"joined_at" json:"joined_at"`
	Enabled   bool       `pg:"enabled" json:"enabled"`
	CreatedBy *int       `pg:"created_by" json:"created_by"`
	UpdatedBy *int       `pg:"updated_by" json:"updated_by"`
	DeletedBy *int       `pg:"deleted_by" json:"deleted_by"`
}
