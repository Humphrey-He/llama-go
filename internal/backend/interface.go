package backend

import "context"

// Message 消息结构
type Message struct {
	Role    string `json:"role"`    // system/user/assistant
	Content string `json:"content"`
}

// GenerateRequest 统一生成请求
type GenerateRequest struct {
	RequestID   string    `json:"request_id"`
	SessionID   string    `json:"session_id"`
	Model       string    `json:"model"`
	Prompt      string    `json:"prompt"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
	TopP        float64   `json:"top_p"`
	TopK        int       `json:"top_k"`
	Stream      bool      `json:"stream"`
	Stop        []string  `json:"stop"`
}

// GenerateResponse 统一生成响应
type GenerateResponse struct {
	ID           string `json:"id"`
	Model        string `json:"model"`
	Text         string `json:"text"`
	FinishReason string `json:"finish_reason"`
	Usage        Usage  `json:"usage"`
}

// StreamChunk 流式数据块
type StreamChunk struct {
	Content string `json:"content"`
	Done    bool   `json:"done"`
}

// Usage token 使用统计
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// BackendInfo 后端信息
type BackendInfo struct {
	Name            string
	SupportsStream  bool
	MaxContextLen   int
	SupportedModels []string
}

// InferenceBackend 推理后端接口
type InferenceBackend interface {
	Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error)
	GenerateStream(ctx context.Context, req *GenerateRequest) (<-chan StreamChunk, error)
	ClearSession(ctx context.Context, sessionID string) error
	Info() BackendInfo
}
