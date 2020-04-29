package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/mitchellh/mapstructure"
	"log"

	"github.com/google/uuid"

	"github.com/maximthomas/blazewall/auth-service/pkg/repo"
	"github.com/spf13/viper"

	"github.com/sirupsen/logrus"
)

type Config struct {
	
}

type Authentication struct {
	Realms map[string]Realm `yaml:"realms"`
	Logger logrus.FieldLogger
}

type Realm struct {
	ID            string
	Modules       map[string]Module    `yaml:"modules"`
	AuthChains    map[string]AuthChain `yaml:"authChains"`
	UserDataStore UserDataStore        `yaml:"userDataStore"`
	UserRepo      repo.UserRepository
	Session       Session `yaml:"session"`
}

type AuthChain struct {
	Modules []ChainModule `yaml:"modules"`
}

type UserDataStore struct {
	Type       string                 `yaml:"type"`
	Properties map[string]interface{} `yaml:"properties,omitempty"`
}

type Module struct {
	Type       string                 `yaml:"type"`
	Properties map[string]interface{} `yaml:"properties,omitempty"`
}

type ChainModule struct {
	ID         string                 `yaml:"id"`
	Properties map[string]interface{} `yaml:"properties,omitempty"`
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
	logger := logrus.New()
	//newLogger.SetFormatter(&logrus.JSONFormatter{})
	//newLogger.SetReportCaller(true)
	auth.Logger = logger
	var configLogger = logger.WithField("module", "config")

	err := viper.UnmarshalKey("authentication", &auth)
	if err != nil { // Handle errors reading the config file
		configLogger.Errorf("Fatal error config file: %s \n", err)
		panic(err)
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

		if realm.UserDataStore.Type == "ldap" {
			prop := realm.UserDataStore.Properties
			repo := &repo.UserLdapRepository{}
			mapstructure.Decode(prop, repo)
			realm.UserRepo = repo
		} else {
			realm.UserRepo = repo.NewInMemoryUserRepository()
		}
		auth.Realms[id] = realm
	}
	configLogger.Infof("got configuration %+v", auth)

	return nil
}

func GetConfig() Authentication {
	return auth
}
