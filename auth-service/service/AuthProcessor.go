package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/maximthomas/blazewall/auth-service/repo"

	"github.com/maximthomas/blazewall/auth-service/models"

	"github.com/gin-gonic/gin"
	"github.com/maximthomas/blazewall/auth-service/config"
)

type AuthProcessor interface {
	ProcessAuthentication(c *gin.Context) (models.Session, bool)
}

const userServiceAuthType = "userService"

func GetAuthProcessor(authConfigEntry config.AuthConfigEntry) (AuthProcessor, error) {
	var ap AuthProcessor

	if authConfigEntry.Type == userServiceAuthType {
		ap, err := NewUserServiceAuthProcessor(authConfigEntry.Parameters)
		return ap, err
	}

	return ap, errors.New("udefuned auth processor")
}

type Credential struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type UserServiceAuthProcessor struct {
	us        repo.UserService
	realmName string
}

func appendCSRF(c *gin.Context, params gin.H) gin.H {
	csrf, ok := c.Get("csrfToken")
	if !ok {
		panic("token not present")
	}
	if params == nil {
		params = gin.H{}
	}
	params["csrfToken"] = csrf
	return params
}

func (ap UserServiceAuthProcessor) ProcessAuthentication(c *gin.Context) (sess models.Session, ok bool) {
	ok = false
	switch c.Request.Method {
	case "GET":
		{
			c.HTML(200, "username-password.html", appendCSRF(c, nil))
			return sess, ok
		}
	case "POST":
		{
			var recievedCredential Credential
			err := c.ShouldBind(&recievedCredential)
			if err != nil {
				c.HTML(http.StatusUnauthorized, "username-password.html", appendCSRF(c, gin.H{"error": "Invalid username or password"}))
				return sess, ok
			}

			//check user exists
			user, exists := ap.us.GetUser(ap.realmName, recievedCredential.Username)
			if !exists {
				c.HTML(http.StatusUnauthorized, "username-password.html", appendCSRF(c, gin.H{"error": "Invalid username or password"}))
				return sess, ok
			}

			valid := ap.us.ValidatePassword(ap.realmName, user.ID, recievedCredential.Password)
			if !valid {
				c.HTML(http.StatusUnauthorized, "username-password.html", appendCSRF(c, gin.H{"error": "Invalid username or password"}))
				return sess, ok
			}

			sess := models.Session{
				UserID: user.ID,
				Realm:  user.Realm,
			}
			sess.Properties = make(map[string]string)
			for key, value := range user.Properties {
				sess.Properties[key] = value
			}
			rolesJSON, _ := json.Marshal(user.Roles)
			sess.Properties["roles"] = string(rolesJSON)
			ok = true
			return sess, ok

		}
	}

	c.HTML(200, "username-password.html", nil)
	return sess, ok
}

func NewUserServiceAuthProcessor(parameters map[string]*interface{}) (UserServiceAuthProcessor, error) {
	var ap UserServiceAuthProcessor
	realmPtr, ok := parameters["realm"]
	if !ok {
		return ap, fmt.Errorf("realm udenfined undefinded %v", parameters)
	}
	endpointPtr, ok := parameters["endpoint"]
	if !ok {
		return ap, fmt.Errorf("endpoint undefinded %v", parameters)
	}
	ap.realmName = fmt.Sprintf("%v", *realmPtr)
	endpoint := fmt.Sprintf("%v", *endpointPtr)
	us := repo.GetUserRestService(ap.realmName, endpoint)
	ap.us = &us
	local := os.Getenv("DEV_LOCAL")
	if local == "true" {
		ap.us = repo.NewDummyUserService()
	}

	return ap, nil
}
