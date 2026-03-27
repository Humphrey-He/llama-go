package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"llama-go/internal/backend"
	"llama-go/internal/model"
	"llama-go/internal/stream"
)

// OpenAIRoutes OpenAI 兼容路由处理器
type OpenAIRoutes struct {
	registry *model.ModelRegistry
}

// NewOpenAIRoutes 创建 OpenAI 路由处理器
func NewOpenAIRoutes(registry *model.ModelRegistry) *OpenAIRoutes {
	return &OpenAIRoutes{
		registry: registry,
	}
}

// ChatCompletions 处理聊天完成请求
func (or *OpenAIRoutes) ChatCompletions(c *gin.Context) {
	var req ChatCompletionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	modelInfo := or.registry.Get(req.Model)
	if modelInfo == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "model not found"})
		return
	}

	sessionID := uuid.New().String()
	genReq := AdaptToGenerateRequest(&req, sessionID)
	genReq.RequestID = uuid.New().String()

	if req.Stream {
		or.handleStreamChat(c, modelInfo, genReq)
	} else {
		or.handleNonStreamChat(c, modelInfo, genReq)
	}
}

// handleNonStreamChat 处理非流式聊天
func (or *OpenAIRoutes) handleNonStreamChat(c *gin.Context, modelInfo *model.ModelInfo, req *backend.GenerateRequest) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	resp, err := modelInfo.Backend.Generate(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, AdaptFromGenerateResponse(resp, req.Model))
}

// handleStreamChat 处理流式聊天
func (or *OpenAIRoutes) handleStreamChat(c *gin.Context, modelInfo *model.ModelInfo, req *backend.GenerateRequest) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	chunks, err := modelInfo.Backend.GenerateStream(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "streaming not supported"})
		return
	}

	sseWriter := stream.NewSSEWriter(c.Writer)

	for chunk := range chunks {
		resp := AdaptStreamChunkToResponse(&chunk, req.RequestID, req.Model)
		if err := sseWriter.WriteChunk(resp); err != nil {
			return
		}
		sseWriter.Flush(flusher)
	}

	sseWriter.WriteDone()
	flusher.Flush()
}

// Models 获取模型列表
func (or *OpenAIRoutes) Models(c *gin.Context) {
	models := or.registry.List()
	data := make([]gin.H, len(models))

	for i, m := range models {
		data[i] = gin.H{
			"id":              m.ID,
			"object":          "model",
			"owned_by":        "llama-go",
			"context_length":  m.ContextLength,
			"supports_stream": m.SupportsStream,
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

// RegisterOpenAIRoutes 注册 OpenAI 兼容路由
func RegisterOpenAIRoutes(r *gin.Engine, registry *model.ModelRegistry) {
	routes := NewOpenAIRoutes(registry)

	v1 := r.Group("/v1")
	{
		v1.POST("/chat/completions", routes.ChatCompletions)
		v1.GET("/models", routes.Models)
	}
}
