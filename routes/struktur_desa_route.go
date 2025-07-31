package routes

import (
	"desa-kepayang-backend/controllers"
	"desa-kepayang-backend/middleware"

	"github.com/gin-gonic/gin"
)

func StrukturDesaRoutes(r *gin.Engine) {
	struktur := r.Group("/struktur-desa")
	{
		struktur.GET("/", controllers.GetAllStrukturDesa)
		struktur.POST("/", middleware.AuthMiddleware(), controllers.CreateStrukturDesa)
		struktur.PUT("/:id", middleware.AuthMiddleware(), controllers.UpdateStrukturDesa)
		struktur.DELETE("/:id", middleware.AuthMiddleware(), controllers.DeleteStrukturDesa)
	}
}
