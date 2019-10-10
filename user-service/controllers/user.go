package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/maximthomas/blazewall/user-service/config"

	"github.com/maximthomas/blazewall/user-service/models"

	"github.com/gin-gonic/gin"
)

func GetUserController() UserController {
	return UserController{
		uc: config.GetUserServiceConfig(),
	}
}

func GetUserControllerByUserServiceConfig(us config.UserServiceConfig) UserController {
	return UserController{
		uc: us,
	}
}

type UserController struct {
	uc config.UserServiceConfig
}

func getGinErrorJSON(err error) interface{} {
	ginErr := gin.Error{
		Err: err,
	}
	return ginErr.JSON()
}

func (us *UserController) GetUser(c *gin.Context) {
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

func (us *UserController) CreateUser(c *gin.Context) {
	var user models.User
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

func (us *UserController) UpdateUser(c *gin.Context) {
	var user models.User
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

func (us *UserController) DeleteUser(c *gin.Context) {
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

func (us *UserController) SetPassword(c *gin.Context) {
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

	var pass models.Password
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

func (us *UserController) ValidatePassword(c *gin.Context) {
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

	var pass models.Password
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

	valPassrodRes := models.ValidatePasswordResult{passwordRes}

	c.JSON(http.StatusOK, valPassrodRes)
}
