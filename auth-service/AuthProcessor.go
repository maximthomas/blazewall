package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type AuthProcessor interface {
	ProcessAuthentication(c *gin.Context) (Session, bool)
}

const plainAuthType = "textFile"

const userServiceAuthType = "userService"

func getAuthProcessor(authConfigEntry AuthConfigEntry) (AuthProcessor, error) {
	var ap AuthProcessor

	if authConfigEntry.Type == plainAuthType {
		ap, err := NewTextFileAuthProcessor(authConfigEntry.Parameters)
		return ap, err
	}
	if authConfigEntry.Type == userServiceAuthType {
		ap, err := NewUserServiceAuthProcessor(authConfigEntry.Parameters)
		return ap, err
	}

	return ap, nil
}

func NewTextFileAuthProcessor(parameters map[string]*interface{}) (TextFileAuthProcessor, error) {
	var ap TextFileAuthProcessor
	filePathPtr, ok := parameters["filePath"]
	if !ok {
		return ap, fmt.Errorf("Text Auth Entry file path undefinded %v", parameters)
	}

	filePath := fmt.Sprintf("%v", *filePathPtr)

	configReader, err := os.Open(filePath)

	if err != nil {
		log.Fatal(err)
		return ap, err
	}

	r := csv.NewReader(configReader)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
			return ap, err
		}
		if len(record) != 2 {
			log.Printf("bad record, skipping %v", record)
			continue
		}
		ap.credentials = append(ap.credentials, Credential{
			Username: record[0],
			Password: record[1],
		})

	}

	return ap, nil
}

type Credential struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type TextFileAuthProcessor struct {
	credentials []Credential
}

func (t TextFileAuthProcessor) ProcessAuthentication(c *gin.Context) (sess Session, ok bool) {

	ok = false
	switch c.Request.Method {
	case "GET":
		{
			c.JSON(200, Credential{})
			return sess, ok
		}
	case "POST":
		{
			var recievedCredential Credential
			err := c.ShouldBindJSON(&recievedCredential)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
				return sess, ok
			}
			for _, credential := range t.credentials {
				if credential.Username == recievedCredential.Username &&
					credential.Password == recievedCredential.Password {
					sess.UserID = recievedCredential.Username
					sess.Properties = make(map[string]string)
					sess.Properties["role"] = "users"
					ok = true
					return sess, ok
				}
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return sess, ok
		}
	}
	c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
	return sess, ok

}

type UserServiceAuthProcessor struct {
	us        UserService
	realmName string
}

func (ap UserServiceAuthProcessor) ProcessAuthentication(c *gin.Context) (sess Session, ok bool) {
	ok = false
	switch c.Request.Method {
	case "GET":
		{
			c.HTML(200, "username-password.html", nil)
			return sess, ok
		}
	case "POST":
		{
			var recievedCredential Credential
			err := c.ShouldBind(&recievedCredential)
			if err != nil {
				c.HTML(http.StatusUnauthorized, "username-password.html", gin.H{"error": "Invalid username or password"})
				return sess, ok
			}

			//check user exists
			user, exists := ap.us.GetUser(ap.realmName, recievedCredential.Username)
			if !exists {
				c.HTML(http.StatusUnauthorized, "username-password.html", gin.H{"error": "Invalid username or password"})
				return sess, ok
			}

			valid := ap.us.ValidatePassword(ap.realmName, user.ID, recievedCredential.Password)
			if !valid {
				c.HTML(http.StatusUnauthorized, "username-password.html", gin.H{"error": "Invalid username or password"})
				return sess, ok
			}

			sess := Session{
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
		return ap, fmt.Errorf("Text Auth Entry file path undefinded %v", parameters)
	}
	ap.realmName = fmt.Sprintf("%v", *realmPtr)
	ap.us = NewDummyUserService()
	return ap, nil
}
