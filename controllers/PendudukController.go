package controllers

import (
	"desa-kepayang-backend/config"
	"desa-kepayang-backend/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ==================================
// =========== [CREATE] =============
// ==================================

func CreatePenduduk(c *gin.Context) {
	var input models.DataPenduduk

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cek apakah RTRW ada
	var rtrw models.RTRW
	if err := config.DB.First(&rtrw, input.IDRTRW).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID RT/RW tidak valid"})
		return
	}

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan data penduduk"})
		return
	}

	// Format response
	response := gin.H{
		"id_penduduk": input.IDPenduduk,
		"id_rtrw":     input.IDRTRW,
		"nama":        input.Nama,
		"agama":       input.Agama,
		"gender":      input.Gender,
		"rtrw": gin.H{
			"id_rtrw": rtrw.IDRTRW,
			"rt":      rtrw.RT,
			"rw":      rtrw.RW,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data penduduk berhasil ditambahkan",
		"data":    response,
	})
}

// ==================================
// =========== [READ] ===============
// ==================================

func GetAllPenduduk(c *gin.Context) {
	var pendudukList []models.DataPenduduk
	var rtrwList []models.RTRW

	// Ambil semua data
	if err := config.DB.Find(&pendudukList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data penduduk"})
		return
	}

	if err := config.DB.Find(&rtrwList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data RTRW"})
		return
	}

	// Buat map untuk lookup RTRW
	rtrwMap := make(map[uint]models.RTRW)
	for _, r := range rtrwList {
		rtrwMap[r.IDRTRW] = r
	}

	// Format response
	var response []gin.H
	for _, p := range pendudukList {
		rtrw := rtrwMap[p.IDRTRW]
		response = append(response, gin.H{
			"id_penduduk": p.IDPenduduk,
			"id_rtrw":     p.IDRTRW,
			"nama":        p.Nama,
			"agama":       p.Agama,
			"gender":      p.Gender,
			"rtrw": gin.H{
				"id_rtrw": rtrw.IDRTRW,
				"rt":      rtrw.RT,
				"rw":      rtrw.RW,
			},
		})
	}

	c.JSON(http.StatusOK, response)
}

func CountPenduduk(c *gin.Context) {
	var count int64

	if err := config.DB.Model(&models.DataPenduduk{}).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung jumlah penduduk"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"jumlah_penduduk": count,
	})
}

func CountPendudukByGender(c *gin.Context) {
	var lakiLaki int64
	var perempuan int64

	// Hitung jumlah laki-laki
	if err := config.DB.Model(&models.DataPenduduk{}).
		Where("gender = ?", "Laki-laki").
		Count(&lakiLaki).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung jumlah laki-laki"})
		return
	}

	// Hitung jumlah perempuan
	if err := config.DB.Model(&models.DataPenduduk{}).
		Where("gender = ?", "Perempuan").
		Count(&perempuan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung jumlah perempuan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"jumlah_laki_laki": lakiLaki,
		"jumlah_perempuan": perempuan,
	})
}

func CountPendudukByAgama(c *gin.Context) {
	type AgamaStat struct {
		Agama  string
		Jumlah int64
	}

	// Daftar agama yang ingin dihitung
	agamaList := []string{"Islam", "Konghucu", "Katolik", "Kristen", "Buddha", "Lain - Lain"}

	var totalPenduduk int64
	if err := config.DB.Model(&models.DataPenduduk{}).Count(&totalPenduduk).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung total penduduk"})
		return
	}

	results := make(map[string]string)

	for _, agama := range agamaList {
		var count int64
		if err := config.DB.Model(&models.DataPenduduk{}).
			Where("agama = ?", agama).
			Count(&count).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung jumlah penduduk untuk agama " + agama})
			return
		}

		// Hitung persentase
		var percentage float64
		if totalPenduduk > 0 {
			percentage = (float64(count) / float64(totalPenduduk)) * 100
		}

		results[agama] = fmt.Sprintf("%d (%.2f%%)", count, percentage)
	}

	c.JSON(http.StatusOK, results)
}

func CountPendudukPerRTRW(c *gin.Context) {
	// Ambil semua data RTRW
	var rtrwList []models.RTRW
	if err := config.DB.Find(&rtrwList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data RTRW"})
		return
	}

	// Hitung total penduduk keseluruhan
	var totalPenduduk int64
	if err := config.DB.Model(&models.DataPenduduk{}).Count(&totalPenduduk).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung total penduduk"})
		return
	}

	// Buat slice untuk menyimpan hasil
	type Result struct {
		IDRTRW        uint    `json:"id_rtrw"`
		RT            string  `json:"rt"`
		RW            string  `json:"rw"`
		Jumlah        int64   `json:"jumlah"`
		Persentase    float64 `json:"persentase"`
		FormattedText string  `json:"formatted_text"`
	}

	var results []Result

	// Hitung jumlah penduduk per RTRW
	for _, rtrw := range rtrwList {
		var count int64
		if err := config.DB.Model(&models.DataPenduduk{}).
			Where("id_rtrw = ?", rtrw.IDRTRW).
			Count(&count).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung penduduk per RTRW"})
			return
		}

		// Hitung persentase
		var percentage float64
		if totalPenduduk > 0 {
			percentage = (float64(count) / float64(totalPenduduk)) * 100
		}

		// Format teks sesuai permintaan
		formatted := fmt.Sprintf("RTRW : %s / %s %d (%.0f%%)",
			rtrw.RT, rtrw.RW, count, percentage)

		results = append(results, Result{
			IDRTRW:        rtrw.IDRTRW,
			RT:            rtrw.RT,
			RW:            rtrw.RW,
			Jumlah:        count,
			Persentase:    percentage,
			FormattedText: formatted,
		})
	}

	c.JSON(http.StatusOK, results)
}

// ==================================
// =========== [UPDATE] =============
// ==================================

func UpdatePenduduk(c *gin.Context) {
	id := c.Param("id")
	var data models.DataPenduduk

	// Cari data penduduk yang akan diupdate
	if err := config.DB.First(&data, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data penduduk tidak ditemukan"})
		return
	}

	// Bind input data
	var input struct {
		IDRTRW uint   `json:"id_rtrw"`
		Nama   string `json:"nama"`
		Agama  string `json:"agama"`
		Gender string `json:"gender"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi input
	if input.Nama == "" || input.Agama == "" || input.Gender == "" || input.IDRTRW == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Semua field harus diisi"})
		return
	}

	// Cek apakah RTRW baru valid
	var rtrw models.RTRW
	if err := config.DB.First(&rtrw, input.IDRTRW).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID RT/RW tidak valid"})
		return
	}

	// Update data
	data.IDRTRW = input.IDRTRW
	data.Nama = input.Nama
	data.Agama = input.Agama
	data.Gender = input.Gender

	if err := config.DB.Save(&data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data penduduk"})
		return
	}

	// Format response
	response := gin.H{
		"id_penduduk": data.IDPenduduk,
		"id_rtrw":     data.IDRTRW,
		"nama":        data.Nama,
		"agama":       data.Agama,
		"gender":      data.Gender,
		"rtrw": gin.H{
			"id_rtrw": rtrw.IDRTRW,
			"rt":      rtrw.RT,
			"rw":      rtrw.RW,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data penduduk berhasil diperbarui",
		"data":    response,
	})
}

// ==================================
// =========== [DELETE] =============
// ==================================

func DeletePenduduk(c *gin.Context) {
	id := c.Param("id")
	var data models.DataPenduduk

	if err := config.DB.First(&data, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data penduduk tidak ditemukan"})
		return
	}

	if err := config.DB.Delete(&data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data penduduk"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data penduduk berhasil dihapus"})
}
