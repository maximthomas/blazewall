package config

import (
	"fmt"
	"github.com/spf13/viper"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfigFileViper(t *testing.T) {
	viper.SetConfigName("auth-config") // name of config file (without extension)
	viper.AddConfigPath("../..")       // optionally look for config in the working directory
	err := viper.ReadInConfig()        // Find and read the config file
	if err != nil {                    // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	InitConfig()
	ac := GetConfig()
	assert.NotNil(t, ac)
	r := auth.Realms["staff"]
	assert.True(t, len(r.AuthChains) > 0)
	assert.Equal(t, "staff", r.ID)
}
