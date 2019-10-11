package main

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

var port = flag.String("p", "8080", "Session service port")

func setupRouter(ss *SessionService) *gin.Engine {
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
	return router
}

func main() {

	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Llongfile)
	redisAddr := getEnv("REDIS_ADDRES", "localhost:6379")
	redisPass := getEnv("REDIS_PASS", "")
	redisDB := getEnvAsInt("REDIS_PASS", 0)

	sr := NewSessionRepositoryRedis(redisAddr, redisPass, redisDB)
	ss := NewSessionService(&sr)

	router := setupRouter(&ss)
	router.Run(":" + *port)
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

// Simple helper function to read an environment variable into integer or return a default value
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}
