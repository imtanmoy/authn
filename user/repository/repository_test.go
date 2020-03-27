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
	defer tests.TruncateTestDB(db)
	ctx := context.Background()

	users := tests.FakeUsers(10)

	err := tests.InsertTestUsers(db, users)
	assert.Nil(t, err)
	usr, err := repo.FindAll(ctx)
	assert.Nil(t, err)
	assert.Equal(t, len(usr), 10)
}

func TestRepository_Save(t *testing.T) {
	defer tests.TruncateTestDB(db)
	ctx := context.Background()

	t.Parallel()

	users := tests.FakeUsers(10)

	for _, u := range users {
		u := u
		t.Run(u.Name, func(t *testing.T) {
			err := repo.Save(ctx, u)
			assert.Nil(t, err)
			assert.NotZero(t, u.CreatedAt)
			assert.NotZero(t, u.UpdatedAt)
		})
	}

}
