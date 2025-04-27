package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/bitsofmandal-com/go-file-upload-api/internal/storage"
	"github.com/google/uuid"
)

// Use a map to define allowed MIME types for better performance
// and to avoid using a switch statement
var allowedMIMEs = map[string]struct{}{
	"image/jpeg":      {},
	"image/png":       {},
	"application/pdf": {},
	"text/plain":      {},
}

type FileService struct {
  Storage storage.FileStoragePort
}

func NewFileService(storage storage.FileStoragePort) *FileService {
  return &FileService{Storage: storage}
}

func (f *FileService) UploadFile(ctx context.Context, file multipart.File, origionalFileName string) (string, int, error) {

  // Allow only selected content types
	ok, contentType, err := isAllowedMimeType(file)
	if !ok {
    return fmt.Sprintf("unsupported type: %s", contentType), http.StatusBadRequest, err
	}

	// Reset file pointer back
	if _, err := file.Seek(0, 0); err != nil {
    return "failed to reset file", http.StatusInternalServerError, err
	}

	// Generate safe filename
	safeName := generateSafeFilename(origionalFileName)

	// Upload using FileService (port)
  fileUploadPath, err :=  f.Storage.Upload(ctx, file, safeName, contentType)
  return fileUploadPath, http.StatusOK, err
}

// detectMIME reads a small buffer to determine the file's MIME type
func isAllowedMimeType(f multipart.File) (bool, string, error) {
	buffer := make([]byte, 512)

	if _, err := f.Read(buffer); err != nil && err != io.EOF {
		return false, "unknown", err
	}

	// Reset the file pointer back to start after reading
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return false, "unknown", err
	}

	fileType := http.DetectContentType(buffer)
	if _, ok := allowedMIMEs[fileType]; !ok {
		return false, fileType, errors.New("unsupported file type")
	}
	return true, fileType, nil
}

func generateSafeFilename(filename string) string {
	newFilename := filepath.Base(filename)                  // basic sanitization
	newFilename = filepath.Clean(newFilename)               // clean up the filename
	newFilename = strings.ReplaceAll(newFilename, " ", "_") // replace spaces with underscores
	return fmt.Sprintf("%s-%s", uuid.NewString(), newFilename)
}