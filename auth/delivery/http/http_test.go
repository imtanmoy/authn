package http

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/go-chi/chi"
	_authUseCase "github.com/imtanmoy/authn/auth/usecase"
	"github.com/imtanmoy/authn/internal/authx"
	"github.com/imtanmoy/authn/tests"
	_userRepo "github.com/imtanmoy/authn/user/repository"
	_userUseCase "github.com/imtanmoy/authn/user/usecase"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var (
	r  = chi.NewRouter()
	db *sql.DB
)

func init() {
	var err error
	db, err = tests.ConnectTestDB("localhost", 5432, "admin", "password", "authn")
	if err != nil {
		log.Fatal(err)
	}
}

func setup() {
	timeoutContext := 30 * time.Millisecond * time.Second
	userRepo := _userRepo.NewRepository(db)

	authxConfig := authx.AuthxConfig{
		SecretKey:             "test",
		AccessTokenExpireTime: 1,
	}

	au := authx.New(userRepo, &authxConfig)

	evt := tests.NewMockEventEmitter()
	userUseCase := _userUseCase.NewUseCase(userRepo, timeoutContext)
	authUseCase := _authUseCase.NewUseCase(userRepo, timeoutContext)
	NewHandler(r, au, authUseCase, userUseCase, evt)
}

func TestAuthHandler_Login(t *testing.T) {
	setup()
	defer tests.TruncateTestDB(db)

	ts := httptest.NewServer(r)
	defer ts.Close()

	tests.SeedUser(db)

	t.Run("Login with correct credentials", func(t *testing.T) {
		//POST login
		payload := &loginPayload{
			Email:    "test@test.com",
			Password: "password",
		}
		bodyRequest, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", ts.URL+"/login", bytes.NewReader(bodyRequest))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Login with wrong credentials", func(t *testing.T) {
		//POST login
		payload := &loginPayload{
			Email:    "wrong@test.com",
			Password: "wrong1234",
		}
		bodyRequest, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", ts.URL+"/login", bytes.NewReader(bodyRequest))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		body := w.Body.Bytes()
		assert.Contains(t, string(body), "invalid credentials")
	})
}

func TestAuthHandler_Register(t *testing.T) {
	setup()
	defer tests.TruncateTestDB(db)

	ts := httptest.NewServer(r)
	defer ts.Close()

	t.Run("Register with success", func(t *testing.T) {
		// POST register
		payload := &registerPayload{
			Name:            "Test",
			Email:           "test@test.com",
			Password:        "password",
			ConfirmPassword: "password",
		}
		bodyRequest, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", ts.URL+"/register", bytes.NewReader(bodyRequest))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var got *UserResponse
		body := w.Body.Bytes()
		err := json.Unmarshal(body, &got)
		assert.Nil(t, err)
		assert.Equal(t, payload.Email, got.Email)
	})

	t.Run("Register failed for invalid email", func(t *testing.T) {
		// POST register
		payload := &registerPayload{
			Name:            "Test",
			Email:           "test.com",
			Password:        "password",
			ConfirmPassword: "password",
		}
		bodyRequest, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", ts.URL+"/register", bytes.NewReader(bodyRequest))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		body := w.Body.Bytes()
		assert.Contains(t, string(body), "The email field must be a valid email address")
	})

	t.Run("Register failed for password and confirm password doest not match", func(t *testing.T) {
		// POST register
		payload := &registerPayload{
			Name:            "Test",
			Email:           "test.com",
			Password:        "password1",
			ConfirmPassword: "password2",
		}
		bodyRequest, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", ts.URL+"/register", bytes.NewReader(bodyRequest))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		body := w.Body.Bytes()
		got := strings.Contains(string(body), "password and confirmation password do not match")
		assert.True(t, got)
	})

	t.Run("Register failed for email already exist", func(t *testing.T) {
		// POST register
		payload := &registerPayload{
			Name:            "Test",
			Email:           "test@test.com",
			Password:        "password",
			ConfirmPassword: "password",
		}
		bodyRequest, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", ts.URL+"/register", bytes.NewReader(bodyRequest))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		body := w.Body.Bytes()
		got := strings.Contains(string(body), "user with this email already exists")
		assert.True(t, got)
	})
}
