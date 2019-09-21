package main

type Session struct {
	ID         string            `json:"id,omitempty"`
	UserID     string            `json:"userId,omitempty"`
	Realm      string            `json:"realm,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
}

type SessionRepository interface {
	//returns sesson and existing flag
	GetSession(id string) (session Session, exists bool)
}

type InMemorySessionRepository struct {
	sessions map[string]Session
}

func (i *InMemorySessionRepository) GetSession(id string) (session Session, exists bool) {
	session, exists = i.sessions[id]
	return session, exists
}

func NewInMemorySessionRepository(sessions map[string]Session) *InMemorySessionRepository {
	return &InMemorySessionRepository{
		sessions: sessions,
	}
}
