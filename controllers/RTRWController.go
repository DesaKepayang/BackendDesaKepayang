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
// =========== [CREATE] =============
// ==================================

func CreateRTRW(c *gin.Context) {
	var input models.RTRW

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	// Sanitasi input
	input.RT = helpers.SanitizeText(input.RT)
	input.RW = helpers.SanitizeText(input.RW)

	// Validasi input
	if strings.TrimSpace(input.RT) == "" || strings.TrimSpace(input.RW) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "RT dan RW tidak boleh kosong"})
		return
	}

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan data RT/RW"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "RT/RW berhasil ditambahkan", "data": input})
}

// ==================================
// =========== [READ] ===============
// ==================================

func GetAllRTRW(c *gin.Context) {
	var data []models.RTRW
	if err := config.DB.Find(&data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data RT/RW"})
		return
	}
	c.JSON(http.StatusOK, data)
}

// ==================================
// =========== [UPDATE] =============
// ==================================

func UpdateRTRW(c *gin.Context) {
	id := c.Param("id")
	var data models.RTRW

	if err := config.DB.First(&data, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data RT/RW tidak ditemukan"})
		return
	}

	var input models.RTRW
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	// Sanitasi & Validasi
	input.RT = helpers.SanitizeText(input.RT)
	input.RW = helpers.SanitizeText(input.RW)

	if strings.TrimSpace(input.RT) == "" || strings.TrimSpace(input.RW) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "RT dan RW tidak boleh kosong"})
		return
	}

	data.RT = input.RT
	data.RW = input.RW

	if err := config.DB.Save(&data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data RT/RW"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data RT/RW berhasil diperbarui", "data": data})
}

// ==================================
// =========== [DELETE] =============
// ==================================

func DeleteRTRW(c *gin.Context) {
	id := c.Param("id")
	var data models.RTRW

	if err := config.DB.First(&data, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data RT/RW tidak ditemukan"})
		return
	}

	if err := config.DB.Delete(&data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data RT/RW"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data RT/RW berhasil dihapus"})
}
