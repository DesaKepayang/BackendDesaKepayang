package routes

import (
	"desa-kepayang-backend/controllers"
	"desa-kepayang-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SambutanRoutes(r *gin.Engine) {
	// Route publik (tanpa login)
	r.GET("/sambutan", controllers.GetSambutan)

	// Route yang butuh autentikasi
	sambutan := r.Group("/sambutan")
	sambutan.Use(middleware.AuthMiddleware())
	{
		sambutan.POST("/", controllers.TambahSambutan)
		sambutan.PATCH("/:id", controllers.UpdateSambutan)
		sambutan.DELETE("/:id", controllers.DeleteSambutan)
	}
}
