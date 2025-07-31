package main

import (
	"desa-kepayang-backend/config"
	"desa-kepayang-backend/models"
	"desa-kepayang-backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Inisialisasi DB dan Migrasi
	config.InitDB()
	config.DB.AutoMigrate(
		&models.SambutanKepalaDesa{},
		&models.Admin{},
		&models.Berita{},
		&models.VisiMisi{},
		&models.StrukturDesa{},
		&models.RTRW{},
	)

	// Inisialisasi router
	r := gin.Default()

	// Registrasi routes sambutan
	routes.SambutanRoutes(r)
	routes.AdminRoutes(r)
	routes.BeritaRoutes(r)
	routes.VisiMisiRoutes(r)
	routes.StrukturDesaRoutes(r)
	routes.RTRWRoutes(r)

	// Root testing
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Halo dari Gin!"})
	})

	// Jalankan server
	r.Run()
}
