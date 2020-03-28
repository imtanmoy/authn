package repository

import (
	"context"
	"github.com/imtanmoy/authn/internal/errorx"
	"github.com/imtanmoy/authn/models"
	"github.com/imtanmoy/authn/organization"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"strings"
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

func (repo *pgxRepository) FindByID(ctx context.Context, id int) (*models.Organization, error) {
	var org models.Organization
	err := repo.conn.QueryRow(ctx, "SELECT id, name, owner_id, created_at, updated_at "+
		"FROM organizations WHERE id = $1 "+
		"AND deleted_at IS NULL", id).
		Scan(&org.ID, &org.Name, &org.OwnerID, &org.CreatedAt, &org.UpdatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, errorx.ErrorNotFound
		}
		return nil, err
	}
	return &org, nil
}

var _ organization.Repository = (*pgxRepository)(nil)

// NewRepository will create an object that represent the organization.Repository interface
func NewPgxRepository(conn *pgx.Conn) organization.Repository {
	return &pgxRepository{conn: conn}
}
