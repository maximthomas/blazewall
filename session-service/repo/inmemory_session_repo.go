package repo

import (
	"errors"
	"github.com/google/uuid"
	"github.com/maximthomas/blazewall/session-service/models"
)

type InMemorySessionRepository []models.Session

func (repo *InMemorySessionRepository) GetSessionByID(id string) (models.Session, bool) {
	var s models.Session
	for _, el := range *repo {
		if el.ID == id {
			return el, true
		}
	}
	return s, false
}

func (repo *InMemorySessionRepository) DeleteSession(id string) error {

	idx := -1
	for i, el := range *repo {
		if el.ID == id {
			idx = i
			break
		}
	}

	if idx == -1 {
		return errors.New("Session id not found")
	}
	*repo = append((*repo)[:idx], (*repo)[idx+1:]...)
	return nil
}

func (repo *InMemorySessionRepository) CreateSession(session models.Session) (models.Session, error) {
	if session.ID == "" {
		session.ID = uuid.New().String()
	}

	*repo = append(*repo, session)
	return session, nil
}

func (repo *InMemorySessionRepository) Find(realm, userID string) []models.Session {
	var sessions []models.Session
	for _, el := range *repo {
		if el.UserID == userID && el.Realm == realm {
			sessions = append(sessions, el)
		}
	}
	return sessions
}

func NewInMemorySessionRepository(sessions []models.Session) *InMemorySessionRepository {
	repo := new(InMemorySessionRepository)
	*repo = sessions
	return repo
}
