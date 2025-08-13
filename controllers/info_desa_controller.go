package controllers

import (
	"desa-kepayang-backend/config"
	"desa-kepayang-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CREATE
func CreateInfoDesa(c *gin.Context) {
	var count int64
	config.DB.Model(&models.InfoDesa{}).Count(&count)
	if count >= 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Hanya boleh ada 1 data info desa. Hapus atau update data lama sebelum menambah baru."})
		return
	}

	var input models.InfoDesa
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	if input.JumlahKK <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Jumlah KK harus lebih dari 0"})
		return
	}

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan info desa"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Info desa berhasil ditambahkan", "data": input})
}

// READ ALL
func GetAllInfoDesa(c *gin.Context) {
	var list []models.InfoDesa
	if err := config.DB.Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data info desa"})
		return
	}
	c.JSON(http.StatusOK, list)
}

// READ BY ID
func GetInfoDesaByID(c *gin.Context) {
	id := c.Param("id")
	var data models.InfoDesa
	if err := config.DB.First(&data, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, data)
}

// UPDATE
func UpdateInfoDesa(c *gin.Context) {
	id := c.Param("id")
	var data models.InfoDesa

	if err := config.DB.First(&data, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}

	// Simpan ID lama
	existingID := data.IDInfo

	var input models.InfoDesa
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	// Timpa semua data dengan input baru
	data = input
	data.IDInfo = existingID // kembalikan ID lama

	if err := config.DB.Save(&data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Info desa berhasil diperbarui", "data": data})
}

// DELETE
func DeleteInfoDesa(c *gin.Context) {
	id := c.Param("id")
	var data models.InfoDesa

	if err := config.DB.First(&data, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}

	if err := config.DB.Delete(&data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Info desa berhasil dihapus"})
}
