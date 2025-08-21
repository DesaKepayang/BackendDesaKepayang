package controllers

import (
	"desa-kepayang-backend/config"
	"desa-kepayang-backend/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ==================================
// ========== [CREATE] ==============
// ==================================

func CreateInfoDesa(c *gin.Context) {
	var input models.InfoDesa
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	if input.Indikator == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Indikator tidak boleh kosong"})
		return
	}

	if input.Jumlah < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Jumlah tidak boleh negatif"})
		return
	}

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan info desa"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Info desa berhasil ditambahkan", "data": input})
}

// ==================================
// ========== [READ] ================
// ==================================

func GetAllInfoDesa(c *gin.Context) {
	var list []models.InfoDesa
	if err := config.DB.Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data info desa"})
		return
	}
	c.JSON(http.StatusOK, list)
}

func GetInfoDasar(c *gin.Context) {
	var infoDesa []models.InfoDesa

	// Ambil semua data indikator yang diperlukan
	if err := config.DB.Where("indikator IN ?", []string{
		"Jumlah Penduduk", "Laki-laki", "Perempuan", "Jumlah KK",
	}).Find(&infoDesa).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data info dasar"})
		return
	}

	// Buat map untuk memudahkan akses
	result := make(map[string]int)
	for _, item := range infoDesa {
		result[item.Indikator] = item.Jumlah
	}

	c.JSON(http.StatusOK, gin.H{
		"Jumlah Penduduk": result["Jumlah Penduduk"],
		"Laki-laki":       result["Laki-laki"],
		"Perempuan":       result["Perempuan"],
		"Jumlah KK":       result["Jumlah KK"],
	})
}

func GetInfoAgama(c *gin.Context) {
	var infoDesa []models.InfoDesa

	// Ambil data indikator agama
	agamaList := []string{"Islam", "Protestan", "Katolik", "Buddha", "Konghucu", "Lain-lain"}
	if err := config.DB.Where("indikator IN ?", agamaList).Find(&infoDesa).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data agama"})
		return
	}

	// Hitung total penduduk
	totalPenduduk := 0
	for _, item := range infoDesa {
		totalPenduduk += item.Jumlah
	}

	// Format hasil
	result := make(map[string]string)
	for _, agama := range agamaList {
		jumlah := 0
		for _, item := range infoDesa {
			if item.Indikator == agama {
				jumlah = item.Jumlah
				break
			}
		}
		// Hitung persentase
		percentage := 0.0
		if totalPenduduk > 0 {
			percentage = (float64(jumlah) / float64(totalPenduduk)) * 100
		}
		result[agama] = fmt.Sprintf("%d (%.2f%%)", jumlah, percentage)
	}

	c.JSON(http.StatusOK, result)
}

func GetInfoRTRW(c *gin.Context) {
	var infoDesa []models.InfoDesa

	// Ambil semua data dengan format "NN / NN"
	if err := config.DB.
		Where("indikator LIKE ?", "__ / __").
		Find(&infoDesa).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data RT/RW"})
		return
	}

	// Hitung total penduduk di semua RT/RW
	totalPenduduk := 0
	for _, item := range infoDesa {
		totalPenduduk += item.Jumlah
	}

	// Format hasil
	result := make(map[string]string)
	for _, item := range infoDesa {
		percentage := 0.0
		if totalPenduduk > 0 {
			percentage = (float64(item.Jumlah) / float64(totalPenduduk)) * 100
		}
		result[item.Indikator] = fmt.Sprintf("%d (%.2f%%)", item.Jumlah, percentage)
	}

	c.JSON(http.StatusOK, result)
}

// ==================================
// ========== [UPDATE] ==============
// ==================================

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

	if input.Indikator == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Indikator tidak boleh kosong"})
		return
	}

	if input.Jumlah < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Jumlah tidak boleh negatif"})
		return
	}

	// Timpa semua data dengan input baru, tetapi pertahankan ID lama
	data = input
	data.IDInfo = existingID

	if err := config.DB.Save(&data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Info desa berhasil diperbarui", "data": data})
}

// ==================================
// ========== [DELETE] ==============
// ==================================

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
