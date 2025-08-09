package controllers

import (
	"context"
	"desa-kepayang-backend/config"
	"desa-kepayang-backend/helpers"
	"desa-kepayang-backend/models"
	"fmt"
	"net/http"
	"time"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

// ==================================
// ========== [CREATE] ==============
// ==================================
func CreateStrukturDesa(c *gin.Context) {
	ctx := context.Background()

	nama := helpers.SanitizeText(c.PostForm("nama"))
	jabatan := helpers.SanitizeText(c.PostForm("jabatan"))

	if nama == "" || jabatan == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama dan jabatan wajib diisi"})
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

	publicID := fmt.Sprintf("struktur_desa/%d_%s", time.Now().Unix(), helpers.RandomString(8))

	uploadRes, err := config.Cloudinary.Upload.Upload(ctx, src, uploader.UploadParams{
		PublicID: publicID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal upload gambar ke Cloudinary"})
		return
	}

	struktur := models.StrukturDesa{
		Nama:    nama,
		Jabatan: jabatan,
		Foto:    uploadRes.SecureURL,
		FotoID:  uploadRes.PublicID,
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
	ctx := context.Background()
	id := c.Param("id")

	var struktur models.StrukturDesa
	if err := config.DB.First(&struktur, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Struktur desa tidak ditemukan"})
		return
	}

	nama := helpers.SanitizeText(c.PostForm("nama"))
	jabatan := helpers.SanitizeText(c.PostForm("jabatan"))

	if nama != "" {
		struktur.Nama = nama
	}
	if jabatan != "" {
		struktur.Jabatan = jabatan
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

		if struktur.FotoID != "" {
			_, _ = config.Cloudinary.Upload.Destroy(ctx, uploader.DestroyParams{
				PublicID: struktur.FotoID,
			})
		}

		src, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuka file"})
			return
		}
		defer src.Close()

		publicID := fmt.Sprintf("struktur_desa/%d_%s", time.Now().Unix(), helpers.RandomString(8))

		uploadRes, err := config.Cloudinary.Upload.Upload(ctx, src, uploader.UploadParams{
			PublicID: publicID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal upload gambar ke Cloudinary"})
			return
		}

		struktur.Foto = uploadRes.SecureURL
		struktur.FotoID = uploadRes.PublicID
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
	ctx := context.Background()
	id := c.Param("id")

	var struktur models.StrukturDesa
	if err := config.DB.First(&struktur, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Struktur desa tidak ditemukan"})
		return
	}

	if struktur.FotoID != "" {
		_, _ = config.Cloudinary.Upload.Destroy(ctx, uploader.DestroyParams{
			PublicID: struktur.FotoID,
		})
	}

	if err := config.DB.Delete(&struktur).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus struktur desa"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Struktur desa berhasil dihapus"})
}
