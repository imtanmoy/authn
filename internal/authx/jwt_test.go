package authx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthx_createToken(t *testing.T) {
	_, err := createToken("test@test.com", "test", 5)
	assert.Nil(t, err)
}

func TestAuthx_parseToken(t *testing.T) {
	token, _ := createToken("test@test.com", "test", 5)
	parsedToken, err := parseToken(token, "test")
	assert.Nil(t, err)
	assert.Equal(t, true, parsedToken.Valid)
	claims, ok := parsedToken.Claims.(*Claims)
	assert.Equal(t, true, ok)
	assert.Equal(t, "test@test.com", claims.Identity)
}
