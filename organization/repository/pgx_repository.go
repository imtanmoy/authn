package repository

import (
	"context"
	"database/sql"
	"github.com/imtanmoy/authn/internal/errorx"
	"github.com/imtanmoy/authn/models"
	"github.com/imtanmoy/authn/organization"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"log"
	"time"
)

type pgxRepository struct {
	conn *pgx.Conn
}

func (repo *pgxRepository) Save(ctx context.Context, org *models.Organization) error {
	lastInsertedID := 0
	var createdAt time.Time
	var updatedAt time.Time
	err := repo.conn.QueryRow(ctx, "INSERT INTO organizations(name, owner_id) "+
		"VALUES ($1,$2) "+
		"RETURNING id, created_at, updated_at",
		org.Name, org.OwnerID).
		Scan(&lastInsertedID, &createdAt, &updatedAt)
	org.ID = lastInsertedID
	org.CreatedAt = createdAt
	org.UpdatedAt = updatedAt
	if err != nil {
		_, ok := err.(*pgconn.PgError)
		if ok {
			return errorx.ErrInternalDB
		}
		return errorx.ErrInternalServer
	}
	return nil
}

var _ organization.Repository = (*pgxRepository)(nil)

// NewRepository will create an object that represent the organization.Repository interface
func NewPgxRepository(db *sql.DB) organization.Repository {
	conn, err := stdlib.AcquireConn(db)
	if err != nil {
		log.Fatal(err)
	}
	return &pgxRepository{conn: conn}
}
