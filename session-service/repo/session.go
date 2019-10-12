package repo

import "github.com/maximthomas/blazewall/session-service/models"

type SessionRepository interface {
	GetSessionByID(id string) (models.Session, bool)
	DeleteSession(id string) error
	CreateSession(session models.Session) (models.Session, error)
	Find(realm, userID string) []models.Session
}