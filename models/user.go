package models

import (
	"context"
	"github.com/go-pg/pg/v9/orm"
	"time"
)

// Organization represent organizations table
type User struct {
	ID             int       `pg:"id,notnull,unique,pk"`
	Name           string    `pg:"name,notnull"`
	Designation    string    `pg:"designation"`
	Email          string    `pg:"email,notnull,unique:idx_email_deleted_at"`
	Password       string    `pg:"password"`
	Enabled        bool      `pg:"enabled,notnull,default:TRUE"`
	OrganizationId int       `pg:"organization_id,notnull"`
	CreatedBy      int       `pg:"created_by,notnull"`
	UpdatedBy      int       `pg:"updated_by,notnull"`
	DeletedBy      int       `pg:"deleted_by,notnull"`
	JoinedAt       time.Time `pg:"joined_at"`
	CreatedAt      time.Time `pg:"created_at,notnull,default:now()"`
	UpdatedAt      time.Time `pg:"updated_at,notnull,default:now()"`
	DeletedAt      time.Time `pg:"deleted_at,soft_delete"`
}

var _ orm.BeforeInsertHook = (*User)(nil)
var _ orm.BeforeUpdateHook = (*User)(nil)

//BeforeInsert hooks
func (o *User) BeforeInsert(ctx context.Context) (context.Context, error) {
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
func (o *User) BeforeUpdate(ctx context.Context) (context.Context, error) {
	o.UpdatedAt = time.Now()
	return ctx, nil
}

type UserResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Designation string `json:"designation"`
	Email       string `json:"email"`
}

func NewUserResponse(u *User) *UserResponse {
	resp := &UserResponse{
		ID:          u.ID,
		Name:        u.Name,
		Designation: u.Designation,
		Email:       u.Email,
	}
	return resp
}

func NewUserListResponse(users []*User) []*UserResponse {
	var list []*UserResponse
	if len(users) == 0 {
		list = make([]*UserResponse, 0)
	}
	for _, u := range users {
		list = append(list, NewUserResponse(u))
	}
	return list
}
