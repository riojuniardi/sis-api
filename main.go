package main

import (
	"log"
	"sis-api/config"
	"sis-api/controllers"
	"sis-api/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.ConnectDB()

	server := gin.Default()

	api := server.Group("/api")
	{

		api.POST("/auth/register", controllers.RegisterUser)
		api.POST("/auth/login", controllers.LoginUser)

		protected := api.Group("/")
		protected.Use(middlewares.RequiredAuth())
		{
			protected.GET("/users/me", controllers.GetCurrentUser)
		}

		admin := api.Group("/admin")
		admin.Use(middlewares.RequiredAdmin())
		{
			admin.GET("/users", controllers.GetAllUsers)
			admin.GET("/users/:id", controllers.GetUserByID)
		}

	}

	server.Run(":8081")
}
