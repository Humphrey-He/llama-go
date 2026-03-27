package backend

import "context"

// Message 消息结构
type Message struct {
	Role    string `json:"role"`    // system/user/assistant
	Content string `json:"content"`
}

// ChatRequest 聊天请求
type ChatRequest struct {
	Model       string     `json:"model"`
	Messages    []Message  `json:"messages"`
	Temperature float32    `json:"temperature"`
	TopP        float32    `json:"top_p"`
	MaxTokens   int        `json:"max_tokens"`
	Stream      bool       `json:"stream"`
	SessionID   string     `json:"session_id"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Usage   Usage  `json:"usage"`
}

// StreamChunk 流式数据块
type StreamChunk struct {
	Delta string `json:"delta"`
	Done  bool   `json:"done"`
}

// Usage token 使用统计
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// LLMBackend LLM 后端接口
type LLMBackend interface {
	Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error)
	StreamChat(ctx context.Context, req ChatRequest) (<-chan StreamChunk, error)
	Health(ctx context.Context) error
}
