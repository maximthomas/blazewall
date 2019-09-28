package main

import (
	"flag"
	"os"

	"github.com/gin-gonic/gin"
)

var yamlConfigFile = flag.String("yc", "", "Yaml config file path")
var port = flag.String("p", "8080", "Gateway service port")

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func getAuthConfigEntry(c *gin.Context, authConfig []AuthConfigEntry) (AuthConfigEntry, bool) {
	var acEntry AuthConfigEntry

	if len(authConfig) == 0 {
		return acEntry, false
	}

	authType := c.Query("authType")
	if authType == "" {
		acEntry = authConfig[0]
	} else {
		for _, acE := range authConfig {
			if authType == acE.Type {
				acEntry = acE
				break
			}
		}
	}
	return acEntry, true
}

const plainAuthType = "textFile"

func processAuthConfig(c *gin.Context, authConfig []AuthConfigEntry) {
	acEntry, ok := getAuthConfigEntry(c, authConfig)
	if !ok {
		c.AbortWithStatusJSON(500, gin.H{"error": "error getting auth config"})
		return
	}

	ap, err := getAuthProcessor(acEntry)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	ap.ProcessAuthentication(c)
}

func main() {
	flag.Parse()

	var ac AuthServiceConfig
	if *yamlConfigFile != "" {
		configReader, err := os.Open(*yamlConfigFile)
		check(err)

		ac, err = NewAuthServiceConfigYaml(configReader)
		check(err)
	}

	router := gin.Default()

	v1 := router.Group("/auth-service/v1")
	{
		for _, realm := range ac.Realms {
			v1.GET("/"+realm.Name, func(c *gin.Context) {
				processAuthConfig(c, realm.AuthConfig)
			})

			v1.POST("/"+realm.Name, func(c *gin.Context) {
				processAuthConfig(c, realm.AuthConfig)
			})
		}
	}

	router.Run(":" + *port)
}
