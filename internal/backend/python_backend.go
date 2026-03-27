package backend

import (
	"context"
)

// PythonBackend Python 推理后端适配器
type PythonBackend struct {
	client *InferenceClient
}

// NewPythonBackend 创建 Python 后端
func NewPythonBackend(baseURL string) *PythonBackend {
	return &PythonBackend{
		client: NewInferenceClient(baseURL),
	}
}

// Generate 生成文本（非流式）
func (pb *PythonBackend) Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error) {
	// 转换为 Python 客户端请求格式
	pythonReq := GenerateRequest{
		RequestID:   req.RequestID,
		SessionID:   req.SessionID,
		Model:       req.Model,
		Prompt:      req.Prompt,
		Messages:    req.Messages,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		TopP:        req.TopP,
		TopK:        req.TopK,
		Stream:      false,
		Stop:        req.Stop,
	}

	result, err := pb.client.Generate(ctx, pythonReq)
	if err != nil {
		return nil, err
	}

	// 转换为统一响应格式
	return &GenerateResponse{
		ID:           req.RequestID,
		Model:        req.Model,
		Text:         result.GeneratedText,
		FinishReason: "stop",
		Usage: Usage{
			PromptTokens:     result.PromptTokens,
			CompletionTokens: result.GeneratedTokens,
			TotalTokens:      result.PromptTokens + result.GeneratedTokens,
		},
	}, nil
}

// GenerateStream 生成文本（流式）
func (pb *PythonBackend) GenerateStream(ctx context.Context, req *GenerateRequest) (<-chan StreamChunk, error) {
	// TODO: 实现流式推理
	// 当前返回空 channel，后续实现 SSE 流式处理
	ch := make(chan StreamChunk)
	close(ch)
	return ch, nil
}

// ClearSession 清空会话
func (pb *PythonBackend) ClearSession(ctx context.Context, sessionID string) error {
	return pb.client.ClearSession(ctx, sessionID)
}

// Info 获取后端信息
func (pb *PythonBackend) Info() BackendInfo {
	return BackendInfo{
		Name:           "python",
		SupportsStream: true,
		MaxContextLen:  2048,
		SupportedModels: []string{
			"tinyllama-chat",
		},
	}
}
