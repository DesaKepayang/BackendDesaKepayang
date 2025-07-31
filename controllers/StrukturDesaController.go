package controllers

import (
	"desa-kepayang-backend/config"
	"desa-kepayang-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ==================================
// ========== [CREATE] ==============
// ==================================

func CreateStrukturDesa(c *gin.Context) {
	nama := c.PostForm("nama")
	jabatan := c.PostForm("jabatan")

	file, err := c.FormFile("foto")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal mengambil file foto"})
		return
	}

	path := "uploads/" + file.Filename
	if err := c.SaveUploadedFile(file, path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file"})
		return
	}

	struktur := models.StrukturDesa{
		Nama:    nama,
		Jabatan: jabatan,
		Foto:    path,
	}

	if err := config.DB.Create(&struktur).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan struktur desa"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Struktur desa berhasil ditambahkan", "data": struktur})
}

// ==================================
// ========== [READ] ================
// ==================================

func GetAllStrukturDesa(c *gin.Context) {
	var data []models.StrukturDesa
	if err := config.DB.Find(&data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}
	c.JSON(http.StatusOK, data)
}

// ==================================
// ========== [UPDATE] ==============
// ==================================

func UpdateStrukturDesa(c *gin.Context) {
	id := c.Param("id")
	var struktur models.StrukturDesa

	if err := config.DB.First(&struktur, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Struktur desa tidak ditemukan"})
		return
	}

	nama := c.PostForm("nama")
	jabatan := c.PostForm("jabatan")

	struktur.Nama = nama
	struktur.Jabatan = jabatan

	file, err := c.FormFile("foto")
	if err == nil {
		path := "uploads/" + file.Filename
		if err := c.SaveUploadedFile(file, path); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file baru"})
			return
		}
		struktur.Foto = path
	}

	if err := config.DB.Save(&struktur).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate struktur desa"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Struktur desa berhasil diupdate", "data": struktur})
}

// ==================================
// ========== [DELETE] ==============
// ==================================

func DeleteStrukturDesa(c *gin.Context) {
	id := c.Param("id")
	var struktur models.StrukturDesa

	if err := config.DB.First(&struktur, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Struktur desa tidak ditemukan"})
		return
	}

	if err := config.DB.Delete(&struktur).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus struktur desa"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Struktur desa berhasil dihapus"})
}
