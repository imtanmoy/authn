package http

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/go-chi/chi"
	_authUseCase "github.com/imtanmoy/authn/auth/usecase"
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
	timeoutContext := 30 * time.Millisecond * time.Second //TODO it will come from config
	userRepo := _userRepo.NewRepository(db)
	evt := tests.NewMockEventEmitter()
	userUseCase := _userUseCase.NewUseCase(userRepo, timeoutContext)
	authUseCase := _authUseCase.NewUseCase(userRepo, timeoutContext)
	NewHandler(r, authUseCase, userUseCase, evt)
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
		got := strings.Contains(string(body), "The email field must be a valid email address")
		assert.True(t, got)
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
