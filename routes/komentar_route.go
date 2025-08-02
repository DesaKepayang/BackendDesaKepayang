package routes

import (
	"desa-kepayang-backend/controllers"
	"desa-kepayang-backend/middleware"

	"github.com/gin-gonic/gin"
)

func KomentarRoutes(r *gin.Engine) {
	k := r.Group("/komentar")
	{
		k.GET("/", controllers.GetAllKomentar)
		k.GET("/:id", controllers.GetKomentarByID)
		k.POST("/", middleware.AuthMiddleware(), controllers.CreateKomentar)
		k.PUT("/:id", middleware.AuthMiddleware(), controllers.UpdateKomentar)
		k.DELETE("/:id", middleware.AuthMiddleware(), controllers.DeleteKomentar)
	}
}
