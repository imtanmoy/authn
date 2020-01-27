package db

import (
	"fmt"
	"testing"
)

func TestDatabaseConn(t *testing.T) {
	_, err := connectDB("admin", "password", "authn", "0.0.0.0:5432")
	fmt.Println(err)
	if err != nil {
		t.Error(fmt.Sprintf("Expected the pointer of db"))
	}
}
