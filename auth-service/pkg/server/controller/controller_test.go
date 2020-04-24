package controller

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/maximthomas/blazewall/auth-service/pkg/auth"
	"github.com/maximthomas/blazewall/auth-service/pkg/config"
	"github.com/maximthomas/blazewall/auth-service/pkg/repo"
	"github.com/prometheus/common/log"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	ac = config.Authentication{
		Realms: map[string]config.Realm{
			"staff": {
				Modules: map[string]config.Module{
					"login": {Type: "login"},
				},
				AuthChains: map[string]config.AuthChain{
					"default": {Modules: []config.ChainModule{
						{
							ID: "login",
						},
					}},
					"sso": {Modules: []config.ChainModule{}},
				},
				UserRepo: repo.NewInMemoryUserRepository(),
			},
		},
	}
	lc = NewLoginController(ac, repo.NewInMemorySessionRepository())
)

func TestControllerLoginByRealmChain(t *testing.T) {
	var tests = []struct {
		expectedStatus int
		realmId        string
		authChainId    string
	}{
		{404, "clients", "users"},
		{404, "staff", "users"},
		{200, "staff", "default"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprint(tt), func(t *testing.T) {
			recorder := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(recorder)
			c.Request = httptest.NewRequest("GET", "/login", nil)

			lc.Login(tt.realmId, tt.authChainId, c)
			assert.Equal(t, tt.expectedStatus, recorder.Result().StatusCode)
			log.Info(recorder.Body.String())
		})
	}
}

func TestLoginPassword(t *testing.T) {
	t.Run("Test auth password", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		bodyStr := `{"callbacks":
[{"name":"login","type":"text","value":"user1"},{"name":"password","type":"password","value":"pass"}]}`
		body := bytes.NewBufferString(bodyStr)
		c.Request = httptest.NewRequest("POST", "/login", body)

		lc.Login("staff", "default", c)
		assert.Equal(t, 200, recorder.Result().StatusCode)
		log.Info(recorder.Body.String())
	})
}

func TestGetSessionState(t *testing.T) {
	t.Run("Test Get New Session State", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		c.Request = httptest.NewRequest("GET", "/login", nil)
		lls := lc.getLoginSessionState(ac.Realms["staff"].AuthChains["default"], ac.Realms["staff"], c)
		assert.Equal(t, 1, len(lls.Modules))
	})

	t.Run("Test Get Existing State", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		c.Request = httptest.NewRequest("GET", "/login", nil)
		lss := lc.getLoginSessionState(ac.Realms["staff"].AuthChains["default"], ac.Realms["staff"], c)
		assert.Equal(t, 1, len(lss.Modules))
		lss.SharedState["key"] = "value"
		lss.UserId = "user1"
		err := lc.updateLoginSessionState(lss)
		assert.Nil(t, err)

		cSecond, _ := gin.CreateTestContext(recorder)
		authCookie := &http.Cookie{
			Name:  auth.AuthCookieName,
			Value: lss.SessionId,
		}
		cSecond.Request = httptest.NewRequest("GET", "/login", nil)
		cSecond.Request.AddCookie(authCookie)
		lssUpdated := lc.getLoginSessionState(ac.Realms["staff"].AuthChains["default"], ac.Realms["staff"], cSecond)
		assert.Equal(t, "value", lssUpdated.SharedState["key"])
		assert.Equal(t, "user1", lssUpdated.UserId)

	})
}
