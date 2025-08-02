package routes

import (
	"desa-kepayang-backend/controllers"
	"desa-kepayang-backend/middleware"

	"github.com/gin-gonic/gin"
)

func JumlahKKRoutes(r *gin.Engine) {
	jk := r.Group("/jumlahkk")
	{
		jk.GET("/", controllers.GetAllJumlahKK)
		jk.GET("/:id", controllers.GetJumlahKKByID)
		jk.POST("/", middleware.AuthMiddleware(), controllers.CreateJumlahKK)
		jk.PUT("/:id", middleware.AuthMiddleware(), controllers.UpdateJumlahKK)
		jk.DELETE("/:id", middleware.AuthMiddleware(), controllers.DeleteJumlahKK)
	}
}
