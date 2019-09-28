package main

import (
	"github.com/google/uuid"
)

type SessionRepository interface {
	CreateSession(userID string, realm string, properties map[string]string)
}

type DummySessionRepository struct {
}

func (sr *DummySessionRepository) CreateSession(userID string, realm string, properties map[string]string) string {
	return uuid.New().String()
}
