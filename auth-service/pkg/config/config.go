package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"

	"github.com/mitchellh/mapstructure"

	"github.com/google/uuid"

	"github.com/maximthomas/blazewall/auth-service/pkg/repo"
	"github.com/spf13/viper"

	"github.com/sirupsen/logrus"
)

type Config struct {
	Authentication Authentication
	Session        SessionSettings
	Logger         logrus.FieldLogger
}

type Authentication struct {
	Realms map[string]Realm `yaml:"realms"`
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

type SessionSettings struct {
	Repo       repo.SessionRepository
	Type       string
	Properties map[string]string
}

var config Config

func InitConfig() error {
	logger := logrus.New()
	//newLogger.SetFormatter(&logrus.JSONFormatter{})
	//newLogger.SetReportCaller(true)
	var configLogger = logger.WithField("module", "config")

	err := viper.Unmarshal(&config)
	auth := &config.Authentication

	config.Logger = logger
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

func GetAuthConfig() Authentication {
	return config.Authentication
}

func GetConfig() Config {
	return config
}
