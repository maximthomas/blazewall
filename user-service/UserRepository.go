package main

import (
	"errors"
)

type UserRepository interface {
	GetUser(realm, userID string) (User, error)
	CreateUser(user User) (User, error)
	UpdateUser(user User) (User, error)
	DeleteUser(realm, userID string) error
	SetPassword(realm, userID, password string) error
	ValidatePassword(realm, userID, password string) (bool, error)
}

type RepoUser struct {
	User     `bson:",inline"`
	Password string `json:"password,omitempty"`
}

type InMemoryUserRepository struct {
	repoUsers []RepoUser
}

var userNotFoudError = errors.New("User not found")

func (ur *InMemoryUserRepository) GetUser(realm, id string) (User, error) {
	var user User
	for _, ru := range ur.repoUsers {
		if ru.ID == id && ru.Realm == realm {
			return ru.User, nil
		}
	}
	return user, userNotFoudError
}

func (ur *InMemoryUserRepository) CreateUser(user User) (User, error) {
	var newUser User
	//check if user exists
	_, err := ur.GetUser(user.Realm, user.ID)
	if err == nil {
		return newUser, errors.New("User already exists")
	}
	//check if realm exists
	realmExists := false
	for _, existingUser := range ur.repoUsers {
		if existingUser.Realm == user.Realm {
			realmExists = true
			break
		}
	}
	if !realmExists {
		return newUser, errors.New("Realm does not exists")
	}

	ur.repoUsers = append(ur.repoUsers, RepoUser{
		User: user,
	})

	return user, nil
}

func (ur *InMemoryUserRepository) UpdateUser(user User) (User, error) {
	var updatedUser User

	for _, ru := range ur.repoUsers {
		if ru.ID == user.ID && ru.Realm == user.Realm {
			ru.Properties = user.Properties
			ru.Roles = user.Roles
			return ru.User, nil
		}
	}

	return updatedUser, userNotFoudError
}

func (ur *InMemoryUserRepository) DeleteUser(realm, userID string) error {
	idx := -1
	for i, ru := range ur.repoUsers {
		if ru.ID == userID && ru.Realm == realm {
			idx = i
		}
	}
	if idx == -1 {
		return userNotFoudError
	}

	ur.repoUsers = append(ur.repoUsers[:idx], ur.repoUsers[idx+1:]...)

	return nil
}

func (ur *InMemoryUserRepository) SetPassword(realm, userID, password string) error {

	for i, ru := range ur.repoUsers {
		if ru.ID == userID && ru.Realm == realm {
			ur.repoUsers[i].Password = password
			return nil
		}
	}
	return userNotFoudError
}

func (ur *InMemoryUserRepository) ValidatePassword(realm, userID, password string) (bool, error) {

	for _, ru := range ur.repoUsers {
		if ru.ID == userID && ru.Realm == realm {
			return ru.Password == password, nil
		}
	}
	return false, userNotFoudError
}
