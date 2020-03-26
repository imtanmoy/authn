package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/imtanmoy/authn/internal/errorx"
	"github.com/imtanmoy/authn/models"
	"github.com/imtanmoy/authn/user"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"log"
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
	rows, _ := repo.conn.Query(ctx, "select * from users")
	var users []*models.User
	if rows.Err() != nil {
		return nil, errorx.ErrInternalDB
	}
	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)
		if err != nil {
			return make([]*models.User, 0), err
		}
		users = append(users, &u)
	}

	//if err != nil {
	//	_, ok := err.(pg.Error)
	//	if ok {
	//		return nil, errorx.ErrInternalDB
	//	} else {
	//		return nil, errorx.ErrInternalServer
	//	}
	//}
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
	//res, err := repo.conn.Exec(ctx,
	//	"INSERT INTO users(name, email, password) "+
	//		"VALUES ($1,$2,$3) "+
	//		"RETURNING id,name,email,password,created_at, updated_at, deleted_at",
	//	u.Name, u.Email, u.Password)

	lastInsertedID := 0
	var createdAt time.Time
	err := repo.conn.QueryRow(ctx, "INSERT INTO users(name, email, password) "+
		"VALUES ($1,$2,$3) "+
		"RETURNING id, created_at",
		u.Name, u.Email, u.Password).
		Scan(&lastInsertedID, &createdAt)
	u.ID =lastInsertedID
	fmt.Println(createdAt)
	if err != nil {
		return err
	}
	return err
}

//func (repo *repository) SaveUserOrganization(ctx context.Context, orgUser *models.UserOrganization) error {
//	db := repo.db.WithContext(ctx)
//	err := db.Insert(orgUser)
//	return err
//}

//func (repo *repository) Find(ctx context.Context, id int) (*models.User, error) {
//	db := repo.db.WithContext(ctx)
//	var u models.User
//	err := db.Model(&u).Where("id = ?", id).Select()
//	if err != nil {
//		if errors.Is(err, pg.ErrNoRows) {
//			return nil, errorx.ErrorNotFound
//		} else {
//			panic(err)
//		}
//	}
//	return &u, err
//}
//
//func (repo *repository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
//	db := repo.db.WithContext(ctx)
//	var u models.User
//	//_, err := db.Query(&u, `
//	//							SELECT "user"."id",
//	//								   "user"."name",
//	//								   "user"."email",
//	//								   "user"."created_at",
//	//								   "user"."updated_at",
//	//								   "user"."deleted_at"
//	//							FROM "users" AS "user"
//	//									 LEFT OUTER JOIN users_organizations AS user_organization ON "user".id = user_organization.user_id
//	//									 LEFT OUTER JOIN organizations organization on user_organization.organization_id = organization.id
//	//							WHERE "user".email = "jjsj"
//	//							  AND "user"."deleted_at" IS NULL`,
//	//	email)
//	err := db.Model(&u).ExcludeColumn("password").Where("email = ?", email).Relation("Organizations").Select()
//	if err != nil {
//		if errors.Is(err, pg.ErrNoRows) {
//			return nil, errorx.ErrorNotFound
//		} else {
//			panic(err)
//		}
//	}
//	return &u, err
//}
//
//func (repo *repository) GetByEmail(ctx context.Context, identity string) (authx.AuthUser, error) {
//	return repo.FindByEmail(ctx, identity)
//}
//
//func (repo *repository) Exists(ctx context.Context, id int) bool {
//	db := repo.db.WithContext(ctx)
//	u := new(models.User)
//	err := db.Model(u).Where("id = ?", id).Select()
//	if err != nil {
//		if errors.Is(err, pg.ErrNoRows) {
//			return false
//		} else {
//			panic(err)
//		}
//	}
//	return u.ID == id
//}
//
//func (repo *repository) ExistsByEmail(ctx context.Context, email string) bool {
//	db := repo.db.WithContext(ctx)
//	u := new(models.User)
//	err := db.Model(u).Where("email = ?", email).Select()
//	if err != nil {
//		if errors.Is(err, pg.ErrNoRows) {
//			return false
//		} else {
//			panic(err)
//		}
//	}
//	return u.Email == email
//}

//func (repo *repository) Delete(ctx context.Context, u *models.User) error {
//	db := repo.db.WithContext(ctx)
//	err := db.Delete(u)
//	err = godbx.ParsePgError(err)
//	return err
//}
//
//func (repo *repository) Update(ctx context.Context, u *models.User) error {
//	db := repo.db.WithContext(ctx)
//	err := db.Update(u)
//	err = godbx.ParsePgError(err)
//	return err
//}
