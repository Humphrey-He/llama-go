package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
	var req backend.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.SessionID == "" {
		req.SessionID = "default"
	}

	if req.Stream {
		h.streamChat(c, req)
	} else {
		h.chat(c, req)
	}
}

// chat 非流式聊天
func (h *Handler) chat(c *gin.Context, req backend.ChatRequest) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	resp, err := h.chatService.Chat(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":      fmt.Sprintf("req-%d", time.Now().UnixNano()),
		"object":  "chat.completion",
		"choices": []gin.H{{
			"message": gin.H{
				"role":    "assistant",
				"content": resp.Content,
			},
			"finish_reason": "stop",
		}},
		"usage": resp.Usage,
	})
}

// streamChat 流式聊天
func (h *Handler) streamChat(c *gin.Context, req backend.ChatRequest) {
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
		if chunk.Done {
			c.SSEvent("message", "[DONE]")
		} else {
			c.SSEvent("message", gin.H{"delta": chunk.Delta})
		}
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
