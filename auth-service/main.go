package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var yamlConfigFile = flag.String("yc", "./test/auth-config.yaml", "Yaml config file path")
var port = flag.String("p", "8080", "Gateway service port")
var sessionServiceEndpoint = flag.String("sess", "http://session-service:8080/session-service/v1/sessions", "Session service endpoint")
var authSessionID = flag.String("sID", "BlazewallSession", "Session service cookie name")

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

func processAuthConfig(c *gin.Context, realm Realm, sr SessionRepository, cookieDomains []string) {
	acEntry, ok := getAuthConfigEntry(c, realm.AuthConfig)
	if !ok {
		c.AbortWithStatusJSON(500, gin.H{"error": "error getting auth config"})
		return
	}

	ap, err := getAuthProcessor(acEntry)
	if err != nil {
		log.Fatalf("error getting auth processor %v for AuthConfig entry: %v", err, acEntry)
		c.AbortWithError(500, err)
		return
	}
	session, ok := ap.ProcessAuthentication(c)
	if ok {
		session.Realm = realm.Name
		session, err = sr.CreateSession(session)
		if err != nil {
			log.Fatalf("error creating session: %v", err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		for _, domain := range cookieDomains {
			c.SetCookie(*authSessionID, session.ID, 0, "/", domain, false, true)
		}
		redirect := c.Query("redirect")
		if redirect == "" {
			redirect = realm.RedirectOnSuccess
		}
		c.Redirect(http.StatusFound, redirect)
		c.Writer.WriteHeaderNow()
		return
	}
}

func processLogout(c *gin.Context, sr SessionRepository) {
	sessionID, err := c.Cookie(*authSessionID)
	if err != nil {
		return
	}
	sr.DeleteSession(sessionID)
}

func setupRouter(ac AuthServiceConfig, sr SessionRepository) *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "static")
	v1 := router.Group("/auth-service/v1")
	{
		for _, realm := range ac.Realms {
			r := realm
			v1.GET("/"+realm.Name, func(c *gin.Context) {
				processAuthConfig(c, r, sr, ac.CookieDomains)
			})

			v1.POST("/"+realm.Name, func(c *gin.Context) {
				processAuthConfig(c, r, sr, ac.CookieDomains)
			})
		}
		v1.GET("/logout", func(c *gin.Context) {
			processLogout(c, sr)
		})
	}
	return router
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

	//sr := DummySessionRepository{}
	sr := RestSessionRepository{endpoint: *sessionServiceEndpoint}
	router := setupRouter(ac, &sr)
	router.Run(":" + *port)
}
