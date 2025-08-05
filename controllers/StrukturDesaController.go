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
// ========== [CREATE] ==============
// ==================================

func CreateStrukturDesa(c *gin.Context) {
	nama := helpers.SanitizeText(c.PostForm("nama"))
	jabatan := helpers.SanitizeText(c.PostForm("jabatan"))

	if nama == "" || jabatan == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama dan jabatan wajib diisi"})
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

	struktur := models.StrukturDesa{
		Nama:    nama,
		Jabatan: jabatan,
		Foto:    path,
	}

	if err := config.DB.Create(&struktur).Error; err != nil {
		os.Remove(path) // cleanup file jika gagal simpan DB
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan struktur desa"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Struktur desa berhasil ditambahkan", "data": struktur})
}

// ==================================
// ========== [READ] ================
// ==================================

func GetAllStrukturDesa(c *gin.Context) {
	var data []models.StrukturDesa
	if err := config.DB.Find(&data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}

	var usedPaths []string
	for _, s := range data {
		usedPaths = append(usedPaths, s.Foto)
	}

	helpers.CleanupUnusedFiles("uploads", usedPaths) // Pakai dari helpers

	c.JSON(http.StatusOK, data)
}

// ==================================
// ========== [UPDATE] ==============
// ==================================

func UpdateStrukturDesa(c *gin.Context) {
	id := c.Param("id")

	var struktur models.StrukturDesa
	if err := config.DB.First(&struktur, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Struktur desa tidak ditemukan"})
		return
	}

	nama := helpers.SanitizeText(c.PostForm("nama"))
	jabatan := helpers.SanitizeText(c.PostForm("jabatan"))

	if nama != "" {
		struktur.Nama = nama
	}
	if jabatan != "" {
		struktur.Jabatan = jabatan
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

		// Hapus file lama
		if struktur.Foto != "" {
			_ = os.Remove(struktur.Foto)
		}
		struktur.Foto = newPath
	}

	if err := config.DB.Save(&struktur).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate struktur desa"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Struktur desa berhasil diupdate", "data": struktur})
}

// ==================================
// ========== [DELETE] ==============
// ==================================

func DeleteStrukturDesa(c *gin.Context) {
	id := c.Param("id")

	var struktur models.StrukturDesa
	if err := config.DB.First(&struktur, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Struktur desa tidak ditemukan"})
		return
	}

	if err := config.DB.Delete(&struktur).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus struktur desa"})
		return
	}

	if struktur.Foto != "" {
		_ = os.Remove(struktur.Foto)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Struktur desa berhasil dihapus"})
}
