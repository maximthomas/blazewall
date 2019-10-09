package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	Realm string `json:"realm,omitempty"`

	ID string `json:"id,omitempty"`

	Roles []string `json:"roles,omitempty"`

	Properties map[string]string `json:"properties,omitempty"`
}

type Password struct {
	Password string `json:"password,omitempty"`
}

type ValidatePasswordResult struct {
	Valid bool `json:"valid,omitempty"`
}

type UserService struct {
	uc UserServiceConfig
}

func getGinErrorJSON(err error) interface{} {
	ginErr := gin.Error{
		Err: err,
	}
	return ginErr.JSON()
}

func (us *UserService) GetUser(c *gin.Context) {
	realm := c.Param("realm")
	userID := c.Param("id")
	ur, ok := us.uc.RealmRepos[realm]
	if !ok {
		log.Printf("realm: %v does not exst", realm)
		c.AbortWithStatusJSON(http.StatusNotFound, getGinErrorJSON(errors.New("Realm does not exist")))
		return
	}
	user, err := ur.GetUser(realm, userID)
	if err != nil {
		log.Printf("error getting user: %v", err)
		c.AbortWithStatusJSON(http.StatusNotFound, getGinErrorJSON(err))
		return
	}
	c.JSON(http.StatusOK, user)
}

func (us *UserService) CreateUser(c *gin.Context) {
	var user User
	log.Printf("creating user")
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, getGinErrorJSON(err))
		return
	}

	ur, ok := us.uc.RealmRepos[user.Realm]
	if !ok {
		log.Printf("realm: %v does not exst", user.Realm)
		c.AbortWithStatusJSON(http.StatusNotFound, getGinErrorJSON(errors.New("Realm does not exist")))
		return
	}

	user, err = ur.CreateUser(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, getGinErrorJSON(err))
		return
	}
	c.JSON(http.StatusOK, user)
}

func (us *UserService) UpdateUser(c *gin.Context) {
	var user User
	log.Printf("creating user")
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, getGinErrorJSON(err))
		return
	}

	realm := c.Param("realm")
	userID := c.Param("id")
	if user.ID != userID || realm != user.Realm {
		c.AbortWithStatusJSON(http.StatusBadRequest, getGinErrorJSON(errors.New("User realm or ID does not match")))
		return
	}

	ur, ok := us.uc.RealmRepos[realm]
	if !ok {
		log.Printf("realm: %v does not exst", realm)
		c.AbortWithStatusJSON(http.StatusNotFound, getGinErrorJSON(errors.New("Realm does not exist")))
		return
	}

	updatedUser, err := ur.UpdateUser(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, getGinErrorJSON(err))
		return
	}
	c.JSON(http.StatusOK, updatedUser)
}

func (us *UserService) DeleteUser(c *gin.Context) {
	realm := c.Param("realm")
	userID := c.Param("id")

	ur, ok := us.uc.RealmRepos[realm]
	if !ok {
		log.Printf("realm: %v does not exst", realm)
		c.AbortWithStatusJSON(http.StatusNotFound, getGinErrorJSON(errors.New("Realm does not exist")))
		return
	}

	_, err := ur.GetUser(realm, userID)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, getGinErrorJSON(err))
		return
	}
	err = ur.DeleteUser(realm, userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, getGinErrorJSON(err))
		return
	}
	c.JSON(http.StatusAccepted, nil)
}

func (us *UserService) SetPassword(c *gin.Context) {
	realm := c.Param("realm")
	userID := c.Param("id")

	ur, ok := us.uc.RealmRepos[realm]
	if !ok {
		log.Printf("realm: %v does not exst", realm)
		c.AbortWithStatusJSON(http.StatusNotFound, getGinErrorJSON(errors.New("Realm does not exist")))
		return
	}

	_, err := ur.GetUser(realm, userID)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, getGinErrorJSON(err))
		return
	}

	var pass Password
	err = c.ShouldBindJSON(&pass)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, getGinErrorJSON(err))
		return
	}

	err = ur.SetPassword(realm, userID, pass.Password)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, getGinErrorJSON(err))
		return
	}

	c.JSON(http.StatusAccepted, nil)
}

func (us *UserService) ValidatePassword(c *gin.Context) {
	realm := c.Param("realm")
	userID := c.Param("id")

	ur, ok := us.uc.RealmRepos[realm]
	if !ok {
		log.Printf("realm: %v does not exst", realm)
		c.AbortWithStatusJSON(http.StatusNotFound, getGinErrorJSON(errors.New("Realm does not exist")))
		return
	}

	_, err := ur.GetUser(realm, userID)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, getGinErrorJSON(err))
		return
	}

	var pass Password
	err = c.ShouldBindJSON(&pass)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, getGinErrorJSON(err))
		return
	}

	passwordRes, err := ur.ValidatePassword(realm, userID, pass.Password)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, getGinErrorJSON(err))
		return
	}

	valPassrodRes := ValidatePasswordResult{passwordRes}

	c.JSON(http.StatusOK, valPassrodRes)
}
