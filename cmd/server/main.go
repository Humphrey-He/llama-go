package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"llama-go/internal/api"
	"llama-go/internal/backend"
	"llama-go/internal/service"
	"llama-go/internal/session"
)

func main() {
	// 初始化后端（vLLM）
	vllmBackend := backend.NewVLLMBackend("http://localhost:8000")

	// 初始化会话存储
	sessionStore := session.NewSessionStore()

	// 初始化服务
	chatService := service.NewChatService(vllmBackend, sessionStore)

	// 初始化 Gin
	r := gin.Default()

	// 注册路由
	api.RegisterRoutes(r, chatService)

	log.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
