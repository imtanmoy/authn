package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitConfig(t *testing.T) {
	conf, err := initViper("../")
	assert.NotNil(t, conf)
	assert.Nil(t, err)
}
