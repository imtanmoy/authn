package repository

import (
	"context"
	"database/sql"
	"github.com/imtanmoy/authn/tests"
	"github.com/imtanmoy/authn/user"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

var db *sql.DB
var repo user.Repository

func init() {
	var err error
	db, err = tests.ConnectTestDB("localhost", 5432, "admin", "password", "authn")
	if err != nil {
		log.Fatal(err)
	}
	repo = NewRepository(db)
}

func TestRepository_FindAll(t *testing.T) {
	tests.TruncateTestDB(db)
	defer tests.TruncateTestDB(db)
	ctx := context.Background()

	users := tests.FakeUsers(10)

	err := tests.InsertTestUsers(db, users)
	usr, err := repo.FindAll(ctx)
	assert.Nil(t, err)
	assert.Equal(t, len(usr), 10)
}

func TestRepository_Save(t *testing.T) {
	tests.TruncateTestDB(db)
	defer tests.TruncateTestDB(db)
	ctx := context.Background()

	t.Parallel()

	users := tests.FakeUsers(10)

	for _, u := range users {
		u := u
		t.Run(u.Name+" -> save", func(t *testing.T) {
			err := repo.Save(ctx, u)
			assert.Nil(t, err)
			assert.NotZero(t, u.CreatedAt)
			assert.NotZero(t, u.UpdatedAt)
		})
	}
}

func TestRepository_ExistsByEmail(t *testing.T) {
	tests.TruncateTestDB(db)
	defer tests.TruncateTestDB(db)
	ctx := context.Background()

	users := tests.FakeUsers(10)

	_ = tests.InsertTestUsers(db, users)

	data := []struct {
		email  string
		result bool
	}{
		{email: "test@notfound.com", result: false},
	}
	for _, u := range users {
		data = append(data, struct {
			email  string
			result bool
		}{email: u.Email, result: true})
	}
	for _, d := range data {
		assert.Equal(t, d.result, repo.ExistsByEmail(ctx, d.email))
	}
}

func TestRepository_ExistsByID(t *testing.T) {
	tests.TruncateTestDB(db)
	defer tests.TruncateTestDB(db)
	ctx := context.Background()

	users := tests.FakeUsers(10)

	_ = tests.InsertTestUsers(db, users)

	data := []struct {
		id     int
		result bool
	}{
		{id: 12, result: false},
		{id: 11, result: false},
	}
	for i, _ := range users {
		data = append(data, struct {
			id     int
			result bool
		}{id: i + 1, result: true})
	}
	for _, d := range data {
		assert.Equal(t, d.result, repo.ExistsByID(ctx, d.id))
	}
}

func TestRepository_FindByID(t *testing.T) {
	tests.TruncateTestDB(db)
	defer tests.TruncateTestDB(db)
	ctx := context.Background()

	users := tests.FakeUsers(10)

	_ = tests.InsertTestUsers(db, users)

	data := []struct {
		id     int
		result bool
	}{
		{id: 12, result: false},
		{id: 11, result: false},
	}
	for i, _ := range users {
		data = append(data, struct {
			id     int
			result bool
		}{id: i + 1, result: true})
	}
	for _, d := range data {
		got, err := repo.FindByID(ctx, d.id)
		if d.result == false {
			assert.Error(t, err)
			assert.Nil(t, got)
			assert.Equal(t, err.Error(), "no rows in result set")
		} else {
			assert.Nil(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, got.ID, d.id)
		}
	}
}

func TestRepository_FindByEmail(t *testing.T) {
	tests.TruncateTestDB(db)
	defer tests.TruncateTestDB(db)
	ctx := context.Background()

	users := tests.FakeUsers(10)

	_ = tests.InsertTestUsers(db, users)

	data := []struct {
		email  string
		result bool
	}{
		{email: "test1@notfound.com", result: false},
		{email: "test2@notfound.com", result: false},
	}
	for _, u := range users {
		data = append(data, struct {
			email  string
			result bool
		}{email: u.Email, result: true})
	}
	for _, d := range data {
		got, err := repo.FindByEmail(ctx, d.email)
		if d.result == false {
			assert.Error(t, err)
			assert.Equal(t, err.Error(), "no rows in result set")
			assert.Nil(t, got)
		} else {
			assert.Nil(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, got.Email, d.email)
		}
	}
}

func TestRepository_GetByEmail(t *testing.T) {
	tests.TruncateTestDB(db)
	defer tests.TruncateTestDB(db)
	ctx := context.Background()

	users := tests.FakeUsers(10)

	_ = tests.InsertTestUsers(db, users)

	data := []struct {
		email  string
		result bool
	}{
		{email: "test1@notfound.com", result: false},
		{email: "test2@notfound.com", result: false},
	}
	for _, u := range users {
		data = append(data, struct {
			email  string
			result bool
		}{email: u.Email, result: true})
	}
	for _, d := range data {
		got, err := repo.GetByEmail(ctx, d.email)
		if d.result == false {
			assert.Error(t, err)
			assert.Equal(t, err.Error(), "no rows in result set")
			assert.Nil(t, got)
		} else {
			assert.Nil(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, got.GetEmail(), d.email)
		}
	}
}

func TestRepository_Delete(t *testing.T) {
	tests.TruncateTestDB(db)
	defer tests.TruncateTestDB(db)
	ctx := context.Background()

	tests.SeedUsers(db)

	users, _ := repo.FindAll(ctx)

	for _, u := range users {
		err := repo.Delete(ctx, u)
		assert.Nil(t, err)
		assert.NotZero(t, u.DeletedAt)
	}
}

func TestRepository_Update(t *testing.T) {
	tests.TruncateTestDB(db)
	defer tests.TruncateTestDB(db)
	ctx := context.Background()

	tests.SeedUsers(db)

	users, _ := repo.FindAll(ctx)

	for _, u := range users {
		u.Name = reverse(u.Name)
	}

	for _, u := range users {
		err := repo.Update(ctx, u)
		assert.Nil(t, err)
	}
}

func reverse(s string) string {
	result := ""
	for _, v := range s {
		result = string(v) + result
	}
	return result
}
