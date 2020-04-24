package cmd

import (
	"github.com/maximthomas/blazewall/auth-service/pkg/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	args := []string{"version", "--config", "../auth-config.yaml"}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	assert.NoError(t, err)
	ac := config.GetConfig()
	r := ac.Realms["staff"]
	assert.True(t, len(r.AuthChains) > 0)

}
