package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDatabaseConn(t *testing.T) {
	conn, err := connectDB("admin", "password", "authn", "0.0.0.0:5432")
	assert.Nil(t, err)
	assert.Equal(t, "DB<Addr=\"0.0.0.0:5432\">", conn.String())
}

func TestDatabaseShutdown(t *testing.T) {
	conn, err := connectDB("admin", "password", "authn", "0.0.0.0:5432")
	assert.Nil(t, err)
	err = closeDB(conn)
	assert.Nil(t, err)
}
