package main

import (
	"desa-kepayang-backend/config"
	"desa-kepayang-backend/models"
	"desa-kepayang-backend/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Inisialisasi DB dan Migrasi
	config.InitDB()
	err := config.DB.AutoMigrate(
		&models.SambutanKepalaDesa{},
		&models.Admin{},
		&models.Berita{},
		&models.VisiMisi{},
		&models.StrukturDesa{},
		&models.RTRW{},
		&models.DataPenduduk{},
	)
	if err != nil {
		log.Fatal("Gagal migrasi DB:", err)
	}

	// Inisialisasi router
	r := gin.Default()

	// Jadikan folder 'uploads/' sebagai folder statis
	r.Static("/uploads", "./uploads")

	// Registrasi routes sambutan
	routes.SambutanRoutes(r)
	routes.AdminRoutes(r)
	routes.BeritaRoutes(r)
	routes.VisiMisiRoutes(r)
	routes.StrukturDesaRoutes(r)
	routes.RTRWRoutes(r)
	routes.PendudukRoutes(r)

	// Root testing
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Halo dari Gin!"})
	})

	// Jalankan server
	r.Run()
}
