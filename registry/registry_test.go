package registry

import (
	"fmt"
	"github.com/imtanmoy/authn/config"
	"github.com/imtanmoy/logx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testConfig() *config.Config {
	con := &config.Config{
		ENVIRONMENT:           "test",
		DEBUG:                 false,
		JwtSecretKey:          "test",
		JwtAccessTokenExpires: 1,
		SERVER: config.Server{
			HOST: "0.0.0.0",
			PORT: 8080,
		},
		DB: config.DB{
			HOST:     "localhost",
			PORT:     5432,
			USERNAME: "admin",
			PASSWORD: "password",
			DBNAME:   "authn",
		},
	}
	return con
}

func TestNewRegistry(t *testing.T) {
	conf := testConfig()
	r := NewRegistry(*conf)
	err := r.Init()
	if err != nil {
		logx.Fatalf("%s : %s", "could not init registry", err)
	}
	assert.NotNil(t, r)
	r.Close()
}

func TestConnectPgx(t *testing.T) {
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", "admin", "password", "localhost", 5432, "authn")
	conn, err := ConnectPgx(connString)
	assert.Nil(t, err)
	assert.NotNil(t, conn)
}

func TestConnectDB(t *testing.T) {
	db, err := connectDB("localhost", 5432, "admin", "password", "authn")
	assert.Nil(t, err)
	assert.NotNil(t, db)
}
