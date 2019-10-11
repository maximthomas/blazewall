package config

import (
	"errors"
	"io"
	"log"
	"net/http/httputil"
	"net/url"

	"github.com/maximthomas/blazewall/gateway-service/policy"

	"github.com/spf13/cast"
	"gopkg.in/yaml.v3"
)

type GatewayConfig struct {
	ProtectedSitesConfig []ProtectedSiteConfig
	SessionID            string    `yaml:"sessionID"`
	Endpoints            Endpoints `yaml:"endpoints"`
}

type Endpoints struct {
	SessionService string `yaml:"sessionService"`
}

type ProtectedSiteConfig struct {
	RequestHost          string
	TargetHost           string
	ProtectedPathsConfig []ProtectedPathConfig
	Proxy                *httputil.ReverseProxy
}

type ProtectedPathConfig struct {
	URLPattern      string
	PolicyValidator policy.PolicyValidator
	AuthURL         string
}

var gc GatewayConfig

func GetConfig() GatewayConfig {
	return gc
}

func Init(reader io.Reader) {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	type yamlPolicyValidator struct {
		PolicyType string                 `yaml:"type"`
		Settings   map[string]interface{} `yaml:"settings"`
	}

	type yamlProtectedPath struct {
		URLPattern      string              `yaml:"urlPattern"`
		PolicyValidator yamlPolicyValidator `yaml:"policyValidator"`
		AuthURL         string              `yaml:"authUrl"`
	}

	type yamlSiteConfig struct {
		RequestHost string              `yaml:"requestHost"`
		TargetHost  string              `yaml:"targetHost"`
		PathsConfig []yamlProtectedPath `yaml:"pathsConfig"`
	}

	type yamlGatewayConfig struct {
		YamlSitesConfig []yamlSiteConfig `yaml:"protectedHosts"`
		SessionID       string           `yaml:"sessionID"`
		Endpoints       Endpoints        `yaml:"endpoints"`
	}

	getPolicyValidator := func(yp yamlPolicyValidator) (policy.PolicyValidator, error) {

		switch yp.PolicyType {
		case "allowed":
			return policy.AllowedPolicyValidator{}, nil
		case "denied":
			return policy.DeniedPolicyValidator{}, nil
		case "authenticated":
			return policy.AuthenticatedUserPolicyValidator{}, nil
		case "realms":
			realms, err := cast.ToStringSliceE(yp.Settings["realms"])
			if err != nil {
				return nil, err
			}
			return policy.RealmsPolicyValidator{Realms: realms}, nil
		}
		return nil, errors.New("Undefined policy type: " + yp.PolicyType)
	}

	var result GatewayConfig

	var config yamlGatewayConfig
	err := yaml.NewDecoder(reader).Decode(&config)
	if err != nil {
		panic(err)
	}

	result.SessionID = config.SessionID
	result.Endpoints = config.Endpoints
	result.ProtectedSitesConfig = make([]ProtectedSiteConfig, len(config.YamlSitesConfig))

	for confIdx, configEntry := range config.YamlSitesConfig {

		paths := make([]ProtectedPathConfig, len(configEntry.PathsConfig))

		for pathIdx, pathEntry := range configEntry.PathsConfig {
			pv, err := getPolicyValidator(pathEntry.PolicyValidator)
			if err != nil {
				panic(err)
			}

			paths[pathIdx] = ProtectedPathConfig{
				URLPattern:      pathEntry.URLPattern,
				PolicyValidator: pv,
				AuthURL:         pathEntry.AuthURL,
			}
		}

		targetURL, err := url.Parse(configEntry.TargetHost)
		if err != nil {
			panic(err)
		}

		result.ProtectedSitesConfig[confIdx] = ProtectedSiteConfig{
			RequestHost:          configEntry.RequestHost,
			TargetHost:           configEntry.TargetHost,
			ProtectedPathsConfig: paths,
			Proxy:                httputil.NewSingleHostReverseProxy(targetURL),
		}
	}
	gc = result
}
