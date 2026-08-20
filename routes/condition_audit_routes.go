package routes

import (
	"sis-api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterConditionAuditRoutes(api *gin.RouterGroup, admin *gin.RouterGroup) {
	// Endpoint publik/terautentikasi umum
	audits := api.Group("/condition-audits")
	{
		audits.GET("", controllers.GetAllConditionAudits)
		audits.GET("/:id", controllers.GetConditionAuditByID)
		audits.GET("/item/:item_id", controllers.GetAuditsByItemID)
	}

	// Endpoint khusus Petugas/Admin
	adminAudits := admin.Group("/condition-audits")
	{
		adminAudits.POST("", controllers.CreateConditionAudit)
		adminAudits.DELETE("/:id", controllers.DeleteConditionAudit)
	}
}
