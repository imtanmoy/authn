package models

import (
	"time"
)

// User represent users table
type User struct {
	ID        int
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
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
