package api

import (
	"comparei-servico-usuario/internal/controller"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, userController *controller.UserController) {
	userRoutes := router.Group("/user")
	{
		userRoutes.GET("/:id", userController.GetUserByID)
	}
}
