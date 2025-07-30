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
	)

	// Inisialisasi router
	r := gin.Default()

	// Registrasi routes sambutan
	routes.SambutanRoutes(r)
	routes.AdminRoutes(r)

	// Root testing
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Halo dari Gin!"})
	})

	// Jalankan server
	r.Run()
}
