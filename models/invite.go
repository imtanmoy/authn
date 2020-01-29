package models

import (
	"context"
	"github.com/go-pg/pg/v9/orm"
	"time"
)

// Invite represent invites table
type Invite struct {
	ID             int       `pg:"id,notnull,unique,pk"`
	Email          string    `pg:"email,notnull"`
	Token          string    `pg:"token,notnull"`
	Status         string    `pg:"status,notnull" sql:"type:invite_status"`
	OrganizationId int       `pg:"organization_id,notnull"`
	UserId         int       `pg:"user_id"`
	InvitedBy      int       `pg:"invited_by,notnull"`
	CreatedAt      time.Time `pg:"created_at,notnull,default:now()"`
	UpdatedAt      time.Time `pg:"updated_at,notnull,default:now()"`
}

func (i *Invite) BeforeUpdate(ctx context.Context) (context.Context, error) {
	i.UpdatedAt = time.Now()
	return ctx, nil
}

func (i *Invite) BeforeInsert(ctx context.Context) (context.Context, error) {
	now := time.Now()
	if i.CreatedAt.IsZero() {
		i.CreatedAt = now
	}
	if i.UpdatedAt.IsZero() {
		i.UpdatedAt = now
	}
	return ctx, nil
}

var _ orm.BeforeInsertHook = (*Invite)(nil)
var _ orm.BeforeUpdateHook = (*Invite)(nil)
