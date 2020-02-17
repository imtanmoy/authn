package usecase

import (
	"github.com/imtanmoy/authn/confirmation"
	"time"
)

type useCase struct {
	contextTimeout time.Duration
}

var _ confirmation.UseCase = (*useCase)(nil)

func NewUseCase(contextTimeout time.Duration) confirmation.UseCase {
	return &useCase{contextTimeout: contextTimeout}
}
