package server

import (
	"github.com/gin-gonic/gin"
	"github.com/maximthomas/blazewall/auth-service/controllers"
	"github.com/maximthomas/blazewall/auth-service/middleware"
)

func GetRouter() *gin.Engine {
	router := gin.Default()

	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "static")
	v1 := router.Group("/auth-service/v1")
	{
		v1.Use(middleware.CSRF())
		v1.GET("/login", func(c *gin.Context) {
			controllers.ProcessAuth(c)
		})
		v1.POST("/login", func(c *gin.Context) {
			controllers.ProcessAuth(c)
		})

		v1.GET("/logout", func(c *gin.Context) {
			controllers.ProcessLogout(c)
		})
	}
	return router
}
