package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

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
	CreateSession(session Session) (Session, error)
	DeleteSession(sessionID string) error
}

type DummySessionRepository struct {
}

func (sr *DummySessionRepository) CreateSession(session Session) (Session, error) {
	session.ID = uuid.New().String()
	return session, nil
}

func (sr *DummySessionRepository) DeleteSession(sessionID string) error {
	return nil
}

type RestSessionRepository struct {
	endpoint string
	client   http.Client
}

func (sr *RestSessionRepository) CreateSession(session Session) (Session, error) {
	var newSession Session
	sessBytes, err := json.Marshal(session)
	if err != nil {
		return newSession, err
	}
	buf := bytes.NewBuffer(sessBytes)
	resp, err := sr.client.Post(sr.endpoint, "application/json", buf)
	if err != nil {
		return newSession, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return newSession, err
	}

	err = json.Unmarshal(body, &newSession)
	if err != nil {
		return newSession, err
	}
	return newSession, err
}

func (sr *RestSessionRepository) DeleteSession(sessionID string) error {
	req, err := http.NewRequest("DELETE", sr.endpoint+"/"+sessionID, nil)
	if err != nil {
		return err
	}
	_, err = sr.client.Do(req)

	return err
}
