package controllers

import (
	"desa-kepayang-backend/config"
	"desa-kepayang-backend/models"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func generateUniqueFileName(originalName string) string {
	timestamp := time.Now().UnixNano()
	random := rand.Intn(9999)
	ext := filepath.Ext(originalName)
	return fmt.Sprintf("%d_%d%s", timestamp, random, ext)
}

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

	// Generate nama file unik
	uniqueFileName := generateUniqueFileName(file.Filename)
	path := "uploads/" + uniqueFileName // Gunakan nama unik

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

	// Ambil semua path foto yang masih digunakan
	var usedPaths []string
	for _, s := range data {
		usedPaths = append(usedPaths, s.Foto)
	}

	// Bersihkan file yang tidak digunakan di folder uploads
	cleanupUnusedFiles("uploads", usedPaths)

	c.JSON(http.StatusOK, data)
}

func cleanupUnusedFiles(folder string, usedPaths []string) {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return // Gagal membaca folder, abaikan
	}

	used := make(map[string]bool)
	for _, path := range usedPaths {
		_, file := filepath.Split(path)
		used[file] = true
	}

	for _, file := range files {
		if !used[file.Name()] {
			os.Remove(filepath.Join(folder, file.Name()))
		}
	}
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
		// Generate nama file unik
		uniqueFileName := generateUniqueFileName(file.Filename)
		newPath := filepath.Join("uploads", uniqueFileName) // Gunakan nama unik

		// Simpan file baru
		if err := c.SaveUploadedFile(file, newPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file baru"})
			return
		}

		// Hapus file lama setelah file baru berhasil disimpan
		if struktur.Foto != "" {
			if err := os.Remove(struktur.Foto); err != nil {
				// Tidak perlu return error, hanya log jika diperlukan
				fmt.Printf("Gagal menghapus file lama: %v\n", err)
			}
		}
		struktur.Foto = newPath
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
