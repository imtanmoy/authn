package db

import (
	"context"
	"strconv"

	"github.com/go-pg/pg/v9"

	"github.com/imtanmoy/authn/config"
)

var DB *pg.DB

type dbLogger struct{}

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	//fmt.Println(q.FormattedQuery())
	return nil
}

func InitDB() error {
	db, err := connectDB(
		config.Conf.DB.USERNAME,
		config.Conf.DB.PASSWORD,
		config.Conf.DB.DBNAME,
		config.Conf.DB.HOST+":"+strconv.Itoa(config.Conf.DB.PORT),
	)
	if err != nil {
		return err
	}
	db.AddQueryHook(dbLogger{})
	DB = db
	return nil
}

func connectDB(username, password, database, address string) (*pg.DB, error) {
	db := pg.Connect(&pg.Options{
		User:     username,
		Password: password,
		Database: database,
		Addr:     address,
	})
	var n int
	_, err := db.QueryOne(pg.Scan(&n), "SELECT 1")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Shutdown() error {
	return closeDB(DB)
}

func closeDB(db *pg.DB) error {
	return db.Close()
}
