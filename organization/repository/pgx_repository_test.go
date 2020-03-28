package repository

import (
	"context"
	"database/sql"
	"github.com/imtanmoy/authn/organization"
	"github.com/imtanmoy/authn/tests"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

var db *sql.DB
var repo organization.Repository

func init() {
	var err error
	db, err = tests.ConnectTestDB("localhost", 5432, "admin", "password", "authn")
	if err != nil {
		log.Fatal(err)
	}
	repo = NewPgxRepository(db)
}

func TestRepository_Save(t *testing.T) {
	tests.TruncateTestDB(db)
	defer tests.TruncateTestDB(db)
	ctx := context.Background()

	t.Parallel()

	tests.SeedUser(db)

	orgs := tests.FakeOrgs(10)

	for _, o := range orgs {
		o := o
		t.Run(o.Name+" -> save", func(t *testing.T) {
			err := repo.Save(ctx, o)
			assert.Nil(t, err)
			assert.NotZero(t, o.CreatedAt)
			assert.NotZero(t, o.UpdatedAt)
		})
	}
}
