package server

import (
	"github.com/gin-gonic/gin"
	"github.com/maximthomas/blazewall/user-service/controllers"
)

func SetupRouter() *gin.Engine {

	uc := controllers.GetUserController()
	return SetupRouterWithController(uc)

}

func SetupRouterWithController(uc controllers.UserController) *gin.Engine {
	router := gin.Default()
	v1 := router.Group("/user-service/v1")
	{
		users := v1.Group("/users")
		{
			users.GET("/:realm/:id", uc.GetUser)
			users.POST("", uc.CreateUser)
			users.PUT("/:realm/:id", uc.UpdateUser)
			users.DELETE("/:realm/:id", uc.DeleteUser)
			users.POST("/:realm/:id/setpassword", uc.SetPassword)
			users.POST("/:realm/:id/validatepassword", uc.ValidatePassword)
		}
	}
	return router
}
