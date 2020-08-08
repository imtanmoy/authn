package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/imtanmoy/authn/events"
	"github.com/imtanmoy/authn/internal/authx"
	"github.com/imtanmoy/authn/internal/errorx"
	"github.com/imtanmoy/authn/models"
	"github.com/imtanmoy/authn/organization"
	"github.com/imtanmoy/httpx"
	param "github.com/oceanicdev/chi-param"
	"gopkg.in/thedevsaddam/govalidator.v1"
	"net/http"
	"net/url"
	"time"
)

type contextKey string

const (
	orgKey contextKey = "organization"
)

type orgCreatePayload struct {
	Name string `json:"name"`
}

func (rp *orgCreatePayload) validate() url.Values {
	rules := govalidator.MapData{
		"name": []string{"required", "min:4", "max:100"},
	}
	opts := govalidator.Options{
		Data:  rp,
		Rules: rules,
	}

	v := govalidator.New(opts)
	e := v.ValidateStruct()
	return e
}

type orgResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	OwnerId   int       `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// orgHandler  represent the http handler for org
type orgHandler struct {
	useCase organization.UseCase
	*authx.Authx
	event events.EventEmitter
}

func (handler *orgHandler) OrgCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if ctx == nil {
			ctx = context.Background()
		}
		id, err := param.Int(r, "id")
		if err != nil {
			httpx.ResponseJSONError(w, r, http.StatusBadRequest, "invalid request parameter", err)
			return
		}
		org, err := handler.useCase.FindByID(ctx, id)
		if err != nil {
			if errors.Is(err, errorx.ErrorNotFound) {
				httpx.ResponseJSONError(w, r, http.StatusNotFound, "organization not found", err)
			} else {
				panic(err)
			}
			return
		}
		ctx = context.WithValue(r.Context(), orgKey, org)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (handler *orgHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	data := &orgCreatePayload{}
	if err := httpx.DecodeJSON(r, data); err != nil {
		var mr *httpx.MalformedRequest
		if errors.As(err, &mr) {
			httpx.ResponseJSONError(w, r, mr.Status, mr.Status, mr.Msg)
			return
		}
		panic(err)
	}

	validationErrors := data.validate()

	if len(validationErrors) > 0 {
		httpx.ResponseJSONError(w, r, 400, "invalid request", validationErrors)
		return
	}

	u, err := handler.GetCurrentUser(r)
	us, ok := u.(*models.User)
	if err != nil || !ok {
		panic(fmt.Sprintf("could not upgrade user to an authable user, type: %T", u))
	}

	var org models.Organization
	org.Name = data.Name
	org.OwnerID = us.ID
	err = handler.useCase.Save(ctx, &org)
	if err != nil {
		httpx.ResponseJSONError(w, r, http.StatusInternalServerError, err)
		return
	}

	httpx.ResponseJSON(w, http.StatusCreated, &orgResponse{
		ID:        org.ID,
		Name:      org.Name,
		OwnerId:   org.OwnerID,
		CreatedAt: org.CreatedAt,
		UpdatedAt: org.UpdatedAt,
	})
	return
}

func (handler *orgHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	org, ok := ctx.Value(orgKey).(*models.Organization)
	if !ok {
		httpx.ResponseJSONError(w, r, http.StatusInternalServerError, httpx.ErrInternalServerError)
		return
	}
	httpx.ResponseJSON(w, http.StatusOK, &orgResponse{
		ID:        org.ID,
		Name:      org.Name,
		OwnerId:   org.OwnerID,
		CreatedAt: org.CreatedAt,
		UpdatedAt: org.UpdatedAt,
	})
	return
}

// NewHandler will initialize the org's resources endpoint
func NewHandler(
	r *chi.Mux,
	aux *authx.Authx,
	useCase organization.UseCase,
	event events.EventEmitter,
) {
	handler := &orgHandler{
		useCase: useCase,
		Authx:   aux,
		event:   event,
	}
	r.Route("/organizations", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(handler.AuthMiddleware)
			r.Post("/", handler.Create)
			r.Group(func(r chi.Router) {
				r.Use(handler.OrgCtx)
				r.Get("/{id}", handler.Get)
				//				r.Put("/{id}", handler.Update)
				//				r.Delete("/{id}", handler.Delete)
			})
		})
	})
}
