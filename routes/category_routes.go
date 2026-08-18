package routes

import (
	"sis-api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterCategoryRoutes(api *gin.RouterGroup, admin *gin.RouterGroup) {
	categories := api.Group("/categories")
	{
		categories.GET("", controllers.GetAllCategories)
		categories.GET("/:id", controllers.GetCategoryByID)
	}

	adminCategories := admin.Group("/categories")
	{
		adminCategories.POST("", controllers.CreateCategory)
		adminCategories.PUT("/:id", controllers.UpdateCategory)
		adminCategories.DELETE("/:id", controllers.DeleteCategory)
	}
}
