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
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

var (
	r    = chi.NewRouter()
	db   *sql.DB
	conn *pgx.Conn
	aux  *authx.Authx
)

func init() {
	var err error
	db, err = tests.ConnectTestDB("localhost", 5432, "admin", "password", "authn")
	if err != nil {
		log.Fatal(err)
	}
	conn, err = stdlib.AcquireConn(db)
	if err != nil {
		log.Fatal(err)
	}
	setup()
}

func setup() {
	timeoutContext := 30 * time.Millisecond * time.Second
	userRepo := _userRepo.NewPgxRepository(conn)
	orgRepo := _orgRepo.NewPgxRepository(conn)

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
	require.NoError(t, err)

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

func TestOrgHandler_Get(t *testing.T) {
	tests.TruncateTestDB(db)
	defer tests.TruncateTestDB(db)

	t.Parallel()

	ts := httptest.NewServer(r)
	defer ts.Close()

	tests.SeedUser(db)

	token, err := aux.GenerateToken("test@test.com")
	require.NoError(t, err)

	orgs := tests.FakeOrgs(10)

	err = tests.InsertTestOrgs(db, orgs)
	require.NoError(t, err)

	data := []struct {
		id     int
		result bool
	}{
		{id: 12, result: false},
		{id: 11, result: false},
	}
	for i, _ := range orgs {
		data = append(data, struct {
			id     int
			result bool
		}{id: i + 1, result: true})
	}

	for _, d := range data {
		t.Run("Org -> "+strconv.Itoa(d.id)+"->GET", func(t *testing.T) {
			req, _ := http.NewRequest("GET", ts.URL+"/organizations/"+strconv.Itoa(d.id), nil)

			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			body := w.Body.Bytes()
			if d.result {
				var got *orgResponse
				assert.Equal(t, http.StatusOK, w.Code)
				err = json.Unmarshal(body, &got)
				assert.Nil(t, err)
				assert.Equal(t, d.id, got.ID)
			} else {
				assert.Equal(t, http.StatusNotFound, w.Code)
				assert.Contains(t, string(body), "organization not found")
			}
		})
	}
}
