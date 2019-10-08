package main

import (
	"errors"

	"github.com/google/uuid"
)

type Session struct {
	ID string `json:"id,omitempty"`

	UserID string `json:"userId,omitempty"`

	Realm string `json:"realm,omitempty"`

	Expired int64 `json:"expired,omitempty"`

	Properties map[string]string `json:"properties,omitempty"`
}

type SessionRepository interface {
	GetSessionByID(id string) (Session, bool)
	DeleteSession(id string) error
	CreateSession(session Session) (Session, error)
	Find(realm, userID string) []Session
}

type InMemorySessionRepository []Session

func (repo *InMemorySessionRepository) GetSessionByID(id string) (Session, bool) {
	var s Session
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

func (repo *InMemorySessionRepository) CreateSession(session Session) (Session, error) {
	if session.ID == "" {
		session.ID = uuid.New().String()
	}

	*repo = append(*repo, session)
	return session, nil
}

func (repo *InMemorySessionRepository) Find(realm, userID string) []Session {
	var sessions []Session
	for _, el := range *repo {
		if el.UserID == userID && el.Realm == realm {
			sessions = append(sessions, el)
		}
	}
	return sessions
}

func NewInMemorySessionRepository(sessions []Session) *InMemorySessionRepository {
	repo := new(InMemorySessionRepository)
	*repo = sessions
	return repo
}
