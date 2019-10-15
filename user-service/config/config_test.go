package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUserServiceConfigYaml(t *testing.T) {
	configReader, err := os.Open("./user-config-test.yaml")
	assert.NoError(t, err)
	Init(configReader)
	usc := GetUserServiceConfig()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(usc.RealmRepos))
}
