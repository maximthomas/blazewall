package config

import (
	"io"

	"github.com/maximthomas/blazewall/user-service/repo"

	"gopkg.in/yaml.v3"
)

type UserServiceConfig struct {
	RealmRepos map[string]repo.UserRepository
}

var usc UserServiceConfig

func Init(reader io.Reader) {
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
		panic(err)
	}

	usc.RealmRepos = make(map[string]repo.UserRepository, len(config.Realms))
	for _, r := range config.Realms {
		switch r.Type {
		case "mongodb":
			ur := repo.NewUserRepositoryMongoDB(r.Parameters["uri"], r.Parameters["db"], r.Parameters["collection"])
			usc.RealmRepos[r.Name] = &ur

		case "inmemory":
			ur := repo.NewInMemoryUserRepository()
			usc.RealmRepos[r.Name] = ur
		default:
			panic("unknown repo type")
		}

	}
}

func GetUserServiceConfig() UserServiceConfig {
	return usc
}
