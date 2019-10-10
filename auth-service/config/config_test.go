package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfigFile(t *testing.T) {
	configReader, err := os.Open("./auth-config-test.yaml")
	assert.NoError(t, err)

	Init(configReader)

	ac := GetConfig()
	assert.Equal(t, len(ac.Realms), 2)

}
