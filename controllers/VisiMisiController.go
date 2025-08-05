package controllers

import (
	"net/http"
	"os"
	"path/filepath"

	"desa-kepayang-backend/config"
	"desa-kepayang-backend/helpers"
	"desa-kepayang-backend/models"

	"github.com/gin-gonic/gin"
)

// ==================================
// ========== [CREATE] ==============
// ==================================

func CreateVisiMisi(c *gin.Context) {
	visi := helpers.SanitizeText(c.PostForm("visi"))
	misi := helpers.SanitizeText(c.PostForm("misi"))

	if visi == "" || misi == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Visi dan misi wajib diisi"})
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

	visimisi := models.VisiMisi{
		Visi: visi,
		Misi: misi,
		Foto: path,
	}

	if err := config.DB.Create(&visimisi).Error; err != nil {
		os.Remove(path) // Hapus file jika gagal simpan DB
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan visi misi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Visi misi berhasil ditambahkan", "data": visimisi})
}

// ==================================
// ========== [READ] ================
// ==================================

func GetAllVisiMisi(c *gin.Context) {
	var data []models.VisiMisi
	if err := config.DB.Find(&data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}

	var usedPaths []string
	for _, v := range data {
		usedPaths = append(usedPaths, v.Foto)
	}

	helpers.CleanupUnusedFiles("uploads", usedPaths) // Pakai dari helpers
	c.JSON(http.StatusOK, data)
}

// ==================================
// ========== [UPDATE] ==============
// ==================================

func UpdateVisiMisi(c *gin.Context) {
	id := c.Param("id")
	var visimisi models.VisiMisi

	if err := config.DB.First(&visimisi, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}

	visi := helpers.SanitizeText(c.PostForm("visi"))
	misi := helpers.SanitizeText(c.PostForm("misi"))

	if visi != "" {
		visimisi.Visi = visi
	}
	if misi != "" {
		visimisi.Misi = misi
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
		if visimisi.Foto != "" {
			_ = os.Remove(visimisi.Foto)
		}
		visimisi.Foto = newPath
	}

	if err := config.DB.Save(&visimisi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil diperbarui", "data": visimisi})
}

// ==================================
// ========== [DELETE] ==============
// ==================================

func DeleteVisiMisi(c *gin.Context) {
	id := c.Param("id")
	var visimisi models.VisiMisi

	if err := config.DB.First(&visimisi, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}

	filePath := visimisi.Foto // Simpan path sebelum dihapus

	if err := config.DB.Delete(&visimisi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data"})
		return
	}

	// Hapus file setelah sukses hapus data
	if filePath != "" {
		_ = os.Remove(filePath)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil dihapus"})
}
