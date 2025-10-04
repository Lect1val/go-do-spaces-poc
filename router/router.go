package router

import (
	"go-do-spaces-poc/config"
	"go-do-spaces-poc/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	h := handler.NewStorageHandler(cfg)
	r.POST("/upload", h.UploadFile)
	r.DELETE("/delete/:key", h.DeleteFile)
	r.GET("/list", h.ListFiles)

	return r
}
