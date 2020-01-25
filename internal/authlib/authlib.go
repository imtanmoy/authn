package authlib

import (
	"context"
	"errors"
	"github.com/imtanmoy/authn/internal/errorx"
	"github.com/imtanmoy/authn/models"
	"github.com/imtanmoy/authn/user"
)

type authLib struct {
	userRepo              user.Repository
	secretKey             string
	accessTokenExpireTime int
}

var auth *authLib

func NewAuthLib(userRepo user.Repository, secretKey string, accessTokenExpireTime int) *authLib {
	return &authLib{userRepo: userRepo, secretKey: secretKey, accessTokenExpireTime: accessTokenExpireTime}
}

func (al *authLib) Init() {
	auth = al
}

func (al *authLib) GetUser(ctx context.Context, identity string) (*models.User, error) {
	if al == nil {
		return nil, errors.New("authlib is not initiated")
	}
	if !al.userRepo.ExistsByEmail(ctx, identity) {
		return nil, errorx.ErrorNotFound
	}
	u, err := al.userRepo.FindByEmail(ctx, identity)
	if err != nil {
		return nil, err
	}
	return u, nil
}
