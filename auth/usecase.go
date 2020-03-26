package auth

import (
	"context"
	"github.com/imtanmoy/authn/models"
)

// UseCase represent the auth's use cases
type UseCase interface {
	ExistsByEmail(ctx context.Context, email string) bool
	FindByEmail(ctx context.Context, email string) (*models.User, error)
}
