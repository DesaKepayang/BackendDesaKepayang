package routes

import (
	"desa-kepayang-backend/controllers"
	"desa-kepayang-backend/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.Engine) {
	admin := r.Group("/admin")
	{
		admin.POST("/register", controllers.TambahAdmin)
		admin.POST("/login", controllers.LoginAdmin)

		adminAuth := admin.Group("/")
		adminAuth.Use(middleware.AuthMiddleware())
		{
			adminAuth.GET("/", controllers.GetAllAdmin)
			adminAuth.PUT("/:id", controllers.UpdateAdmin)
			adminAuth.DELETE("/:id", controllers.DeleteAdmin)
			adminAuth.GET("/me", controllers.GetAdminProfile)
			adminAuth.POST("/logout", controllers.LogoutAdmin)
		}
	}

}
