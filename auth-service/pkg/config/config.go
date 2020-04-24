package config

import (
	"fmt"
	"github.com/maximthomas/blazewall/auth-service/pkg/repo"
	"github.com/spf13/viper"
)

type Authentication struct {
	Realms map[string]Realm `yaml:"realms"`
}

type Realm struct {
	ID         string
	Modules    map[string]Module    `yaml:"modules"`
	AuthChains map[string]AuthChain `yaml:"authChains"`
	DataStore  DataStore            `yaml:"datastore"`
	UserRepo   repo.UserRepository
}

type AuthChain struct {
	Modules []ChainModule `yaml:"modules"`
}

type DataStore struct {
	Type          string `yaml:"type"`
	URL           string `yaml:"url"`
	Authorization string `yaml:"authorization"`
}

type Module struct {
	Type       string            `yaml:"type"`
	Properties map[string]string `yaml:"properties,omitempty"`
}

type ChainModule struct {
	ID         string            `yaml:"id"`
	Properties map[string]string `yaml:"properties,omitempty"`
}

var auth Authentication

func InitConfig() {
	err := viper.UnmarshalKey("authentication", &auth)
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	for id, realm := range auth.Realms {
		realm.UserRepo = repo.NewInMemoryUserRepository()
		realm.ID = id
		auth.Realms[id] = realm
	}
}

func GetConfig() Authentication {
	return auth
}
