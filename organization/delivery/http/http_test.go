package http

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/imtanmoy/authn/internal/authx"
	_orgRepo "github.com/imtanmoy/authn/organization/repository"
	_orgUseCase "github.com/imtanmoy/authn/organization/usecase"
	"github.com/imtanmoy/authn/tests"
	_userRepo "github.com/imtanmoy/authn/user/repository"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	r   = chi.NewRouter()
	db  *sql.DB
	aux *authx.Authx
)

func init() {
	var err error
	db, err = tests.ConnectTestDB("localhost", 5432, "admin", "password", "authn")
	if err != nil {
		log.Fatal(err)
	}
	setup()
}

func setup() {
	timeoutContext := 30 * time.Millisecond * time.Second
	userRepo := _userRepo.NewRepository(db)
	orgRepo := _orgRepo.NewRepository(db)

	authxConfig := authx.AuthxConfig{
		SecretKey:             "test",
		AccessTokenExpireTime: 1,
	}

	aux = authx.New(userRepo, &authxConfig)

	evt := tests.NewMockEventEmitter()
	orgUseCase := _orgUseCase.NewUseCase(orgRepo, timeoutContext)
	NewHandler(r, aux, orgUseCase, evt)
}

func TestOrgHandler_Create(t *testing.T) {
	tests.TruncateTestDB(db)
	defer tests.TruncateTestDB(db)

	ts := httptest.NewServer(r)
	defer ts.Close()

	tests.SeedUser(db)

	token, err := aux.GenerateToken("test@test.com")
	if err != nil {
		panic(err)
	}

	t.Run("Organization Create success", func(t *testing.T) {
		// POST register
		payload := &orgCreatePayload{
			Name: "Test Orgs",
		}
		bodyRequest, _ := json.Marshal(payload)

		req, _ := http.NewRequest("POST", ts.URL+"/organizations", bytes.NewReader(bodyRequest))

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
		var got *orgResponse
		body := w.Body.Bytes()
		err = json.Unmarshal(body, &got)
		assert.Nil(t, err)
		assert.Equal(t, payload.Name, got.Name)
		assert.Equal(t, 1, got.OwnerId)
	})

	t.Run("Organization Create Fail", func(t *testing.T) {
		// POST register
		payload := &orgCreatePayload{
			Name: "",
		}
		bodyRequest, _ := json.Marshal(payload)

		req, _ := http.NewRequest("POST", ts.URL+"/organizations", bytes.NewReader(bodyRequest))

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		body := w.Body.Bytes()
		assert.Contains(t, string(body), "invalid request")
	})
}
