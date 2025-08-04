package main

import (
	"desa-kepayang-backend/config"
	"desa-kepayang-backend/middleware"
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
		&models.JumlahKK{},
		&models.Komentar{},
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
	r.Use(middleware.CORSMiddleware())

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
	routes.JumlahKKRoutes(r)
	routes.KomentarRoutes(r)

	// Jalankan server
	r.Run()
}
