package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"llama-go/internal/backend"
	"llama-go/internal/session"
)

// GenerateHandler 生成处理器
type GenerateHandler struct {
	client         *backend.InferenceClient
	sessionManager *session.SessionManager
}

// NewGenerateHandler 创建生成处理器
func NewGenerateHandler(client *backend.InferenceClient, sm *session.SessionManager) *GenerateHandler {
	return &GenerateHandler{
		client:         client,
		sessionManager: sm,
	}
}

// GenerateRequest 生成请求
type GenerateRequest struct {
	SessionID    string  `json:"session_id"`
	Prompt       string  `json:"prompt" binding:"required"`
	MaxNewTokens int     `json:"max_new_tokens"`
	Temperature  float64 `json:"temperature"`
	TopP         float64 `json:"top_p"`
	TopK         int     `json:"top_k"`
	Stream       bool    `json:"stream"`
}

// Generate 生成文本
func (gh *GenerateHandler) Generate(c *gin.Context) {
	var req GenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置默认值
	if req.SessionID == "" {
		req.SessionID = fmt.Sprintf("session-%d", time.Now().UnixNano())
	}
	if req.MaxNewTokens == 0 {
		req.MaxNewTokens = 128
	}
	if req.Temperature == 0 {
		req.Temperature = 0.7
	}
	if req.TopP == 0 {
		req.TopP = 0.9
	}
	if req.TopK == 0 {
		req.TopK = 50
	}

	// 获取或创建会话
	sessionMeta := gh.sessionManager.GetSession(req.SessionID)
	if sessionMeta == nil {
		sessionMeta = gh.sessionManager.CreateSession(req.SessionID)
	}

	// 调用推理
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	backendReq := backend.GenerateRequest{
		SessionID:    req.SessionID,
		Prompt:       req.Prompt,
		MaxNewTokens: req.MaxNewTokens,
		Temperature:  req.Temperature,
		TopP:         req.TopP,
		TopK:         req.TopK,
		Stream:       req.Stream,
	}

	result, err := gh.client.Generate(ctx, backendReq)
	if err != nil {
		log.Printf("Generate error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 更新会话
	gh.sessionManager.UpdateSession(req.SessionID, result.PromptTokens+result.GeneratedTokens)

	// 记录日志
	log.Printf("Generate: session_id=%s, mode=%s, cache_hit=%v, ttft_ms=%.2f",
		req.SessionID, result.Mode, result.CacheHit, result.TTFTMS)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"session_id":       result.SessionID,
			"generated_text":   result.GeneratedText,
			"prompt_tokens":    result.PromptTokens,
			"generated_tokens": result.GeneratedTokens,
			"cache_hit":        result.CacheHit,
			"mode":             result.Mode,
			"ttft_ms":          result.TTFTMS,
			"total_latency_ms": result.TotalLatencyMS,
		},
	})
}

// ClearSession 清空会话
func (gh *GenerateHandler) ClearSession(c *gin.Context) {
	sessionID := c.Param("id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id required"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := gh.client.ClearSession(ctx, sessionID); err != nil {
		log.Printf("Clear session error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	gh.sessionManager.DeleteSession(sessionID)
	log.Printf("Session cleared: %s", sessionID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "session cleared",
	})
}

// Health 健康检查
func (gh *GenerateHandler) Health(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	if err := gh.client.Health(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Metrics 指标
func (gh *GenerateHandler) Metrics(c *gin.Context) {
	stats := gh.sessionManager.GetStats()
	c.JSON(http.StatusOK, stats)
}
