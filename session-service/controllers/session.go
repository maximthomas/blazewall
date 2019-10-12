package controllers

import (
	"github.com/maximthomas/blazewall/session-service/models"
	"github.com/maximthomas/blazewall/session-service/repo"
	"log"

	"github.com/gin-gonic/gin"
)

type SessionService struct {
	sr repo.SessionRepository
}

func (ss *SessionService) GetSessionByID(c *gin.Context) {
	id := c.Param("id")
	log.Printf("getting session by id: %s", id)
	session, ok := ss.sr.GetSessionByID(id)
	if !ok {
		c.JSON(404, gin.H{"error": "Session not found"})
	} else {
		c.JSON(200, session)
	}
}

func (ss *SessionService) FindSessions(c *gin.Context) {
	realm := c.Query("realm")
	userID := c.Query("userID")
	log.Printf("findins sessions by realm: %s and userID: %s", realm, userID)

	if realm == "" || userID == "" {
		c.JSON(400, gin.H{"error": "Realm and userId not set"})
		return
	}

	sessions := ss.sr.Find(realm, userID)
	if len(sessions) == 0 {
		c.JSON(404, gin.H{"error": "Sessions not found"})
	} else {
		c.JSON(200, sessions)
	}
}

func (ss *SessionService) DeleteSession(c *gin.Context) {
	id := c.Param("id")
	log.Printf("deleting session by id: %s", id)
	err := ss.sr.DeleteSession(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Session not found"})
	} else {
		c.JSON(202, gin.H{"message": "Accepted"})
	}
}

func (ss *SessionService) CreateSession(c *gin.Context) {
	session := models.Session{}
	log.Printf("creating session")
	err := c.ShouldBindJSON(&session)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
		return
	}
	newSession, createErr := ss.sr.CreateSession(session)
	if createErr != nil {
		c.JSON(500, gin.H{"error": createErr})
	} else {
		c.JSON(200, newSession)
	}
}

func NewSessionService(sr repo.SessionRepository) SessionService {
	return SessionService{
		sr: sr,
	}
}
