package main

import (
	"github.com/gin-gonic/gin"
	"llama-go/api"
	"llama-go/internal/config"
	"llama-go/internal/kvcache"
	"llama-go/internal/logger"
)

func main() {
	r := gin.Default()

	// 初始化 KV Cache
	cache := kvcache.NewKVCache()
	logger.InfoLogger.Println("KV Cache initialized")

	// 注册路由
	api.RegisterRoutes(r, cache)
	logger.InfoLogger.Println("API routes registered")

	logger.InfoLogger.Printf("Starting Go server on %s", config.GoServerPort)
	if err := r.Run(config.GoServerPort); err != nil {
		logger.ErrorLogger.Fatalf("Failed to start server: %v", err)
	}
}
