package repo

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/maximthomas/blazewall/auth-service/config"

	"github.com/google/uuid"

	"github.com/maximthomas/blazewall/auth-service/models"
)

type SessionRepository interface {
	CreateSession(session models.Session) (models.Session, error)
	DeleteSession(sessionID string) error
}

type RestSessionRepository struct {
	Endpoint string
	client   http.Client
}

func (sr *RestSessionRepository) CreateSession(session models.Session) (models.Session, error) {
	var newSession models.Session
	sessBytes, err := json.Marshal(session)
	if err != nil {
		return newSession, err
	}
	buf := bytes.NewBuffer(sessBytes)
	resp, err := sr.client.Post(sr.Endpoint, "application/json", buf)
	if err != nil {
		log.Printf("error creating session: %v", err)
		return newSession, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error creating session: %v", err)
		return newSession, err
	}

	err = json.Unmarshal(body, &newSession)
	if err != nil {
		log.Printf("error creating session: %v", err)
		return newSession, err
	}
	log.Printf("created new session: %v", newSession)
	return newSession, err
}

func (sr *RestSessionRepository) DeleteSession(sessionID string) error {
	req, err := http.NewRequest("DELETE", sr.Endpoint+"/"+sessionID, nil)
	if err != nil {
		return err
	}
	_, err = sr.client.Do(req)

	return err
}

type DummySessionRepository struct {
}

func (sr *DummySessionRepository) CreateSession(session models.Session) (models.Session, error) {
	session.ID = uuid.New().String()
	return session, nil
}

func (sr *DummySessionRepository) DeleteSession(sessionID string) error {
	return nil
}

var sr SessionRepository

func InitSessionRepo() {
	ac := config.GetConfig()
	sr = &RestSessionRepository{Endpoint: ac.Endpoints.SessionService}
	local := os.Getenv("DEV_LOCAL")
	if local == "true" {
		sr = &DummySessionRepository{}
	}
}

func GetSessionRepo() SessionRepository {
	return sr
}
