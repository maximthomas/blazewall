package repo

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/maximthomas/blazewall/auth-service/models"
)

type UserService interface {
	GetUser(realm, id string) (models.User, bool)
	ValidatePassword(realm, id, password string) bool
}

type UserRestService struct {
	realm    string
	endpoint string
	client   http.Client
}

func (us *UserRestService) GetUser(realm, id string) (user models.User, exists bool) {

	resp, err := us.client.Get(us.endpoint + "/" + realm + "/" + id)
	if err != nil {
		log.Fatalf("error getting session: %v", err)
		return user, exists
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error getting session: %v", err)
		return user, exists
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Fatalf("error unmarshalling session: %v", err)
		return user, exists
	}
	exists = true
	return user, exists
}

func (us *UserRestService) ValidatePassword(realm, id, password string) (valid bool) {

	pr := models.Password{
		Password: password,
	}

	prBytes, err := json.Marshal(pr)
	if err != nil {
		return valid
	}

	buf := bytes.NewBuffer(prBytes)
	resp, err := us.client.Post(us.endpoint+"/"+realm+"/"+id, "application/json", buf)
	if err != nil {
		log.Fatal(err)
		return valid
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return valid
	}
	var vpr models.ValidatePasswordResult

	err = json.Unmarshal(body, &vpr)
	if err != nil {
		log.Fatal(err)
		return valid
	}
	valid = vpr.Valid
	return valid
}

func GetUserRestService(realm, endpoint string) UserRestService {
	return UserRestService{
		realm:    realm,
		endpoint: endpoint,
	}
}

type DummyUserService struct {
	Users []models.User
	Realm string
}

func (us DummyUserService) GetUser(realm, id string) (user models.User, exists bool) {
	for _, u := range us.Users {
		if u.Realm == realm && u.ID == id {
			user = u
			exists = true
			break
		}
	}
	return user, exists
}

func (us DummyUserService) ValidatePassword(realm, id, password string) (valid bool) {
	if password == "pass" {
		valid = true
	}
	return valid
}

func NewDummyUserService() UserService {

	ds := DummyUserService{}
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
