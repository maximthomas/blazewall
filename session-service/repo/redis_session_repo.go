package repo


import (
	"encoding/json"
	"github.com/maximthomas/blazewall/session-service/models"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/go-redis/redis"
)

type SessionRepositoryRedis struct {
	client *redis.Client
}

func (repo *SessionRepositoryRedis) GetSessionByID(id string) (models.Session, bool) {
	var session models.Session
	val, err := repo.client.Get(id).Result()
	if err != nil {
		return session, false
	}

	err = json.Unmarshal([]byte(val), &session)
	if err != nil {
		log.Printf("Error unmarshalling session %v", err)
		panic(err)
	}
	return session, true
}

func (repo *SessionRepositoryRedis) DeleteSession(id string) error {
	_, err := repo.client.Get(id).Result()
	if err != nil {
		return err
	}
	_, err = repo.client.Del(id).Result()
	return err
}

func (repo *SessionRepositoryRedis) CreateSession(session models.Session) (models.Session, error) {
	if session.ID == "" {
		session.ID = uuid.New().String()
	}
	sessionBytes, err := json.Marshal(session)
	if err != nil {
		return session, err
	}

	duration := time.Unix(0, session.Expired*int64(time.Millisecond)).Sub(time.Now())

	_, err = repo.client.Set(session.ID, string(sessionBytes), duration).Result()
	if err != nil {
		return session, err
	}

	return session, nil
}

func (repo *SessionRepositoryRedis) Find(realm, userID string) []models.Session {
	panic("Method not implemented")
}

func NewSessionRepositoryRedis(addr, pass string, db int) SessionRepositoryRedis {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass, // no password set
		DB:       db,   // use default DB
	})

	_, err := client.Ping().Result()

	if err != nil {
		panic(err)
	}

	return SessionRepositoryRedis{
		client: client,
	}
}
