package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

const UploadDir = "./uploads"

func main() {
	// Ensure upload directory exists
	if err := os.MkdirAll(UploadDir, os.ModePerm); err != nil {
		log.Fatalf("failed to create upload dir: %v", err)
	}

	router := gin.Default()

	// Max upload size: 10MB
	router.MaxMultipartMemory = 10 << 20 // 10 MiB

	router.POST("/upload", handleUpload)

	log.Println("Server running at http://localhost:8080")
	router.Run(":8080")
}

func handleUpload(c *gin.Context) {
	// Multipart form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file not provided"})
		return
	}
	defer file.Close()

	// Save to uploads/
	filename := filepath.Base(header.Filename) // basic sanitization
	outPath := filepath.Join(UploadDir, filename)

	if err := c.SaveUploadedFile(header, outPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to save file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "file uploaded successfully",
		"filepath": outPath,
	})
}