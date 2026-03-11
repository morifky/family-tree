package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

// PhotoStorage defines the behavior expected from any photo storage provider
type PhotoStorage interface {
	SavePhoto(fileHeader *multipart.FileHeader) (string, error)
	DeletePhoto(filename string) error
}

type localPhotoStorage struct {
	baseDir string
}

// NewLocalPhotoStorage initializes a new local disk-based photo storage
func NewLocalPhotoStorage(baseDir string) PhotoStorage {
	_ = os.MkdirAll(baseDir, os.ModePerm)
	return &localPhotoStorage{baseDir: baseDir}
}

// SavePhoto takes a multipart file header and saves it to the disk.
// Returns the file path or an error.
func (s *localPhotoStorage) SavePhoto(fileHeader *multipart.FileHeader) (string, error) {
	if fileHeader == nil {
		return "", fmt.Errorf("file is nil")
	}

	ext := filepath.Ext(fileHeader.Filename)
	if ext == "" {
		ext = ".jpg"
	}

	id, err := gonanoid.Generate("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", 16)
	if err != nil {
		return "", fmt.Errorf("could not generate photo id: %w", err)
	}

	filename := id + ext
	dst := filepath.Join(s.baseDir, filename)

	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("could not open uploaded file: %w", err)
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return "", fmt.Errorf("could not create destination file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		return "", fmt.Errorf("could not copy file data: %w", err)
	}

	return filename, nil
}

func (s *localPhotoStorage) DeletePhoto(filename string) error {
	if filename == "" {
		return nil
	}
	path := filepath.Join(s.baseDir, filename)
	err := os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("could not delete photo: %w", err)
	}
	return nil
}
