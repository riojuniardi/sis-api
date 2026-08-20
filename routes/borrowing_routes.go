package routes

import (
	"sis-api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterBorrowingRoutes(api *gin.RouterGroup, admin *gin.RouterGroup) {
	borrowings := api.Group("/borrowings")
	{
		borrowings.GET("", controllers.GetAllBorrowings)
		borrowings.GET("/:id", controllers.GetBorrowingByID)
		borrowings.POST("", controllers.CreateBorrowing)
	}

	adminBorrowings := admin.Group("/borrowings")
	{
		adminBorrowings.PUT("/:id", controllers.UpdateBorrowing)
		adminBorrowings.DELETE("/:id", controllers.DeleteBorrowing)
		adminBorrowings.PATCH("/:id/approve", controllers.ApproveBorrowing)
		adminBorrowings.PATCH("/:id/return", controllers.ReturnBorrowing)
	}
}
