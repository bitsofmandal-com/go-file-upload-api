package handlers

import (
	"net/http"

	"github.com/bitsofmandal-com/go-file-upload-api/internal/service"
	"github.com/gin-gonic/gin"
)


type FileHandler struct {
	fileService *service.FileService
}

func NewFileHandler(fileService *service.FileService) *FileHandler {
	return &FileHandler{fileService: fileService}
}

func (srv *FileHandler) HandleFileUpload(c *gin.Context) {
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
	
	message, httpCode, err := srv.fileService.UploadFile(c.Request.Context(), file, header.Filename)
	if err != nil {
		c.JSON(httpCode, gin.H{"error": err.Error(), "message" : message})
		return
	}
	c.JSON(httpCode, gin.H{
		"message":   "file uploaded successfully at " + message,
	})
}
