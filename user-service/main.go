package main

import (
	"flag"

	"github.com/gin-gonic/gin"
)

var port = flag.String("p", "8080", "User service port")

func setupRouter(us *UserService) *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/user-service/v1")
	{
		users := v1.Group("/users")
		{
			users.GET("/:realm/:id", us.GetUser)
			users.POST("", us.CreateUser)
			users.PUT("/:realm/:id", us.UpdateUser)
			users.DELETE("/:realm/:id", us.DeleteUser)
			users.POST("/:realm/:id/setpassword", us.SetPassword)
			users.POST("/:realm/:id/validatepassword", us.ValidatePassword)
		}
	}
	return router
}

func main() {

	flag.Parse()

	repo := InMemoryUserRepository{[]RepoUser{
		{
			User: User{
				ID:    "user1",
				Realm: "users",
			},
			Password: "pass",
		},
	}}
	repos := make(map[string]UserRepository)
	repos["users"] = &repo

	uc := UserServiceConfig{
		RealmRepos: repos,
	}

	ss := UserService{
		uc: uc,
	}
	router := setupRouter(&ss)
	router.Run(":" + *port)
}
