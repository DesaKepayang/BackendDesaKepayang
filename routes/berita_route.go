package routes

import (
	"desa-kepayang-backend/controllers"
	"desa-kepayang-backend/middleware"

	"github.com/gin-gonic/gin"
)

func BeritaRoutes(r *gin.Engine) {
	berita := r.Group("/berita")
	{
		berita.GET("/", controllers.GetAllBerita)
		berita.POST("/", middleware.AuthMiddleware(), controllers.CreateBerita)
		berita.PUT("/:id", middleware.AuthMiddleware(), controllers.UpdateBerita)
		berita.DELETE("/:id", middleware.AuthMiddleware(), controllers.DeleteBerita)
	}
}
