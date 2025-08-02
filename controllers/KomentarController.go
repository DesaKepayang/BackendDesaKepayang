package controllers

import (
	"desa-kepayang-backend/config"
	"desa-kepayang-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CREATE
func CreateKomentar(c *gin.Context) {
	var input models.Komentar

	if err := c.ShouldBindJSON(&input); err != nil ||
		input.Nama == "" || input.Email == "" || input.NoHP == "" || input.Komentar == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Semua field harus diisi"})
		return
	}

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan komentar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Komentar berhasil dikirim", "data": input})
}

// READ (semua)
func GetAllKomentar(c *gin.Context) {
	var list []models.Komentar
	if err := config.DB.Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil komentar"})
		return
	}
	c.JSON(http.StatusOK, list)
}

// READ (satu)
func GetKomentarByID(c *gin.Context) {
	id := c.Param("id")
	var komentar models.Komentar

	if err := config.DB.First(&komentar, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Komentar tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, komentar)
}

// UPDATE
func UpdateKomentar(c *gin.Context) {
	id := c.Param("id")
	var komentar models.Komentar

	if err := config.DB.First(&komentar, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Komentar tidak ditemukan"})
		return
	}

	var input models.Komentar
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	komentar.Nama = input.Nama
	komentar.Email = input.Email
	komentar.NoHP = input.NoHP
	komentar.Komentar = input.Komentar

	if err := config.DB.Save(&komentar).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui komentar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Komentar berhasil diperbarui", "data": komentar})
}

// DELETE
func DeleteKomentar(c *gin.Context) {
	id := c.Param("id")
	var komentar models.Komentar

	if err := config.DB.First(&komentar, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Komentar tidak ditemukan"})
		return
	}

	if err := config.DB.Delete(&komentar).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus komentar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Komentar berhasil dihapus"})
}
