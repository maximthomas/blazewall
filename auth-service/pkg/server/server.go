package pkg

import (
	"github.com/gin-gonic/gin"
)

func login(c *gin.Context) {
	c.JSON(404, gin.H{"error": "Session not found"})
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("/auth-service/v1")
	{
		session := v1.Group("/login")
		{
			session.GET("/", login)
		}
	}
	return router
}

func RunServer() {
	router := setupRouter()
	router.Run(":" + "8080")
}
