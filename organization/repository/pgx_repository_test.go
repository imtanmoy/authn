package repository

import (
	"context"
	"database/sql"
	"github.com/imtanmoy/authn/internal/errorx"
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

func TestPgxRepository_FindByID(t *testing.T) {
	tests.TruncateTestDB(db)
	defer tests.TruncateTestDB(db)
	ctx := context.Background()

	//t.Parallel()

	tests.SeedUser(db)

	orgs := tests.FakeOrgs(10)

	err := tests.InsertTestOrgs(db, orgs)
	if err != nil {
		panic(err)
	}

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
		got, err := repo.FindByID(ctx, d.id)
		if d.result == false {
			assert.Error(t, err)
			assert.Nil(t, got)
			assert.Equal(t, err, errorx.ErrorNotFound)
		} else {
			assert.Nil(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, got.ID, d.id)
		}
	}
}
