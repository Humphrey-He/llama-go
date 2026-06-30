package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"llama-go/internal/api"
	"llama-go/internal/config"
	"llama-go/internal/model"
	"llama-go/internal/plugin"
	"llama-go/internal/plugin/llamacpp"
	"llama-go/internal/plugin/ollama"
)

func main() {
	// 加载配置
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Load config: %v", err)
	}

	// 初始化插件管理器
	pm := plugin.NewManager()

	// 注册可用插件
	pluginFactories := map[string]func() plugin.Plugin{
		"llama-cpp": llamacpp.New,
		"ollama":    ollama.New,
	}

	// 初始化启用的插件
	ctx := context.Background()
	for name, factory := range pluginFactories {
		pluginCfg, ok := cfg.Plugins[name]
		if !ok || !pluginCfg.Enabled {
			continue
		}

		p := factory()
		config := map[string]interface{}{
			"base_url": pluginCfg.BaseURL,
		}

		if err := p.Init(ctx, config); err != nil {
			log.Printf("Init plugin %s failed: %v", name, err)
			continue
		}

		if err := pm.Register(name, p); err != nil {
			log.Fatalf("Register plugin %s: %v", name, err)
		}
		log.Printf("Plugin loaded: %s v%s", p.Name(), p.Version())
	}

	// 初始化模型注册表
	registry := model.NewModelRegistry()

	// 注册模型到插件
	for pluginName, pluginCfg := range cfg.Plugins {
		if !pluginCfg.Enabled {
			continue
		}

		p, err := pm.Get(pluginName)
		if err != nil {
			continue
		}

		for _, modelID := range pluginCfg.Models {
			registry.Register(modelID, p)
			log.Printf("Model registered: %s -> %s", modelID, pluginName)
		}
	}

	// 初始化 Gin
	r := gin.Default()

	// 健康检查
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "alive"})
	})

	r.GET("/readyz", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		for _, name := range pm.List() {
			p, _ := pm.Get(name)
			if err := p.Health(ctx); err != nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{
					"status": "unhealthy",
					"plugin": name,
					"error":  err.Error(),
				})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	// 注册 OpenAI 兼容路由
	api.RegisterOpenAIRoutes(r, registry)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting inference server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
