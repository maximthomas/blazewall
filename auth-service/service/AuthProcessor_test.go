package service

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTextFileAuth(t *testing.T) {
	ap := TextFileAuthProcessor{
		credentials: []Credential{
			{
				Username: "admin",
				Password: "pass",
			},
		},
	}

	t.Run("test initialise authentication", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		c.Request, _ = http.NewRequest("GET", "/users/", nil)

		ap.ProcessAuthentication(c)

		assert.Equal(t, recorder.Code, 200)

		responseBody := recorder.Body.String()

		want := `{"username":"","password":""}`
		assert.Equal(t, want, responseBody)

	})

	t.Run("test process bad requesst", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		body := bytes.NewBufferString(`bad request`)
		c.Request, _ = http.NewRequest("POST", "/users/", body)
		ap.ProcessAuthentication(c)

		assert.Equal(t, recorder.Code, 400)
		want := `{"error":"Bad request"}`
		responseBody := recorder.Body.String()
		assert.Equal(t, want, responseBody)

	})

	t.Run("test process invalid username, password", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		body := bytes.NewBufferString(`{"username":"bad","password":"bad"}`)
		c.Request, _ = http.NewRequest("POST", "/users/", body)
		ap.ProcessAuthentication(c)

		assert.Equal(t, recorder.Code, 401)
		want := `{"error":"Invalid username or password"}`
		responseBody := recorder.Body.String()
		assert.Equal(t, want, responseBody)

	})

	t.Run("test process valid username, password", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		body := bytes.NewBufferString(`{"username":"admin","password":"pass"}`)
		c.Request, _ = http.NewRequest("POST", "/users/?redirect=http://protected-resource:8080", body)
		sess, ok := ap.ProcessAuthentication(c)
		assert.True(t, ok)
		assert.Equal(t, "admin", sess.UserID)
	})
}

func TestNewTextFileAuthProcessor(t *testing.T) {
	var filePath interface{}
	filePath = "../test/users.txt"
	params := map[string]*interface{}{
		"filePath": &filePath,
	}

	ap, err := NewTextFileAuthProcessor(params)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(ap.credentials))
}
