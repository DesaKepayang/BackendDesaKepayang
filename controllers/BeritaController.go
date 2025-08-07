package controllers

import (
	"desa-kepayang-backend/helpers"
	"net/http"
	"os"
	"path/filepath"

	"desa-kepayang-backend/config"
	"desa-kepayang-backend/models"

	"github.com/gin-gonic/gin"
)

// ==================================
// =========== [CREATE] =============
// ==================================

func CreateBerita(c *gin.Context) {
	judul := helpers.SanitizeText(c.PostForm("judul"))
	deskripsi := helpers.SanitizeText(c.PostForm("deskripsi"))
	tanggal := helpers.SanitizeText(c.PostForm("tanggal")) // Tambahan

	if judul == "" || deskripsi == "" || tanggal == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Judul, deskripsi, dan tanggal wajib diisi"})
		return
	}

	file, err := c.FormFile("foto")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File foto wajib diunggah"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file"})
		return
	}

	berita := models.Berita{
		Judul:     judul,
		Deskripsi: deskripsi,
		Foto:      path,
		Tanggal:   tanggal,
	}

	if err := config.DB.Create(&berita).Error; err != nil {
		os.Remove(path)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan berita"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Berita berhasil ditambahkan", "data": berita})
}

// ==================================
// =========== [READ] ===============
// ==================================

func GetAllBerita(c *gin.Context) {
	var daftarBerita []models.Berita

	if err := config.DB.Find(&daftarBerita).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data berita"})
		return
	}

	c.JSON(http.StatusOK, daftarBerita)
}

// ==================================
// =========== [UPDATE] =============
// ==================================

func UpdateBerita(c *gin.Context) {
	id := c.Param("id")
	var berita models.Berita

	if err := config.DB.First(&berita, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Berita tidak ditemukan"})
		return
	}

	judul := helpers.SanitizeText(c.PostForm("judul"))
	deskripsi := helpers.SanitizeText(c.PostForm("deskripsi"))
	tanggal := helpers.SanitizeText(c.PostForm("tanggal")) // Tambahan

	if judul != "" {
		berita.Judul = judul
	}
	if deskripsi != "" {
		berita.Deskripsi = deskripsi
	}
	if tanggal != "" {
		berita.Tanggal = tanggal
	}

	file, err := c.FormFile("foto")
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file baru"})
			return
		}

		if berita.Foto != "" {
			_ = os.Remove(berita.Foto)
		}
		berita.Foto = newPath
	}

	if err := config.DB.Save(&berita).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate berita"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Berita berhasil diupdate", "data": berita})
}

// ==================================
// =========== [DELETE] =============
// ==================================

func DeleteBerita(c *gin.Context) {
	id := c.Param("id")
	var berita models.Berita

	if err := config.DB.First(&berita, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Berita tidak ditemukan"})
		return
	}

	if err := config.DB.Delete(&berita).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus berita"})
		return
	}

	if berita.Foto != "" {
		_ = os.Remove(berita.Foto)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Berita berhasil dihapus"})
}
