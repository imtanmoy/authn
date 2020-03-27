package usecase

import (
	"github.com/imtanmoy/authn/organization"
	"time"
)

type useCase struct {
	userRepo       organization.Repository
	contextTimeout time.Duration
}

var _ organization.UseCase = (*useCase)(nil)

// NewUseCase will create new an useCase object representation of user.UseCase interface
func NewUseCase(g organization.Repository, timeout time.Duration) organization.UseCase {
	return &useCase{
		userRepo:       g,
		contextTimeout: timeout,
	}
}
