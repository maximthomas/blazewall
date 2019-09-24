package main

import (
	"flag"

	"github.com/gin-gonic/gin"
)

var port = flag.String("p", "8080", "Session service port")

func main() {

	flag.Parse()

	repo := NewInMemorySessionRepository([]Session{
		{
			ID:     "sess1",
			UserID: "user1",
			Realm:  "users",
		},
	})
	ss := NewSessionService(repo)

	router := gin.Default()

	v1 := router.Group("/session-service/v1")
	{
		session := v1.Group("/sessions")
		{
			session.GET("/:id", ss.getSessionByID)
			session.DELETE("/:id", ss.deleteSession)
			session.POST("/", ss.createSession)
			session.GET("/", ss.findSessions)
		}

	}

	router.Run(":" + *port)
}
