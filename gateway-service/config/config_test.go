package config

import (
	"os"
	"reflect"
	"testing"

	"github.com/maximthomas/blazewall/gateway-service/policy"
)

func TestNewProtectedSitesConfig(t *testing.T) {
	configReader, err := os.Open("./gateway-config-test.yaml")
	if err != nil {
		panic(err)
	}
	Init(configReader)
	gc := GetConfig()

	sitesConfig := gc.ProtectedSitesConfig

	if err != nil {
		panic(err)
	}

	if err != nil {
		t.Errorf("could not get sites config %s", err)
	}

	if len(sitesConfig) != 2 {
		t.Errorf("bad sites config length")
	}

	if sitesConfig[0].RequestHost == "" {
		t.Errorf("could not get sites config %s", err)
	}

	if len(sitesConfig[0].ProtectedPathsConfig) != 2 {
		t.Errorf("could not get protected paths config %s", err)
	}

	if reflect.TypeOf(sitesConfig[0].ProtectedPathsConfig[0].PolicyValidator) != reflect.TypeOf(policy.RealmsPolicyValidator{}) {
		t.Errorf("bad policy validator")
	}
}
