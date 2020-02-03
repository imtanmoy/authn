package models

import "github.com/go-pg/pg/v9/orm"

func init() {
	// Register many to many model so ORM can better recognize m2m relation.
	// This should be done before dependant models are used.
	orm.RegisterTable((*UserOrganization)(nil))
}

type UserOrganization struct {
	tableName struct{} `pg:"user_organization"`

	UserId         int `pg:"user_id,notnull"`
	User           *User
	OrganizationId int `pg:"organization_id,notnull"`
	Organization   *Organization
}
