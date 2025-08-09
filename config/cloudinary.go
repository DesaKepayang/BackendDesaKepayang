package config

import (
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
)

var Cloudinary *cloudinary.Cloudinary

func InitCloudinary() {
	var err error
	// Pastikan CLOUDINARY_URL di Railway / .env, contoh:
	// CLOUDINARY_URL=cloudinary://API_KEY:API_SECRET@CLOUD_NAME
	Cloudinary, err = cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	if err != nil {
		log.Fatal("Gagal koneksi ke Cloudinary:", err)
	}
}
