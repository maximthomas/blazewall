package cmd

import (
	"testing"

	"github.com/maximthomas/blazewall/auth-service/pkg/config"

	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	args := []string{"version", "--config", "../auth-config.yaml"}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	assert.NoError(t, err)
	conf := config.GetConfig()
	r := conf.Authentication.Realms["staff"]
	assert.True(t, len(r.AuthChains) > 0)

}
