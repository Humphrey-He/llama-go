package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"llama-go/internal/backend"
	"llama-go/internal/kvcache"
	"net/http"
)

type GenerateRequest struct {
	Prompt    string `json:"prompt" binding:"required"`
	SessionID string `json:"session_id"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func RegisterRoutes(r *gin.Engine, cache *kvcache.KVCache) {
	api := r.Group("/api")
	{
		// 文本生成（核心接口）
		api.POST("/generate", func(c *gin.Context) {
			var req GenerateRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, APIResponse{
					Success: false,
					Error:   fmt.Sprintf("invalid request: %v", err),
				})
				return
			}

			// 默认会话ID
			sessionID := req.SessionID
			if sessionID == "" {
				sessionID = "default"
			}

			// 1. 获取历史 KV Cache（复用优化）
			cachedKV := cache.Get(sessionID)
			var kvCache interface{}
			if len(cachedKV) > 0 {
				// 将缓存转换为可传递的格式
				kvCache = convertCacheToKV(cachedKV)
			}

			// 2. 调用 Python 推理（传递历史缓存）
			resp, err := backend.CallPythonBackend(req.Prompt, kvCache)
			if err != nil {
				c.JSON(http.StatusInternalServerError, APIResponse{
					Success: false,
					Error:   fmt.Sprintf("inference failed: %v", err),
				})
				return
			}

			// 3. 安全提取 KV 数据
			keys, vals, err := extractKVData(resp.KV)
			if err != nil {
				c.JSON(http.StatusInternalServerError, APIResponse{
					Success: false,
					Error:   fmt.Sprintf("extract kv failed: %v", err),
				})
				return
			}

			// 4. 写入 KV 缓存
			cache.Set(sessionID, &kvcache.CacheEntry{
				Keys:  keys,
				Vals:  vals,
				Token: resp.TokenNum,
			})

			// 5. 返回结果
			c.JSON(http.StatusOK, APIResponse{
				Success: true,
				Data: gin.H{
					"generated_text": resp.Text,
					"cache_size":     cache.GetCacheSize(sessionID),
					"session_id":     sessionID,
					"cache_reused":   len(cachedKV) > 0,
				},
			})
		})

		// 清空缓存
		api.POST("/clear", func(c *gin.Context) {
			var req struct {
				SessionID string `json:"session_id"`
			}
			c.ShouldBindJSON(&req)

			sessionID := req.SessionID
			if sessionID == "" {
				sessionID = "default"
			}

			cache.Clear(sessionID)
			c.JSON(http.StatusOK, APIResponse{
				Success: true,
				Data:    gin.H{"message": "cache cleared", "session_id": sessionID},
			})
		})
	}
}

// extractKVData 安全提取 KV 数据
func extractKVData(kv map[string]interface{}) ([][]float32, [][]float32, error) {
	keysRaw, ok := kv["keys"]
	if !ok {
		return nil, nil, fmt.Errorf("keys not found in kv")
	}

	valsRaw, ok := kv["vals"]
	if !ok {
		return nil, nil, fmt.Errorf("vals not found in kv")
	}

	keys, err := convertToFloat32Array(keysRaw)
	if err != nil {
		return nil, nil, fmt.Errorf("convert keys failed: %w", err)
	}

	vals, err := convertToFloat32Array(valsRaw)
	if err != nil {
		return nil, nil, fmt.Errorf("convert vals failed: %w", err)
	}

	return keys, vals, nil
}

// convertToFloat32Array 转换为 [][]float32
func convertToFloat32Array(data interface{}) ([][]float32, error) {
	arr, ok := data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("data is not array")
	}

	result := make([][]float32, len(arr))
	for i, row := range arr {
		rowArr, ok := row.([]interface{})
		if !ok {
			return nil, fmt.Errorf("row %d is not array", i)
		}

		result[i] = make([]float32, len(rowArr))
		for j, val := range rowArr {
			switch v := val.(type) {
			case float64:
				result[i][j] = float32(v)
			case float32:
				result[i][j] = v
			default:
				return nil, fmt.Errorf("invalid type at [%d][%d]", i, j)
			}
		}
	}

	return result, nil
}

// convertCacheToKV 将缓存转换为可传递的 KV 格式
func convertCacheToKV(entries []*kvcache.CacheEntry) map[string]interface{} {
	if len(entries) == 0 {
		return nil
	}

	// 取最后一个缓存条目
	lastEntry := entries[len(entries)-1]
	return map[string]interface{}{
		"keys": lastEntry.Keys,
		"vals": lastEntry.Vals,
	}
}
