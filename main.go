package main

import (
	"desa-kepayang-backend/config"
	"desa-kepayang-backend/controllers"
	"desa-kepayang-backend/models"
	"desa-kepayang-backend/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Inisialisasi DB dan Migrasi
	config.InitDB()
	db := config.DB

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

	err = db.Exec(`
		ALTER TABLE data_penduduk
		ADD CONSTRAINT fk_penduduk_rtrw
		FOREIGN KEY (id_rtrw) REFERENCES rt_rw(id_rtrw)
		ON DELETE RESTRICT
		ON UPDATE CASCADE
	`).Error

	if err != nil {
		log.Println("Peringatan: Gagal menambahkan foreign key (mungkin sudah ada):", err)
	}

	// Inisialisasi router
	r := gin.Default()

	r.POST("/data_penduduk", controllers.CreatePenduduk)
	r.GET("/data_penduduk", controllers.GetAllPenduduk)
	r.PUT("/data_penduduk/:id", controllers.UpdatePenduduk)
	r.DELETE("/data_penduduk/:id", controllers.DeletePenduduk)

	// Jadikan folder 'uploads/' sebagai folder statis
	r.Static("/uploads", "./uploads")

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
