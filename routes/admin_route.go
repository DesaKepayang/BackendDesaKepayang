package routes

import (
	"desa-kepayang-backend/controllers"
	"desa-kepayang-backend/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.Engine) {
	admin := r.Group("/admin")
	{
		admin.POST("/register", controllers.TambahAdmin) // tanpa auth
		admin.POST("/login", controllers.LoginAdmin)     // login

		adminAuth := admin.Group("/")
		adminAuth.Use(middleware.AuthMiddleware())
		{
			admin.OPTIONS("/logout", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})
			adminAuth.GET("/", controllers.GetAllAdmin)
			adminAuth.PUT("/:id", controllers.UpdateAdmin)
			adminAuth.DELETE("/:id", controllers.DeleteAdmin)
			adminAuth.GET("/me", controllers.GetAdminProfile)
		}
	}

}
