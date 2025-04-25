package localstorage

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

const UploadDir = "./uploads"

func HandleUpload(c *gin.Context) {
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

func HandleGetFile(c *gin.Context) {
  filename := c.Param("filename")
  filePath := filepath.Join(UploadDir, filename)

  // Serve the file
  c.File(filePath)
}