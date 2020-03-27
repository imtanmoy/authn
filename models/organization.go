package models

import (
	"time"
)

// Organization represent organizations table
type Organization struct {
	ID        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Users     []*User
}
