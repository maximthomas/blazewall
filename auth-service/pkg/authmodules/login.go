package authmodules

import (
	"github.com/gin-gonic/gin"
	"github.com/maximthomas/blazewall/auth-service/pkg/models"
)

type LoginModule struct {
	properties map[string]string
	callbacks []models.Callback
}

func (lm *LoginModule) Init(c *gin.Context) error {
	return nil
}
func (lm *LoginModule) Process(s *LoginSessionState, c *gin.Context) (ms ModuleState, cbs []models.Callback, err error) {
	return InProgress, lm.callbacks, err
}

func (lm *LoginModule) ProcessCallbacks(inCbs []models.Callback, s *LoginSessionState, c *gin.Context) (ms ModuleState, cbs []models.Callback, err error) {
	if c.Request.Method == "POST" {
		cbr := models.CallbackRequest{}
		err := c.ShouldBindJSON(&cbr)
		if err != nil {
			return ms, cbs, err
		}
		s.SharedState["test"] = "test"
	}
	return Pass, cbs, err
}

func NewLoginModule(properties map[string]string) *LoginModule {
	return &LoginModule{
		properties: properties,
		callbacks: []models.Callback{
			{
				Name:  "login",
				Type:  "text",
				Value: "",
			},
			{
				Name:  "password",
				Type:  "password",
				Value: "",
			},
		},
	}
}
