package models

import (
	"context"
	"github.com/go-pg/pg/v9/orm"
	"time"
)

func init() {
	// Register many to many model so ORM can better recognize m2m relation.
	// This should be done before dependant models are used.
	orm.RegisterTable((*UserOrganization)(nil))
}

type UserOrganization struct {
	tableName struct{} `pg:"users_organizations"`

	ID             int `pg:"id,notnull,unique,pk"`
	UserId int `pg:"user_id,notnull"`
	User           *User
	OrganizationId int `pg:"organization_id,notnull"`
	Organization   *Organization
	JoinedAt       time.Time `pg:"joined_at"`
	Enabled        bool      `pg:"enabled,notnull,default:TRUE"`
	CreatedBy      int       `pg:"created_by"`
	UpdatedBy      int       `pg:"updated_by"`
	DeletedBy      int       `pg:"deleted_by"`
}

var _ orm.BeforeInsertHook = (*UserOrganization)(nil)
var _ orm.BeforeUpdateHook = (*UserOrganization)(nil)

//BeforeInsert hooks
func (o *UserOrganization) BeforeInsert(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

//BeforeUpdate hooks
func (o *UserOrganization) BeforeUpdate(ctx context.Context) (context.Context, error) {
	return ctx, nil
}
