package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitConfig(t *testing.T) {
	InitConfig()
	assert.NotNil(t, Conf)
}
