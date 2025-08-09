package controllers

import (
	"context"
	"desa-kepayang-backend/helpers"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"desa-kepayang-backend/config"
	"desa-kepayang-backend/models"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

// fungsi bantu untuk membuat string acak
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// ==================================
// =========== [CREATE] =============
// ==================================

func CreateBerita(c *gin.Context) {
	ctx := context.Background()

	judul := helpers.SanitizeText(c.PostForm("judul"))
	deskripsi := helpers.SanitizeText(c.PostForm("deskripsi"))
	tanggal := helpers.SanitizeText(c.PostForm("tanggal"))

	if judul == "" || deskripsi == "" || tanggal == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Judul, deskripsi, dan tanggal wajib diisi"})
		return
	}

	fileHeader, err := c.FormFile("foto")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File foto wajib diunggah"})
		return
	}

	if fileHeader.Size > 2*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran file maksimal 2MB"})
		return
	}

	if !helpers.IsAllowedFileType(fileHeader.Filename) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ekstensi file tidak diizinkan"})
		return
	}

	src, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuka file"})
		return
	}
	defer src.Close()

	publicID := fmt.Sprintf("berita/%d_%s", time.Now().Unix(), helpers.RandomString(8))

	uploadRes, err := config.Cloudinary.Upload.Upload(ctx, src, uploader.UploadParams{
		PublicID: publicID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal upload gambar ke Cloudinary"})
		return
	}

	berita := models.Berita{
		Judul:     judul,
		Deskripsi: deskripsi,
		Foto:      uploadRes.SecureURL,
		FotoID:    uploadRes.PublicID,
		Tanggal:   tanggal,
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

func GetBeritaByID(c *gin.Context) {
	id := c.Param("id")
	var berita models.Berita
	if err := config.DB.First(&berita, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Berita tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, berita)
}

// ==================================
// =========== [UPDATE] =============
// ==================================

func UpdateBerita(c *gin.Context) {
	ctx := context.Background()
	id := c.Param("id")

	var berita models.Berita
	if err := config.DB.First(&berita, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Berita tidak ditemukan"})
		return
	}

	judul := helpers.SanitizeText(c.PostForm("judul"))
	deskripsi := helpers.SanitizeText(c.PostForm("deskripsi"))
	tanggal := helpers.SanitizeText(c.PostForm("tanggal"))

	if judul != "" {
		berita.Judul = judul
	}
	if deskripsi != "" {
		berita.Deskripsi = deskripsi
	}
	if tanggal != "" {
		berita.Tanggal = tanggal
	}

	fileHeader, err := c.FormFile("foto")
	if err == nil {
		if fileHeader.Size > 2*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran file maksimal 2MB"})
			return
		}
		if !helpers.IsAllowedFileType(fileHeader.Filename) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ekstensi file tidak diizinkan"})
			return
		}

		if berita.FotoID != "" {
			_, _ = config.Cloudinary.Upload.Destroy(ctx, uploader.DestroyParams{
				PublicID: berita.FotoID,
			})
		}

		src, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuka file"})
			return
		}
		defer src.Close()

		publicID := fmt.Sprintf("berita/%d_%s", time.Now().Unix(), helpers.RandomString(8))

		uploadRes, err := config.Cloudinary.Upload.Upload(ctx, src, uploader.UploadParams{
			PublicID: publicID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal upload gambar ke Cloudinary"})
			return
		}

		berita.Foto = uploadRes.SecureURL
		berita.FotoID = uploadRes.PublicID
	}

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
	ctx := context.Background()
	id := c.Param("id")

	var berita models.Berita
	if err := config.DB.First(&berita, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Berita tidak ditemukan"})
		return
	}

	if berita.FotoID != "" {
		_, _ = config.Cloudinary.Upload.Destroy(ctx, uploader.DestroyParams{
			PublicID: berita.FotoID,
		})
	}

	if err := config.DB.Delete(&berita).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus berita"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Berita berhasil dihapus"})
}
