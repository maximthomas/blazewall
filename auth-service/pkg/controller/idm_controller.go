package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maximthomas/blazewall/auth-service/pkg/config"
	"github.com/maximthomas/blazewall/auth-service/pkg/repo"
	"github.com/sirupsen/logrus"
)

type IDMController struct {
	sr     repo.SessionRepository
	logger logrus.FieldLogger
}

func NewIDMController(config config.Config) *IDMController {
	logger := config.Logger.WithField("module", "IDMController")
	sr := config.SessionDataStore.Repo
	return &IDMController{sr, logger}
}

func (ic IDMController) Profile(c *gin.Context) {
	sessID := getSessionIdFromRequest(c)
	if sessID == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
	} else {
		s, err := ic.sr.GetSession(sessID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		} else {
			c.JSON(http.StatusOK, s)
		}
	}
}
