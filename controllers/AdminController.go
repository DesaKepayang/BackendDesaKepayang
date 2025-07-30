package controllers

import (
	"desa-kepayang-backend/config"
	"desa-kepayang-backend/models"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// ==================================
// =========== [CREATE] =============
// ==================================

func TambahAdmin(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// Validasi input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid", "detail": err.Error()})
		return
	}

	// Hash password sebelum disimpan
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghash password", "detail": err.Error()})
		return
	}

	// Tetapkan role admin secara default
	admin := models.Admin{
		Username: input.Username,
		Password: string(hashedPassword),
		Role:     "admin",
	}

	// Simpan ke database
	if err := config.DB.Create(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan admin", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Admin berhasil ditambahkan", "data": admin})
}

// =================================
// =========== [READ] ==============
// =================================

func GetAllAdmin(c *gin.Context) {
	var admins []models.Admin

	if err := config.DB.Find(&admins).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data admin", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": admins})
}

func GetAdminProfile(c *gin.Context) {
	// Ambil data dari context
	claims, exists := c.Get("adminData")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Data admin tidak ditemukan di token"})
		return
	}

	mapClaims := claims.(jwt.MapClaims)
	id := uint(mapClaims["id"].(float64)) // konversi dari float64 ke uint

	var admin models.Admin
	if err := config.DB.First(&admin, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": admin})
}

// ==================================
// =========== [UPDATE] =============
// ==================================

func UpdateAdmin(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid", "detail": err.Error()})
		return
	}

	var admin models.Admin
	if err := config.DB.First(&admin, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin tidak ditemukan"})
		return
	}

	// Update username jika disediakan
	if input.Username != "" {
		admin.Username = input.Username
	}

	// Hash dan update password jika disediakan
	if input.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghash password", "detail": err.Error()})
			return
		}
		admin.Password = string(hashedPassword)
	}

	if err := config.DB.Save(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate admin", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Admin berhasil diperbarui", "data": admin})
}

// ==================================
// =========== [DELETE] =============
// ==================================

func DeleteAdmin(c *gin.Context) {
	id := c.Param("id")

	var admin models.Admin
	if err := config.DB.First(&admin, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin tidak ditemukan"})
		return
	}

	if err := config.DB.Delete(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus admin", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Admin berhasil dihapus"})
}

// =================================
// =========== [LOGIN] =============
// =================================

func LoginAdmin(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	var admin models.Admin
	if err := config.DB.Where("username = ?", input.Username).First(&admin).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username tidak ditemukan"})
		return
	}

	// Verifikasi password
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password salah"})
		return
	}

	// Generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       admin.ID,
		"username": admin.Username,
		"role":     admin.Role,
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
