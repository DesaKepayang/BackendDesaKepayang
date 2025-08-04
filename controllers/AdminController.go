package controllers

import (
	"desa-kepayang-backend/config"
	"desa-kepayang-backend/models"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// ==================================
// =========== [CREATE] =============
// ==================================

func TambahAdmin(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required,alphanum,min=3,max=30"`
		Password string `json:"password" binding:"required,min=6,max=100"`
	}

	// Validasi input JSON dan filter karakter khusus
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid", "detail": err.Error()})
		return
	}

	// Cegah SSRF dan abuse dengan sanitasi
	username := strings.TrimSpace(input.Username)
	password := input.Password

	if len(username) == 0 || len(password) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username dan password tidak boleh kosong"})
		return
	}

	// Optional: Cek apakah username sudah ada
	var existing models.Admin
	if err := config.DB.Where("username = ?", username).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username sudah terdaftar"})
		return
	}

	// Hash password dengan bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal meng-hash password"})
		return
	}

	// Simpan ke database secara aman menggunakan GORM (ORM cegah SQL Injection)
	admin := models.Admin{
		Username: username,
		Password: string(hashedPassword),
		Role:     "admin",
	}

	if err := config.DB.Create(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan admin"})
		return
	}

	// Hindari mengembalikan password meskipun sudah di-hash
	c.JSON(http.StatusOK, gin.H{
		"message": "Admin berhasil ditambahkan",
		"data": gin.H{
			"id":       admin.ID,
			"username": admin.Username,
			"role":     admin.Role,
		},
	})
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
	idParam := c.Param("id")

	// Validasi dan sanitasi ID (hindari injection via param ID)
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var input struct {
		Username string `json:"username" binding:"omitempty,alphanum,min=3,max=30"`
		Password string `json:"password" binding:"omitempty,min=6,max=100"`
	}

	// Validasi input JSON dengan aturan
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid", "detail": err.Error()})
		return
	}

	var admin models.Admin
	// GORM parameter binding mencegah SQL Injection
	if err := config.DB.First(&admin, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin tidak ditemukan"})
		return
	}

	// Update username jika disediakan dan aman
	if input.Username != "" {
		admin.Username = strings.TrimSpace(input.Username)
	}

	// Update password jika disediakan
	if input.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal meng-hash password"})
			return
		}
		admin.Password = string(hashedPassword)
	}

	// Simpan perubahan
	if err := config.DB.Save(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate admin"})
		return
	}

	// Hindari mengembalikan field sensitif (seperti hashed password)
	c.JSON(http.StatusOK, gin.H{
		"message": "Admin berhasil diperbarui",
		"data": gin.H{
			"id":       admin.ID,
			"username": admin.Username,
			"role":     admin.Role,
		},
	})
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
		Username string `json:"username" binding:"required,alphanum,min=3,max=30"`
		Password string `json:"password" binding:"required,min=6,max=100"`
	}

	// Validasi input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	username := strings.TrimSpace(input.Username)
	password := input.Password

	if len(username) == 0 || len(password) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username dan password wajib diisi"})
		return
	}

	var admin models.Admin
	if err := config.DB.Where("username = ?", username).First(&admin).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"})
		return
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       admin.ID,
		"username": admin.Username,
		"role":     admin.Role,
		"exp":      time.Now().Add(72 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

	// Set cookie HttpOnly
	secure := false
	if gin.Mode() == gin.ReleaseMode {
		secure = true // hanya diaktifkan di production (HTTPS)
	}

	c.SetCookie("auth_token", tokenString, 72*3600, "/", "", secure, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"token":   tokenString, // ⬅️ tambahkan ini agar frontend bisa pakai token
		"admin": gin.H{
			"id":       admin.ID,
			"username": admin.Username,
			"role":     admin.Role,
		},
	})
}

// =================================
// =========== [LOGOUT] ============
// =================================

func LogoutAdmin(c *gin.Context) {
	// Menghapus cookie dengan mengatur MaxAge menjadi -1
	c.SetCookie("auth_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout berhasil",
	})
}
