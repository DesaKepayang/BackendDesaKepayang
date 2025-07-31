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
	judul := c.PostForm("judul")
	deskripsi := c.PostForm("deskripsi")

	// Ambil file gambar dari form-data
	file, err := c.FormFile("foto")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal mengambil file foto"})
		return
	}

	// Simpan file ke folder lokal
	path := "uploads/" + file.Filename
	if err := c.SaveUploadedFile(file, path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file"})
		return
	}

	// Simpan data berita ke database
	berita := models.Berita{
		Judul:     judul,
		Deskripsi: deskripsi,
		Foto:      path,
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

	// Cari data berita berdasarkan ID
	if err := config.DB.First(&berita, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Berita tidak ditemukan"})
		return
	}

	// Ambil data form
	judul := c.PostForm("judul")
	deskripsi := c.PostForm("deskripsi")

	// Perbarui data judul dan deskripsi
	berita.Judul = judul
	berita.Deskripsi = deskripsi

	// Cek apakah ada file baru dikirim
	file, err := c.FormFile("foto")
	if err == nil {
		// Jika ada file baru, simpan ke folder dan perbarui path
		path := "uploads/" + file.Filename
		if err := c.SaveUploadedFile(file, path); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file baru"})
			return
		}
		berita.Foto = path
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
