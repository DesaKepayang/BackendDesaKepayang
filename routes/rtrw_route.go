package routes

import (
	"desa-kepayang-backend/controllers"
	"desa-kepayang-backend/middleware"

	"github.com/gin-gonic/gin"
)

func RTRWRoutes(r *gin.Engine) {
	rt := r.Group("/rtrw")
	{
		rt.GET("/", controllers.GetAllRTRW)
		rt.POST("/", middleware.AuthMiddleware(), controllers.CreateRTRW)
		rt.PUT("/:id", middleware.AuthMiddleware(), controllers.UpdateRTRW)
		rt.DELETE("/:id", middleware.AuthMiddleware(), controllers.DeleteRTRW)
	}
}
