package utils

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var allowedImageExt = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

var allowedFileExt = map[string]bool{
	".pdf": true,
}

const maxImageSize = 10 << 20 // 10 digeser sebanyak 20 bit, {value} << 20 = {value}mb

func ValidateImage(file *multipart.FileHeader) error {
	// cek ekstensi
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedImageExt[ext] {
		return fmt.Errorf("format file tidak valid, hanya boleh JPG/JPEG/PNG/WEBP")
	}

	// cek ukuran
	if file.Size > maxImageSize {
		return fmt.Errorf("ukuran file terlalu besar, max %dMB", maxImageSize>>20)
	}

	return nil
}

func GenerateFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	return fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
}

func UploadFile(c *gin.Context, fieldName, jenisFile, tipeFile string) (string, error) {
	// ambil file
	file, err := c.FormFile(fieldName)
	if err != nil {
		return "", fmt.Errorf("file %s tidak ditemukan: %v", fieldName, err)
	}

	// validasi ekstensi
	if jenisFile == "image" {
		if err := ValidateImage(file); err != nil {
			return "", err
		}
	}

	// validasi ukuran
	if file.Size > maxImageSize {
		return "", fmt.Errorf("ukuran file terlalu besar, max %d MB", maxImageSize/(1<<20))
	}

	// buat folder path nyata
	basePath := fmt.Sprintf("public/uploads/%s/%s", jenisFile, tipeFile)
	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		return "", fmt.Errorf("gagal membuat folder: %v", err)
	}

	// nama file unik
	newFileName := GenerateFileName(file.Filename)
	savePath := filepath.Join(basePath, newFileName)

	// simpan file
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		return "", fmt.Errorf("gagal menyimpan file: %v", err)
	}

	// return path publik, tanpa mengekspos path nyata
	publicPath := fmt.Sprintf("/static/%s/%s/%s", jenisFile, tipeFile, newFileName)
	return publicPath, nil
}
