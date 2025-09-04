package controllers

import (
	"context"
	"desa-kepayang-backend/config"
	"desa-kepayang-backend/helpers"
	"desa-kepayang-backend/models"
	"fmt"
	"net/http"
	"time"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

// ==================================
// =========== [CREATE] =============
// ==================================

func TambahSambutan(c *gin.Context) {
	ctx := context.Background()

	// Cek apakah sudah ada data sambutan
	var count int64
	config.DB.Model(&models.SambutanKepalaDesa{}).Count(&count)
	if count >= 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Hanya boleh ada 1 data sambutan. Hapus atau update data lama sebelum menambah baru.",
		})
		return
	}

	// Ambil dan validasi input teks
	kataSambutan := helpers.SanitizeText(c.PostForm("kata_sambutan"))
	namaKepalaDesa := helpers.SanitizeText(c.PostForm("nama_kepaladesa"))

	if kataSambutan == "" || namaKepalaDesa == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kata sambutan dan nama kepala desa wajib diisi"})
		return
	}

	// Ambil dan validasi file foto
	fileHeader, err := c.FormFile("foto_kepaladesa")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File foto wajib diunggah"})
		return
	}
	if fileHeader.Size > 2*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran file maksimal 2MB"})
		return
	}
	if !helpers.IsAllowedFileType(fileHeader.Filename) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ekstensi file tidak diizinkan"})
		return
	}

	src, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuka file"})
		return
	}
	defer src.Close()

	publicID := fmt.Sprintf("sambutan/%d_%s", time.Now().Unix(), helpers.RandomString(8))

	uploadRes, err := config.Cloudinary.Upload.Upload(ctx, src, uploader.UploadParams{
		PublicID:       publicID,
		Transformation: "c_fill,g_face,h_600,w_400", // fokus wajah & crop proporsional
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal upload gambar ke Cloudinary"})
		return
	}

	sambutan := models.SambutanKepalaDesa{
		Foto:           uploadRes.SecureURL,
		FotoID:         uploadRes.PublicID,
		KataSambutan:   kataSambutan,
		NamaKepalaDesa: namaKepalaDesa,
	}

	if err := config.DB.Create(&sambutan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan sambutan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sambutan berhasil ditambahkan", "data": sambutan})
}

// =================================
// =========== [READ] ==============
// =================================

func GetSambutan(c *gin.Context) {
	var sambutans []models.SambutanKepalaDesa

	if err := config.DB.Find(&sambutans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": sambutans})
}

// ==================================
// ========== [UPDATE] ==============
// ==================================

func UpdateSambutan(c *gin.Context) {
	ctx := context.Background()
	id := c.Param("id")

	var sambutan models.SambutanKepalaDesa
	if err := config.DB.First(&sambutan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}

	// Update teks jika ada
	kataSambutan := helpers.SanitizeText(c.PostForm("kata_sambutan"))
	namaKepalaDesa := helpers.SanitizeText(c.PostForm("nama_kepaladesa"))

	if kataSambutan != "" {
		sambutan.KataSambutan = kataSambutan
	}
	if namaKepalaDesa != "" {
		sambutan.NamaKepalaDesa = namaKepalaDesa
	}

	// Cek file foto baru
	fileHeader, err := c.FormFile("foto_kepaladesa")
	if err == nil {
		if fileHeader.Size > 2*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran file maksimal 2MB"})
			return
		}
		if !helpers.IsAllowedFileType(fileHeader.Filename) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ekstensi file tidak diizinkan"})
			return
		}

		// Hapus foto lama dari Cloudinary
		if sambutan.FotoID != "" {
			_, _ = config.Cloudinary.Upload.Destroy(ctx, uploader.DestroyParams{
				PublicID: sambutan.FotoID,
			})
		}

		src, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuka file"})
			return
		}
		defer src.Close()

		publicID := fmt.Sprintf("sambutan/%d_%s", time.Now().Unix(), helpers.RandomString(8))

		uploadRes, err := config.Cloudinary.Upload.Upload(ctx, src, uploader.UploadParams{
			PublicID:       publicID,
			Transformation: "c_fill,g_face,h_600,w_400", // fokus wajah & crop proporsional
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal upload gambar ke Cloudinary"})
			return
		}

		sambutan.Foto = uploadRes.SecureURL
		sambutan.FotoID = uploadRes.PublicID
	}

	if err := config.DB.Save(&sambutan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil diperbarui", "data": sambutan})
}

// ==================================
// ========== [DELETE] ==============
// ==================================

func DeleteSambutan(c *gin.Context) {
	ctx := context.Background()
	id := c.Param("id")

	var sambutan models.SambutanKepalaDesa
	if err := config.DB.First(&sambutan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}

	// Hapus foto di Cloudinary
	if sambutan.FotoID != "" {
		_, _ = config.Cloudinary.Upload.Destroy(ctx, uploader.DestroyParams{
			PublicID: sambutan.FotoID,
		})
	}

	if err := config.DB.Delete(&sambutan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil dihapus"})
}
