package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestNewProtectedSitesConfig(t *testing.T) {
	conf := `-
  requestHost: 'http://gateway-service:80'
  targetHost: 'http://protected-resource:8080'
  protectedPathsConfig: 
    -
      urlPattern: '/'
      policyValidator:
        type: realmsPolicyValidator
        settings:
          realms:
            - "staff"
            - "users"
   `
	configReader := strings.NewReader(conf)

	sitesConfig, err := NewProtectedSitesConfigYaml(configReader)

	if sitesConfig == nil || err != nil {
		t.Errorf("could not get sites config %s", err)
	}
}

func TestWriteCong(t *testing.T) {
	/*
		sc := []ProtectedSiteConfig{
			{
				RequestHost: "proxyHost",
				TargetHost:  "test",
				ProtectedPathsConfig: []ProtectedPathConfig{
					{
						URLPattern:      "/",
						PolicyValidator: AllowedPolicyValidator{},
					},
					{
						URLPattern:      "/protected",
						PolicyValidator: DeniedPolicyValidator{},
					},
				},
				proxy: nil,
			},
		}*/

	type yamlPolicyValidator struct {
		PolicyType string                 `yaml:"policyType"`
		Settings   map[string]interface{} `yaml:"settings"`
	}

	type yamlProtectedPaths struct {
		URLPattern      string              `yaml:"urlPattern"`
		PolicyValidator yamlPolicyValidator `yaml:"policyValidator"`
	}

	type yamlConfig struct {
		RequestHost          string               `yaml:"requestHost"`
		TargetHost           string               `yaml:"targetHost"`
		ProtectedPathsConfig []yamlProtectedPaths `yaml:"protectedPathsConfig"`
	}
	cfg := []yamlConfig{
		{
			RequestHost: "te",
			TargetHost:  "asdas",
			ProtectedPathsConfig: []yamlProtectedPaths{
				{
					URLPattern: "patt",
					PolicyValidator: yamlPolicyValidator{
						PolicyType: "RealmPolicy",
						Settings: map[string]interface{}{
							"realms": []string{"staff", "users"},
						},
					},
				},
			},
		},
	}

	var b strings.Builder
	yaml.NewEncoder(&b).Encode(cfg)
	bstr := b.String()
	fmt.Println(bstr)

	configReader := strings.NewReader(bstr)

	var config []yamlConfig
	err := yaml.NewDecoder(configReader).Decode(&config)
	p := config[0].ProtectedPathsConfig[0].PolicyValidator.Settings["realms"]
	rv := reflect.ValueOf(p)
	fmt.Println(rv.Kind())
	if err != nil {
		t.Errorf("could not get sites config %s", err)
	}
}
