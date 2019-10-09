package main

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

type UserServiceConfig struct {
	RealmRepos map[string]UserRepository
}

func NewUserServiceConfigYaml(reader io.Reader) (UserServiceConfig, error) {

	var userServiceConfig UserServiceConfig

	type yamlRealm struct {
		Name       string            `yaml:"realm"`
		Type       string            `yaml:"type"`
		Parameters map[string]string `yaml:"parameters"`
	}
	type yamlRealms struct {
		Realms []yamlRealm `yaml:"realms"`
	}

	var config yamlRealms
	err := yaml.NewDecoder(reader).Decode(&config)
	if err != nil {
		fmt.Println(err)
		return userServiceConfig, err
	}

	userServiceConfig.RealmRepos = make(map[string]UserRepository, len(config.Realms))
	for _, r := range config.Realms {
		if r.Type == "mongodb" {
			ur := NewUserRepositoryMongoDB(r.Parameters["uri"], r.Parameters["db"], r.Parameters["collection"])
			userServiceConfig.RealmRepos[r.Name] = &ur
		}
	}

	return userServiceConfig, nil
}
