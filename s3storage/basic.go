package s3storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
    MaxUploadSize = 10 << 20 // 10 MiB
    S3Bucket      = "your-s3-bucket-name" // Replace with your S3 bucket name
    S3Region      = "your-region"         // Replace with your AWS region, e.g., "us-west-2"
)

var allowedMIMEs = map[string]struct{}{
    "image/jpeg":      {},
    "image/png":       {},
    "application/pdf": {},
    "text/plain":      {},
}

func HandleUpload(c *gin.Context) {
    file, header, err := c.Request.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "file field is required"})
        return
    }
    defer file.Close()

    if header.Size == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "uploaded file is empty"})
        return
    }

    // Detect MIME type
    mimeType, err := detectMIME(file)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to detect file type"})
        return
    }

    if _, ok := allowedMIMEs[mimeType]; !ok {
        c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported file type", "type": mimeType})
        return
    }

    // Reset file pointer
    if _, err := file.Seek(0, io.SeekStart); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reset file reader"})
        return
    }

    // Generate unique filename
    ext := strings.ToLower(filepath.Ext(header.Filename))
    newFilename := fmt.Sprintf("%s%s", uuid.NewString(), ext)

    // Upload to S3
    location, err := uploadToS3(file, newFilename, mimeType)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload to S3", "details": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message":  "file uploaded successfully",
        "filename": newFilename,
        "mime":     mimeType,
        "url":      location,
    })
}

func detectMIME(file multipart.File) (string, error) {
    buffer := make([]byte, 512)
    if _, err := file.Read(buffer); err != nil && err != io.EOF {
        return "", err
    }
    return http.DetectContentType(buffer), nil
}

func uploadToS3(file multipart.File, key, contentType string) (string, error) {
    ctx := context.Background()

    cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(S3Region))
    if err != nil {
        return "", fmt.Errorf("failed to load AWS config: %w", err)
    }

    client := s3.NewFromConfig(cfg)
    uploader := manager.NewUploader(client)

    result, err := uploader.Upload(ctx, &s3.PutObjectInput{
        Bucket:      aws.String(S3Bucket),
        Key:         aws.String(key),
        Body:        file,
        ContentType: aws.String(contentType),
    })
    if err != nil {
        return "", fmt.Errorf("failed to upload file to S3: %w", err)
    }

    return result.Location, nil
}
