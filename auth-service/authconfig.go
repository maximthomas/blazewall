package main

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

type AuthConfigEntry struct {
	Type       string                  `yaml:"type"`
	Parameters map[string]*interface{} `yaml:"parameters"`
}

type Realm struct {
	Name       string            `yaml:"name"`
	AuthConfig []AuthConfigEntry `yaml:"authConfig"`
}

type AuthServiceConfig struct {
	Realms        []Realm  `yaml:"realms"`
	CookieDomains []string `yaml:"cookieDomains"`
}

func NewAuthServiceConfigYaml(reader io.Reader) (AuthServiceConfig, error) {

	var config AuthServiceConfig
	err := yaml.NewDecoder(reader).Decode(&config)
	if err != nil {
		fmt.Println(err)
		return config, err
	}

	return config, nil
}
