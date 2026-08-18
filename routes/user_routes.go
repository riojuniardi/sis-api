package routes

import (
	"sis-api/controllers"
	"sis-api/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(api *gin.RouterGroup, admin *gin.RouterGroup) {
	users := api.Group("/users")
	users.Use(middlewares.RequiredAuth())
	{
		users.GET("/me", controllers.GetCurrentUser)
	}

	adminUsers := admin.Group("/users")
	{
		adminUsers.GET("", controllers.GetAllUsers)
		adminUsers.GET("/:id", controllers.GetUserByID)
	}
}
