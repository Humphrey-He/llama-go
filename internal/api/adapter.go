package api

import (
	"llama-go/internal/backend"
	"time"
)

// ChatCompletionRequest OpenAI 兼容的聊天完成请求
type ChatCompletionRequest struct {
	Model       string                   `json:"model"`
	Messages    []ChatMessage            `json:"messages"`
	Temperature float64                  `json:"temperature"`
	TopP        float64                  `json:"top_p"`
	MaxTokens   int                      `json:"max_tokens"`
	Stream      bool                     `json:"stream"`
	Stop        []string                 `json:"stop"`
}

// ChatMessage 聊天消息
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionResponse OpenAI 兼容的聊天完成响应
type ChatCompletionResponse struct {
	ID      string                   `json:"id"`
	Object  string                   `json:"object"`
	Created int64                    `json:"created"`
	Model   string                   `json:"model"`
	Choices []ChatCompletionChoice   `json:"choices"`
	Usage   ChatCompletionUsage      `json:"usage"`
}

// ChatCompletionChoice 聊天完成选择
type ChatCompletionChoice struct {
	Index        int          `json:"index"`
	Message      ChatMessage  `json:"message"`
	FinishReason string       `json:"finish_reason"`
}

// ChatCompletionUsage 使用统计
type ChatCompletionUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ChatCompletionStreamResponse 流式响应
type ChatCompletionStreamResponse struct {
	ID      string                        `json:"id"`
	Object  string                        `json:"object"`
	Created int64                         `json:"created"`
	Model   string                        `json:"model"`
	Choices []ChatCompletionStreamChoice  `json:"choices"`
}

// ChatCompletionStreamChoice 流式选择
type ChatCompletionStreamChoice struct {
	Index        int                    `json:"index"`
	Delta        ChatCompletionDelta    `json:"delta"`
	FinishReason *string                `json:"finish_reason"`
}

// ChatCompletionDelta 增量内容
type ChatCompletionDelta struct {
	Content string `json:"content"`
}

// AdaptToGenerateRequest 将 OpenAI 请求转换为内部格式
func AdaptToGenerateRequest(req *ChatCompletionRequest, sessionID string) *backend.GenerateRequest {
	messages := make([]backend.Message, len(req.Messages))
	for i, msg := range req.Messages {
		messages[i] = backend.Message{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	return &backend.GenerateRequest{
		SessionID:   sessionID,
		Model:       req.Model,
		Messages:    messages,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		TopP:        req.TopP,
		Stream:      req.Stream,
		Stop:        req.Stop,
	}
}

// AdaptFromGenerateResponse 将内部响应转换为 OpenAI 格式
func AdaptFromGenerateResponse(resp *backend.GenerateResponse, model string) *ChatCompletionResponse {
	return &ChatCompletionResponse{
		ID:      resp.ID,
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   model,
		Choices: []ChatCompletionChoice{
			{
				Index: 0,
				Message: ChatMessage{
					Role:    "assistant",
					Content: resp.Text,
				},
				FinishReason: resp.FinishReason,
			},
		},
		Usage: ChatCompletionUsage{
			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		},
	}
}

// AdaptStreamChunkToResponse 将流式块转换为 OpenAI 格式
func AdaptStreamChunkToResponse(chunk *backend.StreamChunk, id, model string) *ChatCompletionStreamResponse {
	finishReason := (*string)(nil)
	if chunk.Done {
		reason := "stop"
		finishReason = &reason
	}

	return &ChatCompletionStreamResponse{
		ID:      id,
		Object:  "chat.completion.chunk",
		Created: time.Now().Unix(),
		Model:   model,
		Choices: []ChatCompletionStreamChoice{
			{
				Index: 0,
				Delta: ChatCompletionDelta{
					Content: chunk.Content,
				},
				FinishReason: finishReason,
			},
		},
	}
}
