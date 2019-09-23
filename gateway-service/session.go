package main

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

func (i *InMemorySessionRepository) GetSession(id string) (session Session, exists bool) {
	session, exists = i.sessions[id]
	return session, exists
}

//NewInMemorySessionRepository creates new in memory session repository
func NewInMemorySessionRepository(sessions map[string]Session) *InMemorySessionRepository {
	return &InMemorySessionRepository{
		sessions: sessions,
	}
}
