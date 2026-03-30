package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"llama-go/internal/api"
	"llama-go/internal/backend"
	"llama-go/internal/model"
)

func main() {
	// 初始化模型注册表
	registry := model.NewModelRegistry()

	// 注册 Python 后端
	pythonBackend := backend.NewPythonBackend("http://localhost:8000")
	registry.Register("tinyllama-chat", pythonBackend)

	// 注册模拟后端（用于测试）
	mockBackend := backend.NewMockBackend()
	registry.Register("mock-model", mockBackend)

	// 初始化 Gin
	r := gin.Default()

	// 健康检查端点
	r.GET("/healthz", func(c *gin.Context) {
		// Liveness: 仅检查进程是否存活
		c.JSON(http.StatusOK, gin.H{"status": "alive"})
	})

	r.GET("/readyz", func(c *gin.Context) {
		// Readiness: 检查是否能提供完整服务（包括后端连接）
		_, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		// 检查后端连接
		_ = pythonBackend.Info()

		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	// 注册 OpenAI 兼容路由
	api.RegisterOpenAIRoutes(r, registry)

	log.Println("Starting inference server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
