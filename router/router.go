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

	// Lifecycle policy endpoints
	r.POST("/lifecycle/set", h.SetLifecyclePolicy)
	r.GET("/lifecycle/list", h.GetLifecyclePolicies)
	r.DELETE("/lifecycle/delete/:ruleId", h.DeleteLifecyclePolicy)

	return r
}
