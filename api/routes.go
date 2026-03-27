package api

import (
	"github.com/gin-gonic/gin"
	"llama-go/internal/backend"
	"llama-go/internal/kvcache"
	"net/http"
)

func RegisterRoutes(r *gin.Engine, cache *kvcache.KVCache) {
	api := r.Group("/api")
	{
		// 文本生成（核心接口）
		api.POST("/generate", func(c *gin.Context) {
			var req struct {
				Prompt string `json:"prompt" binding:"required"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// 1. 调用 Python 推理
			resp, err := backend.CallPythonBackend(req.Prompt, nil)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// 2. 写入 KV 缓存
			sessionID := "default"
			cache.Set(sessionID, &kvcache.CacheEntry{
				Keys:  resp.KV.(map[string]interface{})["keys"].([][]float32),
				Vals:  resp.KV.(map[string]interface{})["vals"].([][]float32),
				Token: resp.TokenNum,
			})

			// 3. 返回结果 + 缓存状态
			c.JSON(http.StatusOK, gin.H{
				"generated_text": resp.Text,
				"cache_size":     cache.GetCacheSize(sessionID),
				"status":         "success",
			})
		})

		// 清空缓存
		api.POST("/clear", func(c *gin.Context) {
			cache.Clear("default")
			c.JSON(http.StatusOK, gin.H{"message": "cache cleared"})
		})
	}
}
