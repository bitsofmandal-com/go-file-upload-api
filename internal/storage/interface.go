package storage

import (
	"context"
	"mime/multipart"
)

// FileStoragePort defines the contract for file storage operations.
type FileStoragePort interface {
  Upload(ctx context.Context, file multipart.File, filename string, contentType string) (string, error)
}