package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"

	"github.com/google/uuid"

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
	Session    Session `yaml:"session"`
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

type Session struct {
	Type    string     `yaml:"type"`
	Expires int        `yaml:"expires"`
	Jwt     SessionJWT `yaml:"jwt"`
}

type SessionJWT struct {
	Issuer        string `yaml:"issuer"`
	PrivateKeyPem string `yml:"privateKeyPem"`
	PrivateKeyID  string
	PrivateKey    *rsa.PrivateKey
	PublicKey     *rsa.PublicKey
}

var auth Authentication

func InitConfig() error {
	err := viper.UnmarshalKey("authentication", &auth)
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	for id, realm := range auth.Realms {
		realm.UserRepo = repo.NewInMemoryUserRepository()
		realm.ID = id
		if realm.Session.Type == "stateless" {
			jwt := &realm.Session.Jwt
			privateKeyBlock, _ := pem.Decode([]byte(jwt.PrivateKeyPem))
			privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
			if err != nil {
				log.Fatal(err)
				return err
			}
			jwt.PrivateKey = privateKey

			jwt.PublicKey = &privateKey.PublicKey
			jwt.PrivateKeyID = uuid.New().String()

		}
		auth.Realms[id] = realm
	}
	return nil
}

func GetConfig() Authentication {
	return auth
}
