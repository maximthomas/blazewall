package server

import (
	"github.com/gin-gonic/gin"
	"github.com/maximthomas/blazewall/auth-service/pkg/config"
	"github.com/maximthomas/blazewall/auth-service/pkg/repo"
	"github.com/maximthomas/blazewall/auth-service/pkg/server/controller"
)

func setupRouter(auth config.Authentication) *gin.Engine {
	router := gin.Default()

	var loginController = controller.NewLoginController(auth, repo.NewInMemorySessionRepository())

	v1 := router.Group("/auth-service/v1")
	{
		login := v1.Group("/login")
		{
			route := "/:realm/:chain"
			login.GET(route, func(context *gin.Context) {
				realmId := context.Param("realm")
				authChainId := context.Param("chain")
				loginController.Login(realmId, authChainId, context)
			})
			login.POST(route, func(context *gin.Context) {
				realmId := context.Param("realm")
				authChainId := context.Param("chain")
				loginController.Login(realmId, authChainId, context)
			})
		}
	}
	return router
}

func RunServer() {
	ac := config.GetConfig()
	router := setupRouter(ac)
	router.Run(":" + "8080")
}
