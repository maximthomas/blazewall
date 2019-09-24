package main

import (
	"github.com/gin-gonic/gin"
)

type SessionService struct {
	sr SessionRepository
}

func (ss *SessionService) getSessionByID(c *gin.Context) {
	id := c.Param("id")
	session, ok := ss.sr.GetSessionByID(id)
	if !ok {
		c.JSON(404, gin.H{"error": "Session not found"})
	} else {
		c.JSON(200, session)
	}
}

func (ss *SessionService) findSessions(c *gin.Context) {
	realm := c.Query("realm")
	userID := c.Query("userId")

	if realm == "" || userID == "" {
		c.JSON(400, gin.H{"error": "Realm and userId not set"})
		return
	}

	sessions := ss.sr.FindByUserId(realm, userID)
	if len(sessions) == 0 {
		c.JSON(404, gin.H{"error": "Sessions not found"})
	} else {
		c.JSON(200, sessions)
	}
}

func (ss *SessionService) deleteSession(c *gin.Context) {
	id := c.Param("id")
	err := ss.sr.DeleteSession(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Session not found"})
	} else {
		c.JSON(202, gin.H{"message": "Accepted"})
	}
}

func (ss *SessionService) createSession(c *gin.Context) {
	session := Session{}
	err := c.ShouldBindJSON(&session)
	if err != nil {
		c.JSON(500, gin.H{"error": err})
	}
	newSession, createErr := ss.sr.CreateSession(session)
	if createErr != nil {
		c.JSON(500, gin.H{"error": createErr})
	} else {
		c.JSON(200, newSession)
	}
}

func NewSessionService(sr SessionRepository) SessionService {
	return SessionService{
		sr: sr,
	}
}
