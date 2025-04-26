package main

import (
	"log"
	"os"

	"github.com/bitsofmandal-com/go-file-upload-api/localstorage"
	"github.com/gin-gonic/gin"
)

func main() {
	// Ensure upload directory exists
	if err := os.MkdirAll(localstorage.UploadDir, os.ModePerm); err != nil {
		log.Fatalf("failed to create upload dir: %v", err)
	}
	router := gin.Default()

	// NOOB Mistake 1: Not setting the max upload size
	// Max upload size: 10MB
	router.MaxMultipartMemory = 10 << 20 // 10 MiB

	router.POST("/upload", localstorage.HandleUpload)
	router.GET("/uploads/:filename", localstorage.HandleGetFile)

	log.Println("Server running at http://localhost:8080")
	router.Run(":8080")
}
