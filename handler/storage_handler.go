package handler

import (
	"context"
	"net/http"

	"go-do-spaces-poc/config"
	"go-do-spaces-poc/storage"

	"github.com/gin-gonic/gin"
)

type StorageHandler struct {
	cfg    *config.Config
	client *storage.ClientWrapper
}

func NewStorageHandler(cfg *config.Config) *StorageHandler {
	return &StorageHandler{
		cfg: cfg,
	}
}

// POST /upload
func (h *StorageHandler) UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file"})
		return
	}
	defer file.Close()

	client := storage.NewSpacesClient()
	url, err := storage.UploadObject(c, client, h.cfg.DOBucket, "uploads/"+header.Filename, file, header.Size, header.Header.Get("Content-Type"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}

// DELETE /delete/:key
func (h *StorageHandler) DeleteFile(c *gin.Context) {
	key := c.Param("key")
	client := storage.NewSpacesClient()

	if err := storage.DeleteObject(c, client, h.cfg.DOBucket, key); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted successfully"})
}

// GET /list
func (h *StorageHandler) ListFiles(c *gin.Context) {
	client := storage.NewSpacesClient()
	files, err := storage.ListObjects(context.Background(), client, h.cfg.DOBucket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"files": files})
}
