package storage

import (
	"context"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
)

type LocalStorageAdapter struct {
  BasePath string
}

func NewLocalStorageAdapter(basePath string) *LocalStorageAdapter {
  	// Ensure upload directory exists
	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		log.Fatalf("failed to create upload dir: %v", err)
	}
  return &LocalStorageAdapter{BasePath: basePath}
}

func (l *LocalStorageAdapter) Upload(ctx context.Context, file multipart.File, filename string, contentType string) (string, error) {
  path := filepath.Join(l.BasePath, filename)
  outFile, err := os.Create(path)
  if err != nil {
      return "", err
  }
  defer outFile.Close()

  _, err = io.Copy(outFile, file)
  if err != nil {
      return "", err
  }

  return path, nil
}
