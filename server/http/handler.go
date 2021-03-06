package http

import (
	_authDeliveryHttp "github.com/imtanmoy/authn/auth/delivery/http"
	_authUseCase "github.com/imtanmoy/authn/auth/usecase"
	"github.com/imtanmoy/authn/config"
	"github.com/imtanmoy/authn/internal/authx"
	_orgDeliveryHttp "github.com/imtanmoy/authn/organization/delivery/http"
	_orgRepo "github.com/imtanmoy/authn/organization/repository"
	_orgUseCase "github.com/imtanmoy/authn/organization/usecase"
	"github.com/imtanmoy/authn/registry"
	_userRepo "github.com/imtanmoy/authn/user/repository"
	_userUseCase "github.com/imtanmoy/authn/user/usecase"
	"github.com/jackc/pgx/v4/stdlib"
	"log"
	"time"

	"github.com/go-chi/chi"
)

// RegisterHandler configures application resources and routes.
func RegisterHandler(r *chi.Mux, rg registry.Registry) {

	b := rg.Bus()

	timeoutContext := 30 * time.Millisecond * time.Second //TODO it will come from config

	conn, err := stdlib.AcquireConn(rg.DB())
	if err != nil {
		log.Fatal(err)
	}

	orgRepo := _orgRepo.NewPgxRepository(conn)
	userRepo := _userRepo.NewPgxRepository(conn)
	//inviteRepo := _inviteRepo.NewRepository(rg.DB())

	authxConfig := authx.AuthxConfig{
		SecretKey:             config.Conf.JwtSecretKey,
		AccessTokenExpireTime: config.Conf.JwtAccessTokenExpires,
	}

	au := authx.New(userRepo, &authxConfig)

	orgUseCase := _orgUseCase.NewUseCase(orgRepo, timeoutContext)
	userUseCase := _userUseCase.NewUseCase(userRepo, timeoutContext)
	authUseCase := _authUseCase.NewUseCase(userRepo, timeoutContext)
	//invitationUseCase := _inviteUseCase.NewUseCase(inviteRepo, timeoutContext)
	//confirmationUseCase := _confirmationUseCase.NewUseCase(timeoutContext)

	_orgDeliveryHttp.NewHandler(r, au, orgUseCase, b)
	//_userDeliveryHttp.NewHandler(r, userUseCase, orgUseCase, au)
	//_authDeliveryHttp.NewHandler(r, authUseCase, userUseCase, au, b)
	_authDeliveryHttp.NewHandler(r, au, authUseCase, userUseCase, b)
	//_invitationDeliveryHttp.NewHandler(r, invitationUseCase, userUseCase, orgUseCase, au)
	//_confirmationDeliveryHttp.NewHandler(r, confirmationUseCase)
}
