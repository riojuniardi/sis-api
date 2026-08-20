package routes

import (
	"sis-api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterItemRoutes(api *gin.RouterGroup, admin *gin.RouterGroup) {
	items := api.Group("/items")
	{
		items.GET("", controllers.GetAllItems)
		items.GET("/:id", controllers.GetItemByID)
	}

	adminItems := admin.Group("/items")
	{
		adminItems.POST("", controllers.CreateItem)
		adminItems.PUT("/:id", controllers.UpdateItem)
		adminItems.DELETE("/:id", controllers.DeleteItem)
	}
}
