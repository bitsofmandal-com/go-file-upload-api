package main

import (
	"log"

	"github.com/bitsofmandal-com/go-file-upload-api/internal/handlers"
	"github.com/bitsofmandal-com/go-file-upload-api/internal/service"
	"github.com/bitsofmandal-com/go-file-upload-api/internal/storage"
	"github.com/gin-gonic/gin"
)

func main() {
	// initialize storage adapter
	localStorage := storage.NewLocalStorageAdapter("./uploads") // replace with your desired local path
	// s3Storage := storage.NewS3StorageAdapter("your-bucket-name", &s3.Client{}) // replace with actual S3 client
	// initialize file service
	fileService := service.NewFileService(localStorage) // or s3Storage
	fileHandler := handlers.NewFileHandler(fileService) // or s3Storage

	router := gin.Default()

	// NOOB Mistake 1: Not setting the max upload size
	// Max upload size: 10MB
	router.MaxMultipartMemory = 10 << 20 // 10 MiB
	router.POST("/upload", fileHandler.HandleFileUpload)
	// router.GET("/uploads/:filename", service.UploadFile)

	log.Println("Server running at http://localhost:8080")
	router.Run(":8080")
}
