package controllers

import (
	"net/http"

	"desa-kepayang-backend/config"
	"desa-kepayang-backend/models"

	"github.com/gin-gonic/gin"
)

// ==================================
// ========== [CREATE] ==============
// ==================================

func CreateVisiMisi(c *gin.Context) {
	var visimisi models.VisiMisi

	if err := c.ShouldBindJSON(&visimisi); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&visimisi).Error; err != nil {
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

	if err := c.ShouldBindJSON(&visimisi); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
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

	if err := config.DB.Delete(&visimisi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil dihapus"})
}
