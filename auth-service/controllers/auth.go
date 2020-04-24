package controllers

import (
	"github.com/gin-gonic/gin"

	"github.com/maximthomas/blazewall/auth-service/pkg/config"
)

func ProcessAuth(c *gin.Context) {

}

func getRealmForContext(c *gin.Context) (config.Realm, error) {
	return config.Realm{}, nil
}

func getAuthConfigEntry(c *gin.Context) (interface{}, bool) {
	return nil, false
}
