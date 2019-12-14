package server

import (
	"github.com/imtanmoy/authy/db"
	_orgDeliveryHttp "github.com/imtanmoy/authy/organization/delivery/http"
	_orgRepo "github.com/imtanmoy/authy/organization/repository"
	_orgUseCase "github.com/imtanmoy/authy/organization/usecase"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	_chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// New configures application resources and routes.
func New() (*chi.Mux, error) {

	r := chi.NewRouter()
	r.Use(_chiMiddleware.Recoverer)
	r.Use(_chiMiddleware.RequestID)
	r.Use(_chiMiddleware.RealIP)
	r.Use(_chiMiddleware.DefaultCompress)
	r.Use(_chiMiddleware.Timeout(15 * time.Second))

	//r.Use(logging.NewStructuredLogger(logger))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	//r.Use(corsConfig().Handler)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	})

	timeoutContext := 30 * time.Millisecond * time.Second //TODO it will come from config

	orgRepo := _orgRepo.NewRepository(db.DB)

	orgUseCase := _orgUseCase.NewUsecase(orgRepo, timeoutContext)

	_orgDeliveryHttp.NewHandler(r, orgUseCase)

	return r, nil
}
