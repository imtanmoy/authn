package models

import (
	"context"
	"github.com/go-pg/pg/v9/orm"
	"github.com/imtanmoy/authn/internal/authx"
	"time"
)

// User represent users table
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

//BeforeInsert hooks
func (u *User) BeforeInsert(ctx context.Context) (context.Context, error) {
	now := time.Now()
	if u.CreatedAt.IsZero() {
		u.CreatedAt = now
	}
	if u.UpdatedAt.IsZero() {
		u.UpdatedAt = now
	}
	return ctx, nil
}

//BeforeUpdate hooks
func (u *User) BeforeUpdate(ctx context.Context) (context.Context, error) {
	u.UpdatedAt = time.Now()
	return ctx, nil
}

func (u *User) GetEmail() (email string) {
	return u.Email
}

func (u *User) GetPassword() (password string) {
	return u.Password
}

func (u *User) PutPassword(password string) {
	panic("implement me")
}

func (u *User) GetId() (id int) {
	return u.ID
}

func (u *User) GetEnabled() (enabled bool) {
	return u.Enabled
}

var _ orm.BeforeInsertHook = (*User)(nil)
var _ orm.BeforeUpdateHook = (*User)(nil)
var _ authx.AuthableUser = (*User)(nil)
var _ authx.AuthUser = (*User)(nil)

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
