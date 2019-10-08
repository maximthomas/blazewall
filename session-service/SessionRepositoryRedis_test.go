package main

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
)

const existingSessionID = "a5624f6e-63a8-4702-a6fa-3f4e001f61c7"

var existingRedisSesson = Session{
	ID:         existingSessionID,
	UserID:     "user1",
	Realm:      "users",
	Expired:    (time.Now().UnixNano() / int64(time.Millisecond)) + 60*60*24,
	Properties: map[string]string{"foo": "bar", "abc": "bac"},
}

func TestRedisGetSession(t *testing.T) {

	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	sessionBytes, _ := json.Marshal(existingSesson)
	s.Set(existingSessionID, string(sessionBytes))

	sr := NewSessionRepositoryRedis(s.Addr(), "", 0)

	t.Run("get not existing session", func(t *testing.T) {
		_, exists := sr.GetSessionByID("bad")
		assert.False(t, exists)
	})

	t.Run("get existing session", func(t *testing.T) {
		gotSession, exists := sr.GetSessionByID(existingSessionID)
		assert.True(t, exists)
		assert.Equal(t, existingSesson, gotSession)
	})
}

func TestRedisDeleteSession(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	sr := NewSessionRepositoryRedis(s.Addr(), "", 0)

	sessionBytes, _ := json.Marshal(existingSesson)
	s.Set(existingSessionID, string(sessionBytes))

	t.Run("delete not existing session", func(t *testing.T) {
		err := sr.DeleteSession("bad")
		assert.NotNil(t, err)
	})

	t.Run("delete existing session", func(t *testing.T) {
		sess, err := s.Get(existingSessionID)
		assert.NotEmpty(t, sess)

		err = sr.DeleteSession(existingSessionID)
		assert.Nil(t, err)
		_, err = s.Get(existingSessionID)
		assert.NotNil(t, err)
	})
}

func TestRedisCreateSession(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	sr := NewSessionRepositoryRedis(s.Addr(), "", 0)

	t.Run("create session", func(t *testing.T) {

		newRedisSesson := Session{
			UserID:     "user1",
			Realm:      "users",
			Expired:    (time.Now().UnixNano() / int64(time.Millisecond)) + 60*60*24,
			Properties: map[string]string{"foo": "bar", "abc": "bac"},
		}

		created, err := sr.CreateSession(newRedisSesson)
		createdBytes, _ := json.Marshal(created)
		assert.Nil(t, err)

		existing, err := s.Get(created.ID)
		assert.Nil(t, err)
		assert.Equal(t, string(createdBytes), existing)
	})
}
