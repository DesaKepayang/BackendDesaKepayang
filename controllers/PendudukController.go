package controllers

import (
	"desa-kepayang-backend/config"
	"desa-kepayang-backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ==================================
// =========== [CREATE] =============
// ==================================

func CreatePenduduk(c *gin.Context) {
	var input models.DataPenduduk
	contentType := c.GetHeader("Content-Type")

	if contentType == "application/json" {
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format JSON tidak valid"})
			return
		}
	} else {
		idRTRW, _ := strconv.Atoi(c.PostForm("id_rtrw"))
		input.IDRTRW = uint(idRTRW)
		input.Nama = c.PostForm("nama")
		input.Agama = c.PostForm("agama")
		input.Gender = c.PostForm("gender")
	}

	if input.Nama == "" || input.Agama == "" || input.Gender == "" || input.IDRTRW == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Semua field harus diisi"})
		return
	}

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan data penduduk"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data penduduk berhasil ditambahkan", "data": input})
}

// ==================================
// =========== [READ] ===============
// ==================================

func GetAllPenduduk(c *gin.Context) {
	var data []models.DataPenduduk
	if err := config.DB.Preload("Penduduk.RTRW").Find(&data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data penduduk"})
		return
	}
	c.JSON(http.StatusOK, data)
}

// ==================================
// =========== [UPDATE] =============
// ==================================

func UpdatePenduduk(c *gin.Context) {
	id := c.Param("id")
	var data models.DataPenduduk

	if err := config.DB.First(&data, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data penduduk tidak ditemukan"})
		return
	}

	contentType := c.GetHeader("Content-Type")
	var input models.DataPenduduk

	if contentType == "application/json" {
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format JSON tidak valid"})
			return
		}
	} else {
		idRTRW, _ := strconv.Atoi(c.PostForm("id_rtrw"))
		input.IDRTRW = uint(idRTRW)
		input.Nama = c.PostForm("nama")
		input.Agama = c.PostForm("agama")
		input.Gender = c.PostForm("gender")
	}

	if input.Nama == "" || input.Agama == "" || input.Gender == "" || input.IDRTRW == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Semua field harus diisi"})
		return
	}

	data.IDRTRW = input.IDRTRW
	data.Nama = input.Nama
	data.Agama = input.Agama
	data.Gender = input.Gender

	if err := config.DB.Save(&data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data penduduk"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data penduduk berhasil diperbarui", "data": data})
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
