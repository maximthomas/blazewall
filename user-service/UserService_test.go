package main

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/assert"
)

var existingUser = User{
	ID:    "user1",
	Realm: "users",
}

const serviceURL = "http://localhost:8080/user-service/v1/users"

func TestGetUser(t *testing.T) {
	router := getRouter()
	t.Run("test getting existing user", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("GET", serviceURL+"/users/user1", nil)
		router.ServeHTTP(recorder, request)
		assert.Equal(t, 200, recorder.Result().StatusCode)

		var responseUser User
		err := json.Unmarshal([]byte(recorder.Body.String()), &responseUser)
		assert.NoError(t, err)
		assert.Equal(t, existingUser, responseUser)
	})

	t.Run("test getting not existing user", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("GET", serviceURL+"/users/bad", nil)
		router.ServeHTTP(recorder, request)
		assert.Equal(t, 404, recorder.Result().StatusCode)
		assert.Equal(t, `{"error":"User not found"}`, recorder.Body.String())
	})
}

func TestCreateUser(t *testing.T) {
	router := getRouter()
	t.Run("test create new user", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		bodyStr := `{"id": "user2", "realm": "users", "roles": ["admin", "manager"]}`
		body := bytes.NewBufferString(bodyStr)
		request := httptest.NewRequest("POST", serviceURL, body)
		router.ServeHTTP(recorder, request)

		assert.Equal(t, 200, recorder.Result().StatusCode)
		var responseUser User
		responseBody := recorder.Body.String()
		err := json.Unmarshal([]byte(responseBody), &responseUser)
		assert.NoError(t, err)

		wantUser := User{
			ID:    "user2",
			Realm: "users",
			Roles: []string{"admin", "manager"},
		}
		assert.Equal(t, wantUser, responseUser)

	})

	t.Run("try create existing user", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		bodyStr := `{"id": "user1", "realm": "users", "roles": ["admin", "manager"]}`
		body := bytes.NewBufferString(bodyStr)
		request := httptest.NewRequest("POST", serviceURL, body)
		router.ServeHTTP(recorder, request)
		assert.Equal(t, 400, recorder.Result().StatusCode)
		assert.Equal(t, `{"error":"User already exists"}`, recorder.Body.String())
	})

	t.Run("test create user for non existig realm", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		bodyStr := `{"id": "user1", "realm": "staff", "roles": ["admin", "manager"]}`
		body := bytes.NewBufferString(bodyStr)
		request := httptest.NewRequest("POST", serviceURL, body)
		router.ServeHTTP(recorder, request)
		assert.Equal(t, 404, recorder.Result().StatusCode)
		assert.Equal(t, `{"error":"Realm does not exist"}`, recorder.Body.String())
	})

}

func TestUpdateUser(t *testing.T) {
	us := getUserService()
	router := setupRouter(&us)
	t.Run("test update existing user", func(t *testing.T) {
		bodyStr := `{"id": "user1", "realm": "users", "roles": ["admin", "manager"]}`
		body := bytes.NewBufferString(bodyStr)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("PUT", serviceURL+"/users/user1", body)
		router.ServeHTTP(recorder, request)

		var responseUser User
		responseBody := recorder.Body.String()
		err := json.Unmarshal([]byte(responseBody), &responseUser)
		assert.NoError(t, err)
		assert.Equal(t, 200, recorder.Result().StatusCode)
		assertEqualJSON(t, bodyStr, responseBody)
	})

	t.Run("test update non existing user", func(t *testing.T) {
		bodyStr := `{"id": "user2", "realm": "users", "roles": ["admin", "manager"]}`
		body := bytes.NewBufferString(bodyStr)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("PUT", serviceURL+"/users/user2", body)
		router.ServeHTTP(recorder, request)

		assert.Equal(t, 400, recorder.Result().StatusCode)
		assert.Equal(t, `{"error":"User not found"}`, recorder.Body.String())

	})

	t.Run("test update wrong realm user", func(t *testing.T) {
		bodyStr := `{"id": "user1", "realm": "users", "roles": ["admin", "manager"]}`
		body := bytes.NewBufferString(bodyStr)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("PUT", serviceURL+"/staff/user1", body)
		router.ServeHTTP(recorder, request)

		assert.Equal(t, 400, recorder.Result().StatusCode)
		assert.Equal(t, `{"error":"User realm or ID does not match"}`, recorder.Body.String())

	})
}

func TestDeleteUser(t *testing.T) {
	router := getRouter()
	t.Run("test delete existing user", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("DELETE", serviceURL+"/users/user1", nil)
		router.ServeHTTP(recorder, request)
		assert.Equal(t, 202, recorder.Result().StatusCode)
	})

	t.Run("test delete non existing user", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("DELETE", serviceURL+"/staff/user2", nil)
		router.ServeHTTP(recorder, request)
		assert.Equal(t, 404, recorder.Result().StatusCode)
	})
}

func TestSetPasswordUser(t *testing.T) {
	router := getRouter()
	t.Run("test set password existing user", func(t *testing.T) {
		bodyStr := `{"password": "testPassword"}`
		body := bytes.NewBufferString(bodyStr)

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("POST", serviceURL+"/users/user1/setpassword", body)
		router.ServeHTTP(recorder, request)
		assert.Equal(t, 202, recorder.Result().StatusCode)
	})

	t.Run("test set password non existing user", func(t *testing.T) {
		bodyStr := `{"password": "testPassword"}`
		body := bytes.NewBufferString(bodyStr)

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("POST", serviceURL+"/users/user2/setpassword", body)
		router.ServeHTTP(recorder, request)
		assert.Equal(t, 404, recorder.Result().StatusCode)
	})
}

func TestSetValidatePaswordUser(t *testing.T) {
	router := getRouter()
	t.Run("test validate password for non existing user", func(t *testing.T) {
		bodyStr := `{"password": "testPassword"}`
		body := bytes.NewBufferString(bodyStr)

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("POST", serviceURL+"/users/user2/validatepassword", body)
		router.ServeHTTP(recorder, request)
		assert.Equal(t, 404, recorder.Result().StatusCode)
	})

	t.Run("test validate wrong password for existing user", func(t *testing.T) {
		bodyStr := `{"password": "testPassword"}`
		body := bytes.NewBufferString(bodyStr)

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("POST", serviceURL+"/users/user1/validatepassword", body)
		router.ServeHTTP(recorder, request)

		var vpr ValidatePasswordResult
		responseBody := recorder.Body.String()
		err := json.Unmarshal([]byte(responseBody), &vpr)
		assert.NoError(t, err)
		assert.Equal(t, 200, recorder.Result().StatusCode)
		assertEqualJSON(t, bodyStr, `{"result":false}`)

	})

	t.Run("test validate good password for existing user", func(t *testing.T) {
		bodyStr := `{"password": "password1"}`
		body := bytes.NewBufferString(bodyStr)

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("POST", serviceURL+"/users/user1/validatepassword", body)
		router.ServeHTTP(recorder, request)

		var vpr ValidatePasswordResult
		responseBody := recorder.Body.String()
		err := json.Unmarshal([]byte(responseBody), &vpr)
		assert.NoError(t, err)
		assert.Equal(t, 200, recorder.Result().StatusCode)
		assertEqualJSON(t, bodyStr, `{"result":true}`)

	})
}

func getRouter() *gin.Engine {
	us := getUserService()
	return setupRouter(&us)
}

func getUserService() UserService {
	return UserService{
		uc: UserServiceConfig{
			RealmRepos: map[string]UserRepository{"users": getInMemoryUserRepository()},
		},
	}
}

func getInMemoryUserRepository() UserRepository {
	return &InMemoryUserRepository{[]RepoUser{
		{
			User:     existingUser,
			Password: "password1",
		},
	}}
}
func assertEqualJSON(t *testing.T, expected, actual string) bool {
	t.Helper()
	var o1 interface{}
	var o2 interface{}

	var err error
	err = json.Unmarshal([]byte(expected), &o1)
	if err != nil {
		return assert.Fail(t, "error deserializing json")
	}
	err = json.Unmarshal([]byte(actual), &o2)
	if err != nil {
		return assert.Fail(t, "error deserializing json")
	}

	return reflect.DeepEqual(o1, o2)
}
