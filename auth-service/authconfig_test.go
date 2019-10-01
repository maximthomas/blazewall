package main

import (
	"os"
	"testing"

	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
)

func TestReadConfigFile(t *testing.T) {
	configReader, err := os.Open("./test/auth-config.yaml")

	ac, err := NewAuthServiceConfigYaml(configReader)
	assert.NoError(t, err)
	assert.Equal(t, len(ac.Realms), 2)

	clientIDInt, ok := ac.Realms[0].AuthConfig[0].Parameters["clientID"]
	assert.True(t, ok)
	clientID := cast.ToString(clientIDInt)
	assert.NotEmpty(t, clientID)

}
