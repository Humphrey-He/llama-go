package api

import (
	"testing"

	"llama-go/internal/backend"
)

func TestAdaptToGenerateRequest(t *testing.T) {
	req := &ChatCompletionRequest{
		Model: "tinyllama-chat",
		Messages: []ChatMessage{
			{Role: "user", Content: "Hello"},
		},
		Temperature: 0.7,
		TopP:        0.9,
		MaxTokens:   256,
		Stream:      false,
	}

	genReq := AdaptToGenerateRequest(req, "session-123")

	if genReq.Model != "tinyllama-chat" {
		t.Errorf("expected model tinyllama-chat, got %s", genReq.Model)
	}
	if genReq.SessionID != "session-123" {
		t.Errorf("expected session_id session-123, got %s", genReq.SessionID)
	}
	if len(genReq.Messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(genReq.Messages))
	}
	if genReq.Messages[0].Content != "Hello" {
		t.Errorf("expected content Hello, got %s", genReq.Messages[0].Content)
	}
}

func TestAdaptFromGenerateResponse(t *testing.T) {
	resp := &backend.GenerateResponse{
		ID:           "resp-123",
		Model:        "tinyllama-chat",
		Text:         "Hello there",
		FinishReason: "stop",
		Usage: backend.Usage{
			PromptTokens:     10,
			CompletionTokens: 5,
			TotalTokens:      15,
		},
	}

	chatResp := AdaptFromGenerateResponse(resp, "tinyllama-chat")

	if chatResp.ID != "resp-123" {
		t.Errorf("expected id resp-123, got %s", chatResp.ID)
	}
	if chatResp.Object != "chat.completion" {
		t.Errorf("expected object chat.completion, got %s", chatResp.Object)
	}
	if len(chatResp.Choices) != 1 {
		t.Errorf("expected 1 choice, got %d", len(chatResp.Choices))
	}
	if chatResp.Choices[0].Message.Content != "Hello there" {
		t.Errorf("expected content Hello there, got %s", chatResp.Choices[0].Message.Content)
	}
	if chatResp.Usage.TotalTokens != 15 {
		t.Errorf("expected total_tokens 15, got %d", chatResp.Usage.TotalTokens)
	}
}

func TestAdaptStreamChunkToResponse(t *testing.T) {
	chunk := &backend.StreamChunk{
		Content: "Hello",
		Done:    false,
	}

	resp := AdaptStreamChunkToResponse(chunk, "id-123", "tinyllama-chat")

	if resp.Object != "chat.completion.chunk" {
		t.Errorf("expected object chat.completion.chunk, got %s", resp.Object)
	}
	if resp.Choices[0].Delta.Content != "Hello" {
		t.Errorf("expected content Hello, got %s", resp.Choices[0].Delta.Content)
	}
	if resp.Choices[0].FinishReason != nil {
		t.Errorf("expected finish_reason nil, got %v", resp.Choices[0].FinishReason)
	}
}
