package controllers

import (
	"net/http"
	"strings"

	"desa-kepayang-backend/config"
	"desa-kepayang-backend/helpers"
	"desa-kepayang-backend/models"

	"github.com/gin-gonic/gin"
)

// ==================================
// ========== [CREATE] ==============
// ==================================

func CreateVisiMisi(c *gin.Context) {
	var input models.VisiMisi

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	// Sanitasi input
	input.Visi = helpers.SanitizeText(input.Visi)
	input.Misi = helpers.SanitizeText(input.Misi)

	// Validasi input
	if strings.TrimSpace(input.Visi) == "" || strings.TrimSpace(input.Misi) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Visi dan Misi tidak boleh kosong"})
		return
	}

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan visi misi ke database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Visi misi berhasil ditambahkan", "data": input})
}

// ==================================
// ========== [READ] ================
// ==================================

func GetAllVisiMisi(c *gin.Context) {
	var data []models.VisiMisi
	if err := config.DB.Find(&data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data visi misi"})
		return
	}
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

	var input models.VisiMisi
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	// Sanitasi dan validasi input
	input.Visi = helpers.SanitizeText(input.Visi)
	input.Misi = helpers.SanitizeText(input.Misi)

	if strings.TrimSpace(input.Visi) == "" || strings.TrimSpace(input.Misi) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Visi dan Misi tidak boleh kosong"})
		return
	}

	// Update data
	visimisi.Visi = input.Visi
	visimisi.Misi = input.Misi

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

	if err := config.DB.Delete(&visimisi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil dihapus"})
}
