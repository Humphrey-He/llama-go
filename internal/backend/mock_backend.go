package backend

import (
	"context"
)

// MockBackend 模拟后端用于测试
type MockBackend struct {
	responses map[string]string
}

// NewMockBackend 创建模拟后端
func NewMockBackend() *MockBackend {
	return &MockBackend{
		responses: map[string]string{
			"default": "This is a mock response from the inference backend.",
		},
	}
}

// Generate 生成文本（非流式）
func (mb *MockBackend) Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error) {
	response, ok := mb.responses[req.Model]
	if !ok {
		response = mb.responses["default"]
	}

	return &GenerateResponse{
		ID:           req.RequestID,
		Model:        req.Model,
		Text:         response,
		FinishReason: "stop",
		Usage: Usage{
			PromptTokens:     10,
			CompletionTokens: 20,
			TotalTokens:      30,
		},
	}, nil
}

// GenerateStream 生成文本（流式）
func (mb *MockBackend) GenerateStream(ctx context.Context, req *GenerateRequest) (<-chan StreamChunk, error) {
	ch := make(chan StreamChunk, 10)

	go func() {
		defer close(ch)
		response, ok := mb.responses[req.Model]
		if !ok {
			response = mb.responses["default"]
		}

		for _, char := range response {
			select {
			case <-ctx.Done():
				return
			case ch <- StreamChunk{
				Content: string(char),
				Done:    false,
			}:
			}
		}

		ch <- StreamChunk{
			Content: "",
			Done:    true,
		}
	}()

	return ch, nil
}

// ClearSession 清空会话
func (mb *MockBackend) ClearSession(ctx context.Context, sessionID string) error {
	return nil
}

// Info 获取后端信息
func (mb *MockBackend) Info() BackendInfo {
	return BackendInfo{
		Name:            "mock",
		SupportsStream:  true,
		MaxContextLen:   2048,
		SupportedModels: []string{"mock-model"},
	}
}

// SetResponse 设置模拟响应
func (mb *MockBackend) SetResponse(model, response string) {
	mb.responses[model] = response
}
