package api

import (
	"github.com/gin-gonic/gin"
	"llama-go/internal/backend"
	"llama-go/internal/session"
)

// RegisterRoutes 注册路由
func RegisterRoutes(r *gin.Engine, client *backend.InferenceClient, sm *session.SessionManager) {
	handler := NewGenerateHandler(client, sm)

	// 健康检查
	r.GET("/healthz", handler.Health)
	r.GET("/readyz", handler.Health)

	// API 路由
	api := r.Group("/api")
	{
		api.POST("/generate", handler.Generate)
		api.POST("/sessions/:id/clear", handler.ClearSession)
	}

	// 指标
	r.GET("/metrics", handler.Metrics)
}
