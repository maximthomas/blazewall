package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type AuthProcessor interface {
	ProcessAuthentication(c *gin.Context)
}

func getAuthProcessor(authConfigEntry AuthConfigEntry) (AuthProcessor, error) {
	var ap AuthProcessor

	if authConfigEntry.Type == plainAuthType {
		ap, err := NewTextFileAuthProcessor(authConfigEntry.Parameters)
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
	Username string `json:"username"`
	Password string `json:"password"`
}

type TextFileAuthProcessor struct {
	credentials []Credential
}

func (t TextFileAuthProcessor) ProcessAuthentication(c *gin.Context) {

	switch c.Request.Method {
	case "GET":
		{
			c.JSON(200, Credential{})
		}
	case "POST":
		{
			var recievedCredential Credential
			err := c.ShouldBindJSON(&recievedCredential)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
				return
			}
			for _, credential := range t.credentials {
				if credential.Username == recievedCredential.Username &&
					credential.Password == recievedCredential.Password {

					c.SetCookie("session", "sessId", 0, "/", "", false, true)
					c.Redirect(http.StatusFound, c.Query("redirect"))
					c.Writer.WriteHeaderNow()
					return
				}
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}
	default:
		c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
	}
}
