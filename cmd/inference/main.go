package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"llama-go/internal/api"
	"llama-go/internal/backend"
	"llama-go/internal/session"
)

func main() {
	// 初始化推理客户端
	client := backend.NewInferenceClient("http://localhost:8000")

	// 初始化会话管理器
	sessionManager := session.NewSessionManager(100, 3600)

	// 初始化 Gin
	r := gin.Default()

	// 注册路由
	api.RegisterRoutes(r, client, sessionManager)

	log.Println("Starting inference server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
