package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/imtanmoy/authn/internal/authx"
	"github.com/imtanmoy/authn/internal/errorx"
	"github.com/imtanmoy/authn/models"
	"github.com/imtanmoy/authn/user"
	"github.com/imtanmoy/logx"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"log"
	"strings"
	"time"
)

type repository struct {
	conn *pgx.Conn
}

var _ user.Repository = (*repository)(nil)

// NewRepository will create an object that represent the user.Repository interface
func NewRepository(db *sql.DB) user.Repository {
	conn, err := stdlib.AcquireConn(db)
	if err != nil {
		log.Fatal(err)
	}
	return &repository{conn: conn}
}

func (repo *repository) FindAll(ctx context.Context) ([]*models.User, error) {
	rows, _ := repo.conn.Query(ctx, "SELECT id, name, email, created_at, updated_at "+
		"FROM users WHERE deleted_at IS NULL")
	var users []*models.User
	if rows.Err() != nil {
		return nil, errorx.ErrInternalDB
	}
	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return make([]*models.User, 0), err
		}
		users = append(users, &u)
	}
	return users, nil
}

//func (repo *repository) FindAllByOrganizationId(ctx context.Context, id int) ([]*models.User, error) {
//	db := repo.db.WithContext(ctx)
//	var users []*models.User
//	err := db.Model(&users).Where("organization_id = ?", id).Select()
//	err = godbx.ParsePgError(err)
//	return users, err
//}

func (repo *repository) Save(ctx context.Context, u *models.User) error {
	lastInsertedID := 0
	var createdAt time.Time
	var updatedAt time.Time
	err := repo.conn.QueryRow(ctx, "INSERT INTO users(name, email, password) "+
		"VALUES ($1,$2,$3) "+
		"RETURNING id, created_at, updated_at",
		u.Name, u.Email, u.Password).
		Scan(&lastInsertedID, &createdAt, &updatedAt)
	u.ID = lastInsertedID
	u.CreatedAt = createdAt
	u.UpdatedAt = updatedAt
	if err != nil {
		pgerr, ok := err.(*pgconn.PgError)
		fmt.Println(pgerr.SQLState())
		if ok {
			return errorx.ErrInternalDB
		}
		return errorx.ErrInternalServer
	}
	return err
}

//func (repo *repository) SaveUserOrganization(ctx context.Context, orgUser *models.UserOrganization) error {
//	db := repo.db.WithContext(ctx)
//	err := db.Insert(orgUser)
//	return err
//}

func (repo *repository) FindByID(ctx context.Context, id int) (*models.User, error) {
	var u models.User
	err := repo.conn.QueryRow(ctx, "SELECT id, name, email, created_at, updated_at "+
		"FROM users WHERE id = $1 "+
		"AND deleted_at IS NULL", id).
		Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, errorx.ErrorNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (repo *repository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	err := repo.conn.QueryRow(ctx, "SELECT id, name, email, password, created_at, updated_at "+
		"FROM users WHERE email = $1 "+
		"AND deleted_at IS NULL", email).
		Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, errorx.ErrorNotFound
		}
		return nil, err
	}
	return &u, err
}

func (repo *repository) GetByEmail(ctx context.Context, identity string) (authx.AuthUser, error) {
	return repo.FindByEmail(ctx, identity)
}

func (repo *repository) ExistsByID(ctx context.Context, id int) bool {
	found := 0
	err := repo.conn.QueryRow(ctx, "SELECT COUNT(*) AS found FROM users WHERE id = $1 AND deleted_at IS NULL", id).
		Scan(&found)
	if err != nil {
		logx.Fatal(err)
	}
	return found == 1
}

func (repo *repository) ExistsByEmail(ctx context.Context, email string) bool {
	found := 0
	err := repo.conn.QueryRow(ctx, "SELECT COUNT(*) AS found FROM users WHERE email = $1 AND deleted_at IS NULL", email).
		Scan(&found)
	if err != nil {
		logx.Fatal(err)
	}
	return found > 0
}

func (repo *repository) Delete(ctx context.Context, u *models.User) error {
	now := time.Now().UTC()
	_, err := repo.conn.Exec(ctx, "UPDATE users SET deleted_at = $1 WHERE id = $2", now, u.ID)
	u.DeletedAt = now
	return err
}

func (repo *repository) Update(ctx context.Context, u *models.User) error {
	now := time.Now().UTC()
	_, err := repo.conn.Exec(ctx, "UPDATE users SET name = $1, updated_at= $2 WHERE id = $3", u.Name, now, u.ID)
	return err
}
