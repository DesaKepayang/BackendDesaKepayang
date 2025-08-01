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
