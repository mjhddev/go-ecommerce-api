package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mjhddev/go-ecommerce-api/internal/errs"
)

const (
	MaxFileSize = 5 << 20 // 5 MB
)

var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

func SaveProductImage(file *multipart.FileHeader) (string, error) {
	// Validasi ukuran
	if file.Size > MaxFileSize {
		return "", errs.ErrImageTooLarge
	}

	// Validasi ekstensi
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[ext] {
		return "", errs.ErrInvalidImage
	}

	// Pastikan folder ada
	dir := "uploads/products"
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", err
	}

	// Generate nama file unik
	filename := fmt.Sprintf("%d_%s%s",
		time.Now().Unix(),
		uuid.New().String(),
		ext,
	)

	dst := filepath.Join(dir, filename)

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, src); err != nil {
		return "", err
	}

	return "/" + dst, nil
}

func DeleteFile(path string) error {
	if path == "" {
		return nil
	}

	path = strings.TrimPrefix(path, "/")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	return os.Remove(path)
}
