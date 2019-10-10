package repo

import (
	"context"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/maximthomas/blazewall/user-service/models"
	"github.com/stretchr/testify/assert"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func startDockerCompose() {
	dir, err := os.Getwd()
	checkErr(err)
	cmd := exec.Command("docker-compose", "-f", "docker-compose-mongodb.yaml", "--project-directory", dir, "up", "-d")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	checkErr(err)
}

func stopDockerCompose() {
	dir, err := os.Getwd()
	checkErr(err)
	cmd := exec.Command("docker-compose", "-f", "docker-compose-mongodb.yaml", "--project-directory", dir, "down")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	checkErr(err)
}

func TestUserRepositoryMongoDBCreateUser(t *testing.T) {

	repo := getRepo()
	t.Run("test create unexisting user", func(t *testing.T) {
		user := models.User{
			Realm:      "users",
			Roles:      []string{"admin"},
			Properties: map[string]string{"foo": "bar"},
		}
		newUser, err := repo.CreateUser(user)
		assert.Nil(t, err)
		assert.NotEmpty(t, newUser.ID)
		assert.Equal(t, user.Roles, newUser.Roles)
		assert.Equal(t, user.Properties, newUser.Properties)
	})

	t.Run("test create existing user", func(t *testing.T) {
		//panic("not implemented")
	})
}

func TestUserRepositoryMongoDBGetUser(t *testing.T) {

	repo := getRepo()
	t.Run("test getting unexisting user", func(t *testing.T) {
		_, err := repo.GetUser("users", "bad")
		assert.NotNil(t, err)
	})

	t.Run("test getting existing user", func(t *testing.T) {

		user := models.User{
			Realm:      "users",
			Roles:      []string{"admin"},
			Properties: map[string]string{"foo": "bar"},
		}
		newUser, err := repo.CreateUser(user)
		assert.Nil(t, err)

		existingUser, err := repo.GetUser("users", newUser.ID)
		assert.Nil(t, err)
		assert.Equal(t, newUser, existingUser)
	})
}

func TestUserRepositoryMongoDBUpdateUser(t *testing.T) {

	repo := getRepo()
	t.Run("test getting unexisting user", func(t *testing.T) {
		user := models.User{
			ID:         "bad",
			Realm:      "users",
			Roles:      []string{"admin"},
			Properties: map[string]string{"foo": "bar"},
		}
		_, err := repo.UpdateUser(user)
		assert.NotNil(t, err)
	})

	t.Run("test update existing user", func(t *testing.T) {

		user := models.User{
			Realm:      "users",
			Roles:      []string{"admin"},
			Properties: map[string]string{"foo": "bar"},
		}
		newUser, err := repo.CreateUser(user)
		assert.Nil(t, err)

		updateUser := newUser
		updateUser.Roles = []string{"admin, manager"}
		updateUser.Properties = map[string]string{"aaa": "bbb"}

		updatedUser, err := repo.UpdateUser(updateUser)
		assert.Nil(t, err)
		assert.Equal(t, updateUser, updatedUser)

	})
}

func TestUserRepositoryMongoDBDeleteUser(t *testing.T) {

	repo := getRepo()
	t.Run("test deleting unexisting user", func(t *testing.T) {
		err := repo.DeleteUser("users", "bad")
		assert.NotNil(t, err)
	})

	t.Run("test deleting existing user", func(t *testing.T) {
		user := models.User{
			Realm:      "users",
			Roles:      []string{"admin"},
			Properties: map[string]string{"foo": "bar"},
		}
		newUser, err := repo.CreateUser(user)
		assert.Nil(t, err)

		newUserID := newUser.ID
		newUserRealm := newUser.Realm
		err = repo.DeleteUser(newUser.Realm, newUserID)
		assert.Nil(t, err)

		_, err = repo.GetUser(newUserRealm, newUserID)

		assert.NotNil(t, err)
	})
}

func TestUserRepositoryMongoDBSetValidatePassword(t *testing.T) {
	repo := getRepo()

	user := models.User{
		Realm:      "users",
		Roles:      []string{"admin"},
		Properties: map[string]string{"foo": "bar"},
	}
	newUser, err := repo.CreateUser(user)
	assert.Nil(t, err)

	t.Run("test set password unexisting user", func(t *testing.T) {
		t.Run("set password unexisting user", func(t *testing.T) {
			err := repo.SetPassword("users", "bad", "newPassword")
			assert.NotNil(t, err)
		})
	})

	t.Run("test validate password unexisting user", func(t *testing.T) {
		_, err := repo.ValidatePassword("users", "bad", "newPassword")
		assert.NotNil(t, err)
	})

	t.Run("test set & validate password existing user", func(t *testing.T) {
		err := repo.SetPassword(newUser.Realm, newUser.ID, "newPassword")
		assert.Nil(t, err)

		ok, err := repo.ValidatePassword(newUser.Realm, newUser.ID, "newPassword")
		assert.Nil(t, err)
		assert.True(t, ok)

		ok, err = repo.ValidatePassword(newUser.Realm, newUser.ID, "bad")
		assert.Nil(t, err)
		assert.False(t, ok)
	})
}

func getRepo() UserRepositoryMongoDB {
	repo := NewUserRepositoryMongoDB("mongodb://root:example@localhost:27017", "test_users", "users")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	repo.client.Database(repo.db).Drop(ctx)
	return repo
}
