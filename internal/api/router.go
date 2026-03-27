package api

import (
	"github.com/gin-gonic/gin"
	"llama-go/internal/service"
)

// RegisterRoutes 注册路由
func RegisterRoutes(r *gin.Engine, cs *service.ChatService) {
	handler := NewHandler(cs)

	// 健康检查
	r.GET("/healthz", handler.Health)
	r.GET("/readyz", handler.Health)

	// API 路由
	v1 := r.Group("/v1")
	{
		v1.POST("/chat/completions", handler.ChatCompletions)
		v1.POST("/sessions/:id/clear", handler.ClearSession)
	}
}
