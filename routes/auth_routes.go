package routes

import (
	"sis-api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(api *gin.RouterGroup) {
	auth := api.Group("/auth")
	{
		auth.POST("/register", controllers.RegisterUser)
		auth.POST("/login", controllers.LoginUser)
	}
}
