package routes

import (
	"desa-kepayang-backend/controllers"
	"desa-kepayang-backend/middleware"

	"github.com/gin-gonic/gin"
)

func PendudukRoutes(r *gin.Engine) {
	pd := r.Group("/penduduk")
	{
		pd.GET("/", controllers.GetAllPenduduk)
		pd.GET("/jumlah", controllers.CountPenduduk)
		pd.GET("/count/agama", controllers.CountPendudukByAgama)
		pd.GET("/count/perrtrw", controllers.CountPendudukPerRTRW)
		pd.GET("/jumlah-gender", controllers.CountPendudukByGender)
		pd.GET("/search", middleware.AuthMiddleware(), controllers.SearchPendudukByName)
		pd.POST("/", middleware.AuthMiddleware(), controllers.CreatePenduduk)
		pd.PUT("/:id", middleware.AuthMiddleware(), controllers.UpdatePenduduk)
		pd.DELETE("/:id", middleware.AuthMiddleware(), controllers.DeletePenduduk)
	}
}
