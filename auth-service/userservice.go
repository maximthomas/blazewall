package main

type User struct {
	ID         string            `json:"id,omitempty"`
	Realm      string            `json:"realm,omitempty"`
	Roles      []string          `json:"roles,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
}

type UserService interface {
	GetUser(realm, id string) (User, bool)
	ValidatePassword(realm, id, password string) bool
}

type UserRestService struct {
	Realm string
}

func (us UserRestService) GetUser(realm, id string) (user User, exists bool) {
	return user, exists
}

func (us UserRestService) ValidatePassword(realm, id, password string) (valid bool) {
	return valid
}

type DummyUserService struct {
	Users []User
	Realm string
}

func (us DummyUserService) GetUser(realm, id string) (user User, exists bool) {
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
	ds.Users = []User{
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
