package controllers

import (
	"desa-kepayang-backend/config"
	"desa-kepayang-backend/models"
	"net/http"
	"os"
	"regexp"
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
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// Bind JSON tanpa validasi bawaan Gin
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal parsing input", "detail": err.Error()})
		return
	}

	// Sanitasi dan ambil nilai
	username := strings.TrimSpace(input.Username)
	password := input.Password

	// Validasi manual: tidak boleh kosong
	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username dan password tidak boleh kosong"})
		return
	}

	// Validasi panjang username
	if len(username) < 3 || len(username) > 30 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username harus 3-30 karakter"})
		return
	}

	// Validasi hanya huruf, angka, dan spasi
	matched, err := regexp.MatchString(`^[a-zA-Z0-9 ]+$`, username)
	if err != nil || !matched {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username hanya boleh huruf, angka, dan spasi"})
		return
	}

	// Validasi panjang password
	if len(password) < 6 || len(password) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password harus 6-100 karakter"})
		return
	}

	// Cek apakah username sudah ada
	var existing models.Admin
	if err := config.DB.Where("username = ?", username).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username sudah terdaftar"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal meng-hash password"})
		return
	}

	// Simpan admin baru ke database
	admin := models.Admin{
		Username: username,
		Password: string(hashedPassword),
		Role:     "admin",
	}

	if err := config.DB.Create(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan admin", "detail": err.Error()})
		return
	}

	// Respon sukses
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

	// Validasi ID
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// Bind JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal parsing input", "detail": err.Error()})
		return
	}

	// Cari admin berdasarkan ID
	var admin models.Admin
	if err := config.DB.First(&admin, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin tidak ditemukan"})
		return
	}

	// Validasi & update username jika ada
	if input.Username != "" {
		username := strings.TrimSpace(input.Username)

		if len(username) < 3 || len(username) > 30 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username harus 3-30 karakter"})
			return
		}

		matched, err := regexp.MatchString(`^[a-zA-Z0-9 ]+$`, username)
		if err != nil || !matched {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username hanya boleh huruf, angka, dan spasi"})
			return
		}

		admin.Username = username
	}

	// Validasi & update password jika ada
	if input.Password != "" {
		if len(input.Password) < 6 || len(input.Password) > 100 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password harus 6-100 karakter"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal meng-hash password"})
			return
		}
		admin.Password = string(hashedPassword)
	}

	// Simpan ke database
	if err := config.DB.Save(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Gagal mengupdate admin",
			"detail": err.Error(),
		})
		return
	}

	// Respon sukses
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
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid", "detail": err.Error()})
		return
	}

	username := strings.TrimSpace(input.Username)
	password := input.Password

	// Validasi manual
	if len(username) < 3 || len(username) > 30 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username harus 3-30 karakter"})
		return
	}

	matched, err := regexp.MatchString(`^[a-zA-Z0-9 ]+$`, username)
	if err != nil || !matched {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username hanya boleh huruf, angka, dan spasi"})
		return
	}

	if len(password) < 6 || len(password) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password harus 6-100 karakter"})
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

	secure := false
	if gin.Mode() == gin.ReleaseMode {
		secure = true
	}

	c.SetCookie("auth_token", tokenString, 72*3600, "/", "", secure, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"token":   tokenString,
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
