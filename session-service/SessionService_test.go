package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
)

var existingSesson = Session{
	ID:     "sess1",
	UserID: "user1",
	Realm:  "users",
}

func getNewSessionService() SessionService {

	repo := NewInMemorySessionRepository([]Session{
		existingSesson,
	})
	return NewSessionService(repo)
}

func TestSessionGetById(t *testing.T) {

	ss := getNewSessionService()
	t.Run("get bad session", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		c.Params = gin.Params{
			gin.Param{Key: "id", Value: "bad"},
		}

		ss.getSessionByID(c)

		assert.Equal(t, recorder.Result().StatusCode, 404)
	})

	t.Run("get good session", func(t *testing.T) {
		t.Run("get bad session", func(t *testing.T) {
			recorder := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(recorder)
			c.Params = gin.Params{
				gin.Param{Key: "id", Value: "sess1"},
			}

			ss.getSessionByID(c)

			assert.Equal(t, recorder.Result().StatusCode, 200)

			var responseSession Session
			err := json.Unmarshal([]byte(recorder.Body.String()), &responseSession)
			assert.NoError(t, err)
			assert.Equal(t, responseSession, existingSesson)
		})
	})
}

func TestSessionsFind(t *testing.T) {
	ss := getNewSessionService()
	t.Run("try to execute bad request", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		c.Request, _ = http.NewRequest("GET", "/?userId=1", nil)
		ss.findSessions(c)
		assert.Equal(t, recorder.Result().StatusCode, 400)
	})

	t.Run("try to execute valid request", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		c.Request, _ = http.NewRequest("GET", "/?userId=user1&realm=users", nil)
		ss.findSessions(c)
		assert.Equal(t, recorder.Result().StatusCode, 200)

		var responseSessions []Session
		err := json.Unmarshal([]byte(recorder.Body.String()), &responseSessions)
		assert.NoError(t, err)
		assert.Equal(t, len(responseSessions), 1)
		assert.Equal(t, responseSessions[0], existingSesson)
	})

	t.Run("try to search not exiting users sessions", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		c.Request, _ = http.NewRequest("GET", "/?userId=bad_user&realm=users", nil)
		ss.findSessions(c)
		assert.Equal(t, recorder.Result().StatusCode, 404)
	})
}

func TestDeleteSession(t *testing.T) {
	ss := getNewSessionService()
	ss.sr.CreateSession(Session{
		ID:     "session2",
		UserID: "user2",
	})
	t.Run("try to delete not existing session", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		c.Params = gin.Params{
			gin.Param{Key: "id", Value: "bad"},
		}
		ss.deleteSession(c)
		assert.Equal(t, recorder.Result().StatusCode, 404)
	})

	t.Run("try to delete existing session", func(t *testing.T) {
		sr := ss.sr.(*InMemorySessionRepository)
		assert.Equal(t, len(*sr), 2)
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		c.Params = gin.Params{
			gin.Param{Key: "id", Value: "session2"},
		}
		ss.deleteSession(c)
		assert.Equal(t, recorder.Result().StatusCode, 202)
		assert.Equal(t, len(*sr), 1)
	})
}

func TestCreateSession(t *testing.T) {

	ss := getNewSessionService()

	t.Run("try to create session from bad data", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		body := bytes.NewBufferString("asdasd")
		c.Request, _ = http.NewRequest("POST", "/", body)

		ss.createSession(c)
		assert.Equal(t, recorder.Result().StatusCode, 500)
	})

	t.Run("try to create request from good data", func(t *testing.T) {

		newSession := Session{
			UserID: "user2",
			Realm:  "users",
		}
		bodyBytes, err := json.Marshal(newSession)
		assert.NoError(t, err)
		bodyStr := string(bodyBytes)
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		body := bytes.NewBufferString(bodyStr)
		c.Request, _ = http.NewRequest("POST", "/", body)

		ss.createSession(c)
		assert.Equal(t, recorder.Result().StatusCode, 200)

		var createdSession Session
		responseBody := recorder.Body.String()
		unmarsahErr := json.Unmarshal([]byte(responseBody), &createdSession)
		assert.NoError(t, unmarsahErr)
		assert.Equal(t, createdSession.UserID, newSession.UserID)
		assertNotEmpty(t, createdSession.ID)
	})
}
