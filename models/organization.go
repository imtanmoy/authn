package models

import (
	"context"
	"github.com/go-pg/pg/v9/orm"
	"time"
)

// Organization represent organizations table
type Organization struct {
	ID        int32     `pg:"id,notnull,unique"`
	Name      string    `pg:"name,notnull"`
	CreatedAt time.Time `pg:"created_at,notnull,default:now()"`
	UpdatedAt time.Time `pg:"updated_at,notnull,default:now()"`
	//Users []*User `pg:"fk:organization_id"`
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
