package controllers

import (
	"desa-kepayang-backend/config"
	"desa-kepayang-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ==================================
// =========== [CREATE] =============
// ==================================

func CreateJumlahKK(c *gin.Context) {
	var count int64
	config.DB.Model(&models.JumlahKK{}).Count(&count)
	if count >= 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Hanya boleh ada 1 data jumlah KK. Hapus atau update data lama sebelum menambah baru."})
		return
	}

	var input struct {
		JumlahKK int `json:"jumlahkk"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	// Validasi nilai
	if input.JumlahKK <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Jumlah KK harus lebih dari 0"})
		return
	}

	data := models.JumlahKK{
		JumlahKK: input.JumlahKK,
	}

	if err := config.DB.Create(&data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan jumlah KK"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Jumlah KK berhasil ditambahkan", "data": data})
}

// ==================================
// =========== [READ] ===============
// ==================================

func GetAllJumlahKK(c *gin.Context) {
	var list []models.JumlahKK
	if err := config.DB.Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data jumlah KK"})
		return
	}
	c.JSON(http.StatusOK, list)
}

// READ (satu)
func GetJumlahKKByID(c *gin.Context) {
	id := c.Param("id")
	var data models.JumlahKK
	if err := config.DB.First(&data, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, data)
}

// ==================================
// =========== [UPDATE] =============
// ==================================

func UpdateJumlahKK(c *gin.Context) {
	id := c.Param("id")
	var data models.JumlahKK

	if err := config.DB.First(&data, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}

	var input struct {
		JumlahKK int `json:"jumlahkk"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	if input.JumlahKK <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Jumlah KK harus lebih dari 0"})
		return
	}

	data.JumlahKK = input.JumlahKK

	if err := config.DB.Save(&data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Jumlah KK berhasil diperbarui", "data": data})
}

// ==================================
// =========== [DELETE] =============
// ==================================

func DeleteJumlahKK(c *gin.Context) {
	id := c.Param("id")
	var data models.JumlahKK

	if err := config.DB.First(&data, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}

	if err := config.DB.Delete(&data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Jumlah KK berhasil dihapus"})
}
