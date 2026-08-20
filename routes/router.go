package routes

import (
	"sis-api/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(server *gin.Engine) {
	api := server.Group("/api")

	admin := api.Group("/admin")
	admin.Use(middlewares.RequiredAdmin())

	RegisterAuthRoutes(api)
	RegisterUserRoutes(api, admin)
	RegisterCategoryRoutes(api, admin)
	RegisterLocationRoutes(api, admin)
	RegisterItemRoutes(api, admin)
}
