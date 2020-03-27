package models

import (
	"time"
)

// Invitation represent invites table
type Invitation struct {
	ID             int
	Email          string
	Token          string
	Status         string
	OrganizationId int
	UserId         int
	InvitedBy      int
	AcceptedAt     time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
