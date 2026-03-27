package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"llama-go/internal/api"
	"llama-go/internal/backend"
	"llama-go/internal/model"
	"llama-go/internal/session"
)

func main() {
	// 初始化会话管理器
	sessionManager := session.NewSessionManager(100, 3600)

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

	// 注册 OpenAI 兼容路由
	api.RegisterOpenAIRoutes(r, registry)

	// 注册旧 API 路由（兼容性）
	api.RegisterRoutes(r, backend.NewInferenceClient("http://localhost:8000"), sessionManager)

	log.Println("Starting inference server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
