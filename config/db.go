package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"desa-kepayang-backend/models"
)

var DB *gorm.DB

func InitDB() {
	// Hanya muat .env jika dijalankan di lokal
	if os.Getenv("RAILWAY_ENVIRONMENT") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Peringatan: Gagal memuat file .env (tidak masalah jika di production)")
		}
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}

	err = db.AutoMigrate(
		&models.SambutanKepalaDesa{},
		&models.Admin{},
		&models.Berita{},
	)
	if err != nil {
		log.Fatal("Gagal migrasi database:", err)
	}

	DB = db
	fmt.Println("Berhasil koneksi dan migrasi ke database!")
}
