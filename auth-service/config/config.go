package config

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
	Name              string            `yaml:"name"`
	RedirectOnSuccess string            `yaml:"redirectOnSuccess"`
	AuthConfig        []AuthConfigEntry `yaml:"authConfig"`
}

type AuthServiceConfig struct {
	Realms        []Realm   `yaml:"realms"`
	CookieDomains []string  `yaml:"cookieDomains"`
	SessionID     string    `yaml:"sessionID"`
	Endpoints     Endpoints `yaml:"endpoints"`
}

type Endpoints struct {
	SessionService string `yaml:"sessionService"`
}

var ac AuthServiceConfig

func Init(reader io.Reader) {
	err := yaml.NewDecoder(reader).Decode(&ac)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func GetConfig() AuthServiceConfig {
	return ac
}
