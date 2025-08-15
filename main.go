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
	// Inisialisasi DB
	config.InitDB()

	// Inisialisasi Cloudinary
	config.InitCloudinary()

	db := config.DB

	// Migrasi tabel baru
	err := config.DB.AutoMigrate(
		&models.SambutanKepalaDesa{},
		&models.Admin{},
		&models.Berita{},
		&models.VisiMisi{},
		&models.StrukturDesa{},
		&models.RTRW{},
		&models.DataPenduduk{},
		&models.InfoDesa{}, // sekarang sudah struktur baru (3 kolom)
		&models.Komentar{},
	)

	if err != nil {
		log.Fatal("Gagal migrasi DB:", err)
	}

	// Tambah foreign key
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

	// ====== ROUTER ======
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	// Routes
	routes.SambutanRoutes(r)
	routes.AdminRoutes(r)
	routes.BeritaRoutes(r)
	routes.VisiMisiRoutes(r)
	routes.StrukturDesaRoutes(r)
	routes.RTRWRoutes(r)
	routes.PendudukRoutes(r)
	routes.InfoDesaRoutes(r)
	routes.KomentarRoutes(r)
	routes.WebSocketRoutes(r)

	// Jalankan server
	r.Run()
}
