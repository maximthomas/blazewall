package main

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

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

func getRealmForContext(c *gin.Context, ac AuthServiceConfig) (Realm, error) {
	var realm Realm
	if len(ac.Realms) == 0 {
		return realm, errors.New("No realm configured")
	}

	realmName := c.Query("realm")
	if realmName == "" {
		return ac.Realms[0], nil
	}
	for _, r := range ac.Realms {
		if r.Name == realmName {
			return r, nil
		}
	}

	return realm, errors.New("No realm found")
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

func processAuthConfig(c *gin.Context, sr SessionRepository, ac AuthServiceConfig) {

	realm, err := getRealmForContext(c, ac)

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "error getting realm"})
		return
	}

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
		for _, domain := range ac.CookieDomains {
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

const secret = "csrf_secret"

func CSRFMiddleware() gin.HandlerFunc {

	getToken := func(tsStr string) string {
		h := sha1.New()
		io.WriteString(h, tsStr+"-"+secret)
		token := base64.URLEncoding.EncodeToString(h.Sum(nil)) + "|" + tsStr
		return token
	}

	return func(c *gin.Context) {
		ts := time.Now().UnixNano() / int64(time.Millisecond)
		tsStr := strconv.FormatInt(ts, 10)
		token := getToken(tsStr)
		c.Set("csrfToken", token)
		if c.Request.Method == "POST" {
			token := c.Request.FormValue("csrfToken")
			if token == "" {
				panic("token not present")
			}
			tokenParts := strings.Split(token, "|")
			if len(tokenParts) != 2 {
				panic("bad token")
			}
			tsStr := tokenParts[1]
			calcToken := getToken(tsStr)
			if token != calcToken {
				panic("bad token!")
			}
		}
		c.Next()
	}
}

func setupRouter(ac AuthServiceConfig, sr SessionRepository) *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "static")
	v1 := router.Group("/auth-service/v1")
	{
		v1.Use(CSRFMiddleware())
		v1.GET("/login", func(c *gin.Context) {
			processAuthConfig(c, sr, ac)
		})
		v1.POST("/login", func(c *gin.Context) {
			processAuthConfig(c, sr, ac)
		})

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
