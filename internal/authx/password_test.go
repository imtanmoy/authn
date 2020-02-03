package authx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthx_HashPassword(t *testing.T) {
	_, err := hashPassword("password")
	assert.Nil(t, err)
}
func TestAuthx_ComparePassword(t *testing.T) {
	hash, _ := hashPassword("password")
	b := comparePasswords("password", hash)
	assert.Equal(t, b, false)
	b2 := comparePasswords("password2", hash)
	assert.NotEqual(t, b2, false)
}
