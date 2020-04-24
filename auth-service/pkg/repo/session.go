package repo

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/maximthomas/blazewall/auth-service/pkg/models"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type SessionRepository interface {
	CreateSession(session models.Session) (models.Session, error)
	DeleteSession(id string) error
	GetSession(id string) (models.Session, error)
	UpdateSession(id string, session models.Session) error
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

func (sr *RestSessionRepository) DeleteSession(id string) error {
	req, err := http.NewRequest("DELETE", sr.Endpoint+"/"+id, nil)
	if err != nil {
		return err
	}
	_, err = sr.client.Do(req)

	return err
}

func (sr *RestSessionRepository) UpdateSession(id string, session models.Session) error {
	return nil
}

type InMemorySessionRepository struct {
	sessions map[string]models.Session
}

func (sr *InMemorySessionRepository) CreateSession(session models.Session) (models.Session, error) {
	if session.ID == "" {
		session.ID = uuid.New().String()
	}
	sr.sessions[session.ID] = session
	return session, nil
}

func (sr *InMemorySessionRepository) DeleteSession(id string) error {
	if _, ok := sr.sessions[id]; ok {
		delete(sr.sessions, id)
		return nil
	} else {
		return errors.New("session does not exist")
	}
}

func (sr *InMemorySessionRepository) GetSession(id string) (models.Session, error) {
	if session, ok := sr.sessions[id]; ok {
		return session, nil
	} else {
		return models.Session{}, errors.New("session does not exist")
	}
}

func (sr *InMemorySessionRepository) UpdateSession(id string, session models.Session) error {
	if _, ok := sr.sessions[id]; ok {
		sr.sessions[id] = session
		return nil
	} else {
		return errors.New("session does not exist")
	}
}

func NewSessionRepository() SessionRepository {
	//ac := config.GetConfig()
	//sr = &RestSessionRepository{Endpoint: ac.Endpoints.SessionService}
	local := os.Getenv("DEV_LOCAL")
	if local == "true" {
		return &InMemorySessionRepository{
			sessions: make(map[string]models.Session),
		}
	}
	return nil
}

func NewInMemorySessionRepository() SessionRepository {
	return &InMemorySessionRepository{
		sessions: make(map[string]models.Session),
	}
}
