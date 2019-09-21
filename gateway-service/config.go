package main

import (
	"errors"
	"fmt"
	"io"
	"net/http/httputil"
	"net/url"

	"github.com/spf13/cast"
	"gopkg.in/yaml.v2"
)

type ProtectedSiteConfig struct {
	RequestHost          string
	TargetHost           string
	ProtectedPathsConfig []ProtectedPathConfig
	proxy                *httputil.ReverseProxy
}

type ProtectedPathConfig struct {
	URLPattern      string
	PolicyValidator PolicyValidator
}

func NewProtectedSitesConfigYaml(reader io.Reader) ([]ProtectedSiteConfig, error) {

	type yamlPolicyValidator struct {
		PolicyType string                 `yaml:"type"`
		Settings   map[string]interface{} `yaml:"settings"`
	}

	type yamlProtectedPath struct {
		URLPattern      string              `yaml:"urlPattern"`
		PolicyValidator yamlPolicyValidator `yaml:"policyValidator"`
	}

	type yamlConfig struct {
		RequestHost string              `yaml:"requestHost"`
		TargetHost  string              `yaml:"targetHost"`
		PathsConfig []yamlProtectedPath `yaml:"pathsConfig"`
	}

	getPolicyValidator := func(yp yamlPolicyValidator) (PolicyValidator, error) {

		switch yp.PolicyType {
		case "allowed":
			return AllowedPolicyValidator{}, nil
		case "denied":
			return DeniedPolicyValidator{}, nil
		case "realms":
			realms, err := cast.ToStringSliceE(yp.Settings["realms"])
			if err != nil {
				return nil, err
			}
			return RealmsPolicyValidator{Realms: realms}, nil
		}
		return nil, errors.New("Undefined policy type: " + yp.PolicyType)
	}

	var config []yamlConfig
	err := yaml.NewDecoder(reader).Decode(&config)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	result := make([]ProtectedSiteConfig, len(config))

	for confIdx, configEntry := range config {

		paths := make([]ProtectedPathConfig, len(configEntry.PathsConfig))

		for pathIdx, pathEntry := range configEntry.PathsConfig {
			pv, err := getPolicyValidator(pathEntry.PolicyValidator)
			if err != nil {
				return nil, err
			}

			paths[pathIdx] = ProtectedPathConfig{
				URLPattern:      pathEntry.URLPattern,
				PolicyValidator: pv,
			}
		}

		targetURL, err := url.Parse(configEntry.TargetHost)
		if err != nil {
			return nil, err
		}

		result[confIdx] = ProtectedSiteConfig{
			RequestHost:          configEntry.RequestHost,
			TargetHost:           configEntry.TargetHost,
			ProtectedPathsConfig: paths,
			proxy:                httputil.NewSingleHostReverseProxy(targetURL),
		}
	}
	return result, nil
}
