package routes

import (
	"sis-api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterLocationRoutes(api *gin.RouterGroup, admin *gin.RouterGroup) {
	location := api.Group("/locations")
	{
		location.GET("", controllers.GetLocations)
		location.GET("/:id", controllers.GetLocationByID)
	}

	adminLocations := admin.Group("/locations")
	{
		adminLocations.POST("", controllers.CreateLocation)
		adminLocations.PUT("/:id", controllers.UpdateLocation)
		adminLocations.DELETE("/:id", controllers.DeleteLocation)
	}
}
