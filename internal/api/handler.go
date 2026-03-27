package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"llama-go/internal/backend"
	"llama-go/internal/service"
)

// Handler API 处理器
type Handler struct {
	chatService *service.ChatService
}

// NewHandler 创建 API 处理器
func NewHandler(cs *service.ChatService) *Handler {
	return &Handler{
		chatService: cs,
	}
}

// ChatCompletions 聊天完成接口
func (h *Handler) ChatCompletions(c *gin.Context) {
	var req ChatCompletionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sessionID := uuid.New().String()
	genReq := AdaptToGenerateRequest(&req, sessionID)
	genReq.RequestID = uuid.New().String()

	if req.Stream {
		h.streamChat(c, genReq)
	} else {
		h.chat(c, genReq)
	}
}

// chat 非流式聊天
func (h *Handler) chat(c *gin.Context, req *backend.GenerateRequest) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	resp, err := h.chatService.Chat(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, AdaptFromGenerateResponse(resp, req.Model))
}

// streamChat 流式聊天
func (h *Handler) streamChat(c *gin.Context, req *backend.GenerateRequest) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	ch, err := h.chatService.StreamChat(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	for chunk := range ch {
		resp := AdaptStreamChunkToResponse(&chunk, req.RequestID, req.Model)
		c.SSEvent("message", resp)
	}
}

// ClearSession 清空会话
func (h *Handler) ClearSession(c *gin.Context) {
	sessionID := c.Param("id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id required"})
		return
	}

	if err := h.chatService.ClearSession(sessionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "session cleared"})
}

// Health 健康检查
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
