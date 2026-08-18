package main

import (
	"log"
	"sis-api/config"
	"sis-api/routes"

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

	routes.SetupRoutes(server)

	server.Run(":8080")
}
