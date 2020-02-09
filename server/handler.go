package server

import (
	_authDeliveryHttp "github.com/imtanmoy/authn/auth/delivery/http"
	_authUseCase "github.com/imtanmoy/authn/auth/usecase"
	"github.com/imtanmoy/authn/config"
	"github.com/imtanmoy/authn/db"
	"github.com/imtanmoy/authn/internal/authx"
	_invitationDeliveryHttp "github.com/imtanmoy/authn/invitation/delivery/http"
	_inviteRepo "github.com/imtanmoy/authn/invitation/repository"
	_inviteUseCase "github.com/imtanmoy/authn/invitation/usecase"
	_orgDeliveryHttp "github.com/imtanmoy/authn/organization/delivery/http"
	_orgRepo "github.com/imtanmoy/authn/organization/repository"
	_orgUseCase "github.com/imtanmoy/authn/organization/usecase"
	_userDeliveryHttp "github.com/imtanmoy/authn/user/delivery/http"
	_userRepo "github.com/imtanmoy/authn/user/repository"
	_userUseCase "github.com/imtanmoy/authn/user/usecase"
	"time"

	"github.com/go-chi/chi"
	_chiMiddleware "github.com/go-chi/chi/middleware"
)

// New configures application resources and routes.
func New() (*chi.Mux, error) {

	r := chi.NewRouter()
	r.Use(_chiMiddleware.Recoverer)
	r.Use(_chiMiddleware.RequestID)
	r.Use(_chiMiddleware.RealIP)
	r.Use(_chiMiddleware.DefaultCompress)
	r.Use(_chiMiddleware.Timeout(15 * time.Second))
	r.Use(_chiMiddleware.Logger)
	r.Use(_chiMiddleware.AllowContentType("application/json"))
	r.Use(_chiMiddleware.Heartbeat("/heartbeat"))
	//r.Use(render.SetContentType(3)) //render.ContentTypeJSON resolve value 3

	timeoutContext := 30 * time.Millisecond * time.Second //TODO it will come from config

	orgRepo := _orgRepo.NewRepository(db.DB)
	userRepo := _userRepo.NewRepository(db.DB)
	inviteRepo := _inviteRepo.NewRepository(db.DB)

	authxConfig := authx.AuthxConfig{
		SecretKey:             config.Conf.JWT_SECRET_KEY,
		AccessTokenExpireTime: config.Conf.JWT_ACCESS_TOKEN_EXPIRES,
	}

	au := authx.New(userRepo, &authxConfig)

	orgUseCase := _orgUseCase.NewUseCase(orgRepo, timeoutContext)
	userUseCase := _userUseCase.NewUseCase(userRepo, timeoutContext)
	authUseCase := _authUseCase.NewUseCase(userRepo, timeoutContext)
	invitationUseCase := _inviteUseCase.NewUseCase(inviteRepo, timeoutContext)

	_orgDeliveryHttp.NewHandler(r, orgUseCase, au)
	_userDeliveryHttp.NewHandler(r, userUseCase, orgUseCase, au)
	_authDeliveryHttp.NewHandler(r, authUseCase, userUseCase, au)
	_invitationDeliveryHttp.NewHandler(r, invitationUseCase, userUseCase, orgUseCase, au)

	return r, nil
}
