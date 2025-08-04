package helpers

import (
	"fmt"
	"math/rand"
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
