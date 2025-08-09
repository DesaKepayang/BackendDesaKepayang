package helpers

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Fungsi untuk membuat nama file unik dan aman
func GenerateUniqueFileName(originalName string) string {
	timestamp := time.Now().UnixNano()
	random := rand.Intn(9999)
	ext := filepath.Ext(originalName)
	safeExt := strings.ToLower(ext)
	if safeExt != ".jpg" && safeExt != ".jpeg" && safeExt != ".png" && safeExt != ".webp" {
		return ""
	}
	return fmt.Sprintf("%d_%d%s", timestamp, random, safeExt)
}

// Fungsi untuk menyaring karakter tidak aman
func SanitizeText(input string) string {
	return strings.TrimSpace(strings.ReplaceAll(input, "<", ""))
}

// Fungsi untuk membersihkan file tidak terpakai [BARU]
func CleanupUnusedFiles(folder string, usedPaths []string) {
	files, err := os.ReadDir(folder)
	if err != nil {
		return
	}

	used := make(map[string]bool)
	for _, path := range usedPaths {
		_, file := filepath.Split(path)
		used[file] = true
	}

	for _, file := range files {
		if !file.IsDir() && !used[file.Name()] {
			os.Remove(filepath.Join(folder, file.Name()))
		}
	}
}

func IsAllowedFileType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	allowedExt := []string{".jpg", ".jpeg", ".png", ".webp"}

	for _, allowed := range allowedExt {
		if ext == allowed {
			return true
		}
	}
	return false
}

// Fungsi untuk membuat string acak
func RandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
