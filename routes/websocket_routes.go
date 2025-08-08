package routes

import (
	"desa-kepayang-backend/controllers"

	"github.com/gin-gonic/gin"
)

func WebSocketRoutes(r *gin.Engine) {
	r.GET("/ws", controllers.WebSocketHandler)
}
