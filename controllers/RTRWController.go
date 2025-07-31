package controllers

import (
	"desa-kepayang-backend/config"
	"desa-kepayang-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ==================================
// =========== [DELETE] =============
// ==================================

func CreateRTRW(c *gin.Context) {
	rt := c.PostForm("rt")
	rw := c.PostForm("rw")

	data := models.RTRW{
		RT: rt,
		RW: rw,
	}

	if err := config.DB.Create(&data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan data RT/RW"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "RT/RW berhasil ditambahkan", "data": data})
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

	rt := c.PostForm("rt")
	rw := c.PostForm("rw")

	data.RT = rt
	data.RW = rw

	if err := config.DB.Save(&data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data RT/RW"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data RT/RW berhasil diperbarui", "data": data})
}

// ==================================
// =========== [CREATE] =============
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
