package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

//Session struct represents session object from session service
type Session struct {
	ID         string            `json:"id,omitempty"`
	UserID     string            `json:"userId,omitempty"`
	Realm      string            `json:"realm,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
}

//SessionRepository is generic interface for different implementations
type SessionRepository interface {
	//returns sesson and existing flag
	GetSession(id string) (session Session, exists bool)
}

//InMemorySessionRepository is a Session repository, stores sessions in memory
type InMemorySessionRepository struct {
	sessions map[string]Session
}

func (sr *InMemorySessionRepository) GetSession(id string) (session Session, exists bool) {
	session, exists = sr.sessions[id]
	return session, exists
}

//NewInMemorySessionRepository creates new in memory session repository
func NewInMemorySessionRepository(sessions map[string]Session) *InMemorySessionRepository {
	return &InMemorySessionRepository{
		sessions: sessions,
	}
}

type RestSessionRepository struct {
	endpoint string
	client   http.Client
}

func (sr *RestSessionRepository) GetSession(id string) (session Session, exists bool) {

	resp, err := sr.client.Get(sr.endpoint + "/" + id)
	if err != nil {
		log.Fatalf("error getting session: %v", err)
		return session, exists
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error getting session: %v", err)
		return session, exists
	}

	err = json.Unmarshal(body, &session)
	if err != nil {
		log.Fatalf("error unmarshalling session: %v", err)
		return session, exists
	}
	exists = true
	return session, exists

}
