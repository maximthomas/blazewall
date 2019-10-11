package repo

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/maximthomas/blazewall/gateway-service/config"

	"github.com/maximthomas/blazewall/gateway-service/models"
)

//SessionRepository is generic interface for different implementations
type SessionRepository interface {
	//returns sesson and existing flag
	GetSession(id string) (session models.Session, exists bool)
}

var sr SessionRepository

func GetSessionRepository() SessionRepository {
	return sr
}

func Init() {
	gc := config.GetConfig()

	sr = &RestSessionRepository{
		endpoint: gc.Endpoints.SessionService,
	}
}

//InMemorySessionRepository is a Session repository, stores sessions in memory
type InMemorySessionRepository struct {
	sessions map[string]models.Session
}

func (sr *InMemorySessionRepository) GetSession(id string) (session models.Session, exists bool) {
	session, exists = sr.sessions[id]
	return session, exists
}

//NewInMemorySessionRepository creates new in memory session repository
func NewInMemorySessionRepository(sessions map[string]models.Session) *InMemorySessionRepository {
	return &InMemorySessionRepository{
		sessions: sessions,
	}
}

type RestSessionRepository struct {
	endpoint string
	client   http.Client
}

func (sr *RestSessionRepository) GetSession(id string) (session models.Session, exists bool) {

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
