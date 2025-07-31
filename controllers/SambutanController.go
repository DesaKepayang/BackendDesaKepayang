package controllers

import (
	"desa-kepayang-backend/config"
	"desa-kepayang-backend/models"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// ==================================
// =========== [CREATE] =============
// ==================================

func TambahSambutan(c *gin.Context) {
	// Ambil file foto
	file, err := c.FormFile("foto_kepaladesa")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gambar tidak ditemukan", "detail": err.Error()})
		return
	}

	// Simpan file ke folder /uploads/
	path := "uploads/" + file.Filename
	if err := c.SaveUploadedFile(file, path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan gambar", "detail": err.Error()})
		return
	}

	// Ambil teks sambutan
	kataSambutan := c.PostForm("kata_sambutan")
	namaKepalaDesa := c.PostForm("nama_kepaladesa")

	// Buat objek dan simpan
	sambutan := models.SambutanKepalaDesa{
		Foto:           path,
		KataSambutan:   kataSambutan,
		NamaKepalaDesa: namaKepalaDesa, // Tambahan
	}

	if err := config.DB.Create(&sambutan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan sambutan", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sambutan berhasil ditambahkan", "data": sambutan})
}

// =================================
// =========== [READ] ==============
// =================================

func GetSambutan(c *gin.Context) {
	var sambutans []models.SambutanKepalaDesa

	if err := config.DB.Find(&sambutans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": sambutans})
}

// ==================================
// ========== [UPDATE] ==============
// ==================================

func UpdateSambutan(c *gin.Context) {
	id := c.Param("id")

	var sambutan models.SambutanKepalaDesa
	if err := config.DB.First(&sambutan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}

	// Cek apakah ada file baru
	file, err := c.FormFile("foto_kepaladesa")
	if err == nil {
		// Hapus file lama jika ada
		if sambutan.Foto != "" {
			os.Remove(sambutan.Foto)
		}

		// Simpan file baru
		path := "uploads/" + file.Filename
		if err := c.SaveUploadedFile(file, path); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan gambar", "detail": err.Error()})
			return
		}
		sambutan.Foto = path
	}

	// Update kata sambutan
	kataSambutan := c.PostForm("kata_sambutan")
	if kataSambutan != "" {
		sambutan.KataSambutan = kataSambutan
	}

	if err := config.DB.Save(&sambutan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data", "detail": err.Error()})
		return
	}

	namaKepalaDesa := c.PostForm("nama_kepaladesa")
	if namaKepalaDesa != "" {
		sambutan.NamaKepalaDesa = namaKepalaDesa
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil diperbarui", "data": sambutan})
}

// ==================================
// ========== [DELETE] ==============
// ==================================

func DeleteSambutan(c *gin.Context) {
	id := c.Param("id")

	var sambutan models.SambutanKepalaDesa
	if err := config.DB.First(&sambutan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}

	// Hapus file foto jika ada
	if sambutan.Foto != "" {
		os.Remove(sambutan.Foto)
	}

	if err := config.DB.Delete(&sambutan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil dihapus"})
}
