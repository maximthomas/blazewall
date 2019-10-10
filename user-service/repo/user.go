package repo

import (
	"github.com/maximthomas/blazewall/user-service/models"
)

type UserRepository interface {
	GetUser(realm, userID string) (models.User, error)
	CreateUser(user models.User) (models.User, error)
	UpdateUser(user models.User) (models.User, error)
	DeleteUser(realm, userID string) error
	SetPassword(realm, userID, password string) error
	ValidatePassword(realm, userID, password string) (bool, error)
}
