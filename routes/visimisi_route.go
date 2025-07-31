package routes

import (
	"desa-kepayang-backend/controllers"
	"desa-kepayang-backend/middleware"

	"github.com/gin-gonic/gin"
)

func VisiMisiRoutes(r *gin.Engine) {
	// Route GET (publik)
	r.GET("/visimisi", controllers.GetAllVisiMisi)

	// Route yang dilindungi dengan middleware (POST, PUT, DELETE)
	protected := r.Group("/visimisi")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/", controllers.CreateVisiMisi)
		protected.PUT("/:id", controllers.UpdateVisiMisi)
		protected.DELETE("/:id", controllers.DeleteVisiMisi)
	}
}
