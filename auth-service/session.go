package main

import (
	"github.com/google/uuid"
)

//Session struct represents session object from session service
type Session struct {
	ID         string            `json:"id,omitempty"`
	UserID     string            `json:"userId,omitempty"`
	Realm      string            `json:"realm,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
}

type SessionRepository interface {
	CreateSession(session Session) Session
}

type DummySessionRepository struct {
}

func (sr *DummySessionRepository) CreateSession(session Session) Session {
	session.ID = uuid.New().String()
	return session
}
