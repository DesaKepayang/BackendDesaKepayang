package routes

import (
	"desa-kepayang-backend/controllers"

	"github.com/gin-gonic/gin"
)

func VisiMisiRoutes(r *gin.Engine) {
	route := r.Group("/visimisi")
	{
		route.GET("/", controllers.GetAllVisiMisi)
		route.POST("/", controllers.CreateVisiMisi)
		route.PUT("/:id", controllers.UpdateVisiMisi)
		route.DELETE("/:id", controllers.DeleteVisiMisi)
	}
}
