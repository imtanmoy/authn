package http

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-pg/pg/v9"
	"github.com/imtanmoy/authn/internal/authx"
	"github.com/imtanmoy/authn/models"
	_userRepo "github.com/imtanmoy/authn/user/repository"
	_userUseCase "github.com/imtanmoy/authn/user/usecase"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	r  = chi.NewRouter()
	DB *pg.DB
)

func init() {
	DB, _ = connectDB("admin", "password", "authn", "0.0.0.0:5432")
	setup()
}

func connectDB(username, password, database, address string) (*pg.DB, error) {
	connect := pg.Connect(&pg.Options{
		User:     username,
		Password: password,
		Database: database,
		Addr:     address,
	})
	return connect, nil
}

func setup() {
	timeoutContext := 30 * time.Millisecond * time.Second //TODO it will come from config
	userRepo := _userRepo.NewRepository(DB)
	authxConfig := authx.AuthxConfig{
		SecretKey:             "test",
		AccessTokenExpireTime: 1,
	}

	au := authx.New(userRepo, &authxConfig)

	userUseCase := _userUseCase.NewUseCase(userRepo, timeoutContext)
	NewHandler(r, userUseCase, au)
}

func teardown() {
	_, _ = DB.Exec("TRUNCATE TABLE users, organizations RESTART IDENTITY;")
}

func TestUserHandler_List(t *testing.T) {
	defer teardown()
	ts := httptest.NewServer(r)
	defer ts.Close()
	req, _ := http.NewRequest("GET", ts.URL+"/users", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var expected []*models.UserResponse
	body := w.Body.Bytes()
	err := json.Unmarshal(body, &expected)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(expected))
}

func insertOrganization() {
	_, _ = DB.Exec("INSERT INTO organizations(name, created_at, updated_at) VALUES ('Example', now(), now());")
}

func insertUser() {
	_, _ = DB.Exec("INSERT INTO users(name, designation, email, password, enabled, organization_id, created_by, updated_by, deleted_by,	joined_at, created_at, updated_at, deleted_at)	VALUES ('Test', 'dev', 'test@test.com', 'password', bool(1), 1, NULL, NULL, NULL, now(), now(), now(), null);")
}

func TestUserHandler_Create(t *testing.T) {
	defer teardown()

	insertOrganization()

	ts := httptest.NewServer(r)
	defer ts.Close()

	// POST user
	payload := &userPayload{
		Name:            "Test",
		Email:           "test@test.com",
		Password:        "password",
		ConfirmPassword: "password",
		Designation:     "DEv",
		OrganizationId:  1,
	}
	bodyRequest, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", ts.URL+"/users", bytes.NewReader(bodyRequest))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var expected *models.UserResponse
	body := w.Body.Bytes()
	err := json.Unmarshal(body, &expected)
	assert.Nil(t, err)
	assert.Equal(t, "test@test.com", expected.Email)

	// Bad Request for same email
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

}

func TestUserHandler_Get(t *testing.T) {
	defer teardown()
	insertOrganization()
	insertUser()
	ts := httptest.NewServer(r)
	defer ts.Close()
	req, _ := http.NewRequest("GET", ts.URL+"/users/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.Bytes()
	var expected *models.UserResponse
	err := json.Unmarshal(body, &expected)
	assert.Nil(t, err)
	assert.Equal(t, "test@test.com", expected.Email)

	//Not Found
	req, _ = http.NewRequest("GET", ts.URL+"/users/100", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUserHandler_Update(t *testing.T) {
	defer teardown()
	insertOrganization()
	insertUser()
	ts := httptest.NewServer(r)
	defer ts.Close()
	payload := &userUpdatePayload{
		Name:            "Test Update",
		Password:        "password1234",
		ConfirmPassword: "password1234",
		Designation:     "tester",
	}
	bodyRequest, _ := json.Marshal(payload)
	req, _ := http.NewRequest("PUT", ts.URL+"/users/1", bytes.NewReader(bodyRequest))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	body := w.Body.Bytes()
	var expected *models.UserResponse
	err := json.Unmarshal(body, &expected)
	assert.Nil(t, err)
	assert.Equal(t, "test@test.com", expected.Email)
	assert.Equal(t, "Test Update", expected.Name)
	assert.Equal(t, "tester", expected.Designation)

	//Not Found
	req, _ = http.NewRequest("PUT", ts.URL+"/users/100", bytes.NewReader(bodyRequest))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUserHandler_Delete(t *testing.T) {
	defer teardown()
	insertOrganization()
	insertUser()
	ts := httptest.NewServer(r)
	defer ts.Close()
	req, _ := http.NewRequest("DELETE", ts.URL+"/users/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	//Not Found
	req, _ = http.NewRequest("DELETE", ts.URL+"/users/100", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
