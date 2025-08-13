package routes

import (
	"desa-kepayang-backend/controllers"
	"desa-kepayang-backend/middleware"

	"github.com/gin-gonic/gin"
)

func InfoDesaRoutes(r *gin.Engine) {
	route := r.Group("/info-desa")
	{
		route.GET("/", controllers.GetAllInfoDesa)
		route.GET("/:id", controllers.GetInfoDesaByID)
		route.POST("/", middleware.AuthMiddleware(), controllers.CreateInfoDesa)
		route.PUT("/:id", middleware.AuthMiddleware(), controllers.UpdateInfoDesa)
		route.DELETE("/:id", middleware.AuthMiddleware(), controllers.DeleteInfoDesa)
	}
}
