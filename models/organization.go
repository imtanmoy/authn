package models

import (
	"time"
)

// Organization represent organizations table
type Organization struct {
	ID        int
	Name      string
	OwnerID   int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	Users     []*User
}
