package repo

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/maximthomas/blazewall/auth-service/pkg/models"
)

type UserRepository interface {
	GetUser(id string) (models.User, bool)
	ValidatePassword(id, password string) bool
	CreateUser(user models.User) (models.User, error)
	UpdateUser(user models.User) error
	SetPassword(id, password string) error
}

type UserRestRepository struct {
	realm    string
	endpoint string
	client   http.Client
}

func (ur *UserRestRepository) GetUser(id string) (user models.User, exists bool) {

	resp, err := ur.client.Get(ur.endpoint + "/users/" + id)
	if err != nil {
		log.Printf("error getting user: %v", err)
		return user, exists
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		log.Printf("got bad response from user service: %v", resp)
		return user, exists
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error getting user: %v", err)
		return user, exists
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Printf("error unmarshalling user: %v", err)
		return user, exists
	}
	log.Printf("got user user: %v", user)
	exists = true
	return user, exists
}

func (ur *UserRestRepository) ValidatePassword(id, password string) (valid bool) {

	pr := models.Password{
		Password: password,
	}

	prBytes, err := json.Marshal(pr)
	if err != nil {
		return valid
	}

	buf := bytes.NewBuffer(prBytes)
	resp, err := ur.client.Post(ur.endpoint+"/users/"+id+"/validatepassword", "application/json", buf)
	if err != nil {
		log.Printf("error validating password: %v", err)
		return valid
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		log.Printf("got bad response from user service: %v", resp)
		return valid
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error validating password: %v", err)
		return valid
	}
	var vpr models.ValidatePasswordResult

	err = json.Unmarshal(body, &vpr)
	if err != nil {
		log.Printf("error validating password: %v", err)
		return valid
	}
	valid = vpr.Valid

	log.Printf("password validation result for user: %v %v", id, valid)

	return valid
}

func (ur *UserRestRepository) CreateUser(user models.User) (models.User, error) {
	return user, nil
}

func (ur *UserRestRepository) UpdateUser(user models.User) error {
	return nil
}
func (ur *UserRestRepository) SetPassword(id, password string) error {
	return nil
}

func NewUserRestRepository(realm, endpoint string) UserRestRepository {
	return UserRestRepository{
		realm:    realm,
		endpoint: endpoint,
	}
}

type InMemoryUserRepository struct {
	Users []models.User
	Realm string
}

func (ur InMemoryUserRepository) GetUser(id string) (user models.User, exists bool) {
	for _, u := range ur.Users {
		if u.ID == id {
			user = u
			exists = true
			break
		}
	}
	return user, exists
}

func (ur InMemoryUserRepository) ValidatePassword(id, password string) (valid bool) {
	if password == "password" {
		valid = true
	}
	return valid
}

func (ur InMemoryUserRepository) CreateUser(user models.User) (models.User, error) {
	return user, nil
}

func (ur InMemoryUserRepository) UpdateUser(user models.User) error {
	return nil
}
func (ur InMemoryUserRepository) SetPassword(id, password string) error {
	return nil
}

func NewInMemoryUserRepository() UserRepository {

	ds := InMemoryUserRepository{}
	ds.Users = []models.User{
		{
			ID:    "user1",
			Realm: "users",
			Roles: []string{"admin"},
		},
		{
			ID:    "user2",
			Realm: "users",
			Roles: []string{"manager"},
		},
		{
			ID:    "staff1",
			Realm: "staff",
			Roles: []string{"head_of_it"},
		},
	}
	return ds
}

func NewUserRepository() UserRepository {
	//ac := config.GetConfig()
	//sr = &RestSessionRepository{Endpoint: ac.Endpoints.SessionService}
	local := os.Getenv("DEV_LOCAL")
	if local == "true" {
		return NewInMemoryUserRepository()
	}
	return nil
}
