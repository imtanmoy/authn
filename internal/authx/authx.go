package authx

import (
	"context"
	"github.com/imtanmoy/authn/internal/errorx"
	"github.com/imtanmoy/authn/models"
	"github.com/imtanmoy/authn/user"
)

type authBackend struct {
	userRepo              user.Repository
	secretKey             string
	accessTokenExpireTime int
}

var authenticator *authBackend

func New(userRepo user.Repository, secretKey string, accessTokenExpireTime int) *authBackend {
	return &authBackend{userRepo: userRepo, secretKey: secretKey, accessTokenExpireTime: accessTokenExpireTime}
}

func (al *authBackend) Init() {
	authenticator = al
}

func getUser(ctx context.Context, identity string) (*models.User, error) {
	if !authenticator.userRepo.ExistsByEmail(ctx, identity) {
		return nil, errorx.ErrorNotFound
	}
	u, err := authenticator.userRepo.FindByEmail(ctx, identity)
	if err != nil {
		return nil, err
	}
	return u, nil
}
