package authmodules

import (
	"github.com/gin-gonic/gin"
	"github.com/maximthomas/blazewall/auth-service/pkg/auth"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegistration_Process(t *testing.T) {
	b := BaseAuthModule{
		properties: map[string]interface{}{
			keyAdditionalFileds: []AdditionalFiled{{
				dataStore: "name",
				prompt:    "Name",
			},
			},
		},
	}
	rm := NewRegistrationModule(b)
	t.Run("Test request callbacks", func(t *testing.T) {

		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		c.Request = httptest.NewRequest("GET", "/login", nil)
		lss := &auth.LoginSessionState{}
		status, cbs, err := rm.Process(lss, c)
		log.Print(status, cbs, err)
		assert.Equal(t, 3, len(cbs))
		assert.NoError(t, err)
		assert.Equal(t, auth.InProgress, status)
		assert.Equal(t, http.StatusOK, recorder.Code)
	})
}
