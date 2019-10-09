package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUserServiceConfigYaml(t *testing.T) {
	configReader, err := os.Open("./test/user-service.yaml")

	serviceConfig, err := NewUserServiceConfigYaml(configReader)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(serviceConfig.RealmRepos))
}
