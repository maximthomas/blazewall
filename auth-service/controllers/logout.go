package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/maximthomas/blazewall/auth-service/config"
	"github.com/maximthomas/blazewall/auth-service/repo"
)

func ProcessLogout(c *gin.Context) {
	ac := config.GetConfig()
	sessionID, err := c.Cookie(ac.SessionID)
	if err != nil {
		return
	}
	sr := repo.GetSessionRepo()
	sr.DeleteSession(sessionID)
}
