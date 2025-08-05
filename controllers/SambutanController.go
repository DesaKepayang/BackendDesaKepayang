package controllers

import (
	"desa-kepayang-backend/config"
	"desa-kepayang-backend/helpers"
	"desa-kepayang-backend/models"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// ==================================
// =========== [CREATE] =============
// ==================================

func TambahSambutan(c *gin.Context) {
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

	// Ambil dan validasi file
	file, err := c.FormFile("foto_kepaladesa")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File foto wajib diunggah", "detail": err.Error()})
		return
	}

	if file.Size > 2*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran file maksimal 2MB"})
		return
	}

	uniqueFileName := helpers.GenerateUniqueFileName(file.Filename)
	if uniqueFileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ekstensi file tidak diizinkan"})
		return
	}

	path := filepath.Join("uploads", uniqueFileName)
	if err := c.SaveUploadedFile(file, path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan gambar", "detail": err.Error()})
		return
	}

	// Buat dan simpan data
	sambutan := models.SambutanKepalaDesa{
		Foto:           path,
		KataSambutan:   kataSambutan,
		NamaKepalaDesa: namaKepalaDesa,
	}

	if err := config.DB.Create(&sambutan).Error; err != nil {
		os.Remove(path) // hapus file jika DB gagal
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan sambutan", "detail": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data", "detail": err.Error()})
		return
	}

	// Ambil semua path yang sedang digunakan
	var usedPaths []string
	for _, s := range sambutans {
		usedPaths = append(usedPaths, s.Foto)
	}

	// Bersihkan file yang tidak digunakan
	helpers.CleanupUnusedFiles("uploads", usedPaths)

	c.JSON(http.StatusOK, gin.H{"data": sambutans})
}

// ==================================
// ========== [UPDATE] ==============
// ==================================

func UpdateSambutan(c *gin.Context) {
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

	// Cek dan proses file baru jika ada
	file, err := c.FormFile("foto_kepaladesa")
	if err == nil {
		if file.Size > 2*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran file maksimal 2MB"})
			return
		}

		uniqueFileName := helpers.GenerateUniqueFileName(file.Filename)
		if uniqueFileName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ekstensi file tidak diizinkan"})
			return
		}

		newPath := filepath.Join("uploads", uniqueFileName)
		if err := c.SaveUploadedFile(file, newPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan gambar baru", "detail": err.Error()})
			return
		}

		// Hapus file lama
		if sambutan.Foto != "" {
			_ = os.Remove(sambutan.Foto)
		}
		sambutan.Foto = newPath
	}

	if err := config.DB.Save(&sambutan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil diperbarui", "data": sambutan})
}

// ==================================
// ========== [DELETE] ==============
// ==================================

func DeleteSambutan(c *gin.Context) {
	id := c.Param("id")

	var sambutan models.SambutanKepalaDesa
	if err := config.DB.First(&sambutan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}

	if err := config.DB.Delete(&sambutan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data", "detail": err.Error()})
		return
	}

	if sambutan.Foto != "" {
		_ = os.Remove(sambutan.Foto)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil dihapus"})
}
