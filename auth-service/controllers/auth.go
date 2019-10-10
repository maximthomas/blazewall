package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/maximthomas/blazewall/auth-service/models"

	"github.com/maximthomas/blazewall/auth-service/repo"
	"github.com/maximthomas/blazewall/auth-service/service"

	"github.com/gin-gonic/gin"

	"github.com/maximthomas/blazewall/auth-service/config"
)

func ProcessAuth(c *gin.Context) {

	realm, err := getRealmForContext(c)

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "error getting realm"})
		return
	}

	acEntry, ok := getAuthConfigEntry(c, realm.AuthConfig)
	if !ok {
		c.AbortWithStatusJSON(500, gin.H{"error": "error getting auth config"})
		return
	}

	ap, err := service.GetAuthProcessor(acEntry)
	if err != nil {
		log.Fatalf("error getting auth processor %v for AuthConfig entry: %v", err, acEntry)
		c.AbortWithError(500, err)
		return
	}

	session, ok := ap.ProcessAuthentication(c)
	if ok {
		processAuthenticationSuccess(c, session, realm)
	}
}

func processAuthenticationSuccess(c *gin.Context, session models.Session, realm config.Realm) {

	ac := config.GetConfig()
	sr := repo.GetSessionRepo()

	session.Realm = realm.Name
	session, err := sr.CreateSession(session)
	if err != nil {
		log.Fatalf("error creating session: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	for _, domain := range ac.CookieDomains {
		c.SetCookie(ac.SessionID, session.ID, 0, "/", domain, false, true)
	}
	redirect := c.Query("redirect")
	if redirect == "" {
		redirect = realm.RedirectOnSuccess
	}
	c.Redirect(http.StatusFound, redirect)
	c.Writer.WriteHeaderNow()
	return
}

func getRealmForContext(c *gin.Context) (config.Realm, error) {
	ac := config.GetConfig()
	var realm config.Realm
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

func getAuthConfigEntry(c *gin.Context, authConfig []config.AuthConfigEntry) (config.AuthConfigEntry, bool) {
	var acEntry config.AuthConfigEntry

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
