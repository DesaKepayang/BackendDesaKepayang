package controllers

import (
	"net/http"

	"desa-kepayang-backend/config"
	"desa-kepayang-backend/models"

	"github.com/gin-gonic/gin"
)

// ==================================
// =========== [CREATE] =============
// ==================================

func CreateBerita(c *gin.Context) {
	var berita models.Berita

	if err := c.ShouldBindJSON(&berita); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&berita).Error; err != nil {
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

	// Cek apakah data berita dengan ID tersebut ada
	if err := config.DB.First(&berita, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Berita tidak ditemukan"})
		return
	}

	// Bind data JSON ke struct berita
	if err := c.ShouldBindJSON(&berita); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Simpan perubahan
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

	// Cari berita berdasarkan ID
	if err := config.DB.First(&berita, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Berita tidak ditemukan"})
		return
	}

	// Hapus berita
	if err := config.DB.Delete(&berita).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus berita"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Berita berhasil dihapus"})
}
