package handler

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"go-do-spaces-poc/config"
	"go-do-spaces-poc/storage"

	"github.com/gin-gonic/gin"
)

type StorageHandler struct {
	cfg *config.Config
}

func NewStorageHandler(cfg *config.Config) *StorageHandler {
	return &StorageHandler{cfg: cfg}
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

	// (optional) custom folder path, defaults to "uploads/"
	path := c.DefaultPostForm("path", "uploads/")
	key := path + header.Filename

	url, err := storage.UploadObject(
		c, client, h.cfg.DOBucket, key,
		file, header.Size, header.Header.Get("Content-Type"),
	)
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

// POST /lifecycle/set
func (h *StorageHandler) SetLifecyclePolicy(c *gin.Context) {
	var req struct {
		Prefix         string `json:"prefix" binding:"required"`
		ExpirationDays int32  `json:"expiration_days" binding:"required,min=1"`
		RuleID         string `json:"rule_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate a default rule ID if not provided
	if req.RuleID == "" {
		req.RuleID = "delete-" + req.Prefix + "-rule"
	}

	client := storage.NewSpacesClient()
	err := storage.SetLifecyclePolicy(c, client, h.cfg.DOBucket, req.Prefix, req.ExpirationDays, req.RuleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "lifecycle policy set successfully",
		"rule_id":         req.RuleID,
		"prefix":          req.Prefix,
		"expiration_days": req.ExpirationDays,
	})
}

// GET /lifecycle/list
func (h *StorageHandler) GetLifecyclePolicies(c *gin.Context) {
	client := storage.NewSpacesClient()
	rules, err := storage.GetLifecyclePolicy(c, client, h.cfg.DOBucket)
	if err != nil {
		// If no lifecycle configuration exists, return empty list instead of error
		if storage.IsNoLifecycleConfigError(err) {
			c.JSON(http.StatusOK, gin.H{"rules": []gin.H{}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Format rules for better readability
	var formattedRules []gin.H
	for _, rule := range rules {
		ruleData := gin.H{
			"id":     rule.ID,
			"status": rule.Status,
		}

		// Extract prefix
		if rule.Prefix != nil {
			ruleData["prefix"] = *rule.Prefix
		}

		// Extract expiration days
		if rule.Expiration != nil && rule.Expiration.Days != nil {
			ruleData["expiration_days"] = *rule.Expiration.Days
		}

		formattedRules = append(formattedRules, ruleData)
	}

	c.JSON(http.StatusOK, gin.H{"rules": formattedRules})
}

// DELETE /lifecycle/delete/:ruleId
func (h *StorageHandler) DeleteLifecyclePolicy(c *gin.Context) {
	ruleID := c.Param("ruleId")
	if ruleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "rule_id is required"})
		return
	}

	client := storage.NewSpacesClient()
	err := storage.DeleteLifecyclePolicy(c, client, h.cfg.DOBucket, ruleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "lifecycle policy deleted successfully", "rule_id": ruleID})
}

// GET /download/folder?prefix=<folder-path>
func (h *StorageHandler) DownloadFolder(c *gin.Context) {
	prefix := c.Query("prefix")
	if prefix == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "prefix query parameter is required"})
		return
	}

	// Ensure prefix ends with / if not already
	if !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}

	client := storage.NewSpacesClient()

	// List all objects with the given prefix
	keys, err := storage.ListObjectsByPrefix(c, client, h.cfg.DOBucket, prefix)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(keys) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no files found with the given prefix"})
		return
	}

	// Create a buffer to write our archive to
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// Download each file and add it to the zip
	for _, key := range keys {
		// Skip directory markers (keys ending with /)
		if strings.HasSuffix(key, "/") {
			continue
		}

		// Download the file
		data, err := storage.DownloadObject(c, client, h.cfg.DOBucket, key)
		if err != nil {
			zipWriter.Close()
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to download %s: %v", key, err)})
			return
		}

		// Create a file in the zip archive
		// Use the relative path (remove the prefix) for better organization
		relativePath := strings.TrimPrefix(key, prefix)
		if relativePath == "" {
			relativePath = filepath.Base(key)
		}

		zipFile, err := zipWriter.Create(relativePath)
		if err != nil {
			zipWriter.Close()
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create zip entry for %s: %v", key, err)})
			return
		}

		_, err = zipFile.Write(data)
		if err != nil {
			zipWriter.Close()
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to write %s to zip: %v", key, err)})
			return
		}
	}

	// Close the zip writer to finalize the archive
	err = zipWriter.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to finalize zip archive"})
		return
	}

	// Generate a filename for the zip
	zipFilename := strings.TrimSuffix(strings.TrimSuffix(prefix, "/"), "/")
	if zipFilename == "" {
		zipFilename = "download"
	} else {
		// Use only the last part of the path as filename
		zipFilename = filepath.Base(zipFilename)
	}
	zipFilename = zipFilename + ".zip"

	// Set headers and send the zip file
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", zipFilename))
	c.Data(http.StatusOK, "application/zip", buf.Bytes())
}
