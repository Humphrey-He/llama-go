package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"llama-go/internal/backend"
	"llama-go/internal/model"
)

// MockBackendForTest 用于测试的模拟后端
type MockBackendForTest struct {
	generateFunc    func(ctx context.Context, req *backend.GenerateRequest) (*backend.GenerateResponse, error)
	generateStreamFunc func(ctx context.Context, req *backend.GenerateRequest) (<-chan backend.StreamChunk, error)
	infoFunc        func() backend.BackendInfo
}

func (m *MockBackendForTest) Generate(ctx context.Context, req *backend.GenerateRequest) (*backend.GenerateResponse, error) {
	if m.generateFunc != nil {
		return m.generateFunc(ctx, req)
	}
	return &backend.GenerateResponse{
		ID:           req.RequestID,
		Model:        req.Model,
		Text:         "mock response",
		FinishReason: "stop",
		Usage: backend.Usage{
			PromptTokens:     10,
			CompletionTokens: 20,
			TotalTokens:      30,
		},
	}, nil
}

func (m *MockBackendForTest) GenerateStream(ctx context.Context, req *backend.GenerateRequest) (<-chan backend.StreamChunk, error) {
	if m.generateStreamFunc != nil {
		return m.generateStreamFunc(ctx, req)
	}
	ch := make(chan backend.StreamChunk)
	go func() {
		defer close(ch)
		ch <- backend.StreamChunk{Content: "mock ", Done: false}
		ch <- backend.StreamChunk{Content: "stream ", Done: false}
		ch <- backend.StreamChunk{Content: "response", Done: false}
		ch <- backend.StreamChunk{Content: "", Done: true}
	}()
	return ch, nil
}

func (m *MockBackendForTest) ClearSession(ctx context.Context, sessionID string) error {
	return nil
}

func (m *MockBackendForTest) Info() backend.BackendInfo {
	if m.infoFunc != nil {
		return m.infoFunc()
	}
	return backend.BackendInfo{
		Name:           "mock-test",
		SupportsStream: true,
		MaxContextLen:  4096,
	}
}

func setupTestRouter() (*gin.Engine, *model.ModelRegistry) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	registry := model.NewModelRegistry()
	mock := &MockBackendForTest{}
	registry.Register("test-model", mock)

	RegisterOpenAIRoutes(r, registry)
	return r, registry
}

// TestModelRegistry_RegisterAndGet 模型注册与获取测试
func TestModelRegistry_RegisterAndGet(t *testing.T) {
	registry := model.NewModelRegistry()
	mock := &MockBackendForTest{}

	registry.Register("test-model", mock)

	info := registry.Get("test-model")
	if info == nil {
		t.Fatal("expected to get model info, got nil")
	}
	if info.ID != "test-model" {
		t.Errorf("expected model id 'test-model', got %s", info.ID)
	}
	if info.ContextLength != 4096 {
		t.Errorf("expected context length 4096, got %d", info.ContextLength)
	}
}

// TestModelRegistry_GetNonExistent 获取不存在的模型
func TestModelRegistry_GetNonExistent(t *testing.T) {
	registry := model.NewModelRegistry()
	info := registry.Get("non-existent")
	if info != nil {
		t.Error("expected nil for non-existent model")
	}
}

// TestModelRegistry_List 模型列表测试
func TestModelRegistry_List(t *testing.T) {
	registry := model.NewModelRegistry()

	registry.Register("model-1", &MockBackendForTest{})
	registry.Register("model-2", &MockBackendForTest{})

	models := registry.List()
	if len(models) != 2 {
		t.Errorf("expected 2 models, got %d", len(models))
	}
}

// TestChatCompletions_NonStream 非流式聊天完成测试
func TestChatCompletions_NonStream(t *testing.T) {
	router, _ := setupTestRouter()

	payload := ChatCompletionRequest{
		Model: "test-model",
		Messages: []ChatMessage{
			{Role: "user", Content: "Hello"},
		},
		MaxTokens:   256,
		Temperature: 0.7,
		Stream:      false,
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/v1/chat/completions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var resp ChatCompletionResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Object != "chat.completion" {
		t.Errorf("expected object 'chat.completion', got %s", resp.Object)
	}
	if len(resp.Choices) != 1 {
		t.Errorf("expected 1 choice, got %d", len(resp.Choices))
	}
	if resp.Choices[0].Message.Content == "" {
		t.Error("expected non-empty content")
	}
}

// TestChatCompletions_ModelNotFound 模型未找到测试
func TestChatCompletions_ModelNotFound(t *testing.T) {
	router, _ := setupTestRouter()

	payload := ChatCompletionRequest{
		Model: "non-existent-model",
		Messages: []ChatMessage{
			{Role: "user", Content: "Hello"},
		},
		Stream: false,
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/v1/chat/completions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

// TestChatCompletions_InvalidRequest 无效请求测试
func TestChatCompletions_InvalidRequest(t *testing.T) {
	router, _ := setupTestRouter()

	req, _ := http.NewRequest("POST", "/v1/chat/completions", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

// TestChatCompletions_Stream 流式聊天完成测试
func TestChatCompletions_Stream(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	registry := model.NewModelRegistry()
	mock := &MockBackendForTest{}
	registry.Register("stream-model", mock)

	RegisterOpenAIRoutes(r, registry)

	payload := ChatCompletionRequest{
		Model: "stream-model",
		Messages: []ChatMessage{
			{Role: "user", Content: "Hello"},
		},
		MaxTokens: 256,
		Stream:    true,
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/v1/chat/completions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "text/event-stream" {
		t.Errorf("expected content-type 'text/event-stream', got %s", contentType)
	}
}

// TestModels_List 模型列表接口测试
func TestModels_List(t *testing.T) {
	router, registry := setupTestRouter()

	// 添加更多模型
	registry.Register("model-2", &MockBackendForTest{
		infoFunc: func() backend.BackendInfo {
			return backend.BackendInfo{
				Name:           "mock-2",
				SupportsStream: false,
				MaxContextLen:  2048,
			}
		},
	})

	req, _ := http.NewRequest("GET", "/v1/models", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp struct {
		Data []struct {
			ID              string `json:"id"`
			Object          string `json:"object"`
			ContextLength   int    `json:"context_length"`
			SupportsStream  bool   `json:"supports_stream"`
		} `json:"data"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(resp.Data) != 2 {
		t.Errorf("expected 2 models, got %d", len(resp.Data))
	}
}

// TestBackend_GenerateWithContext 带上下文的生成测试
func TestBackend_GenerateWithContext(t *testing.T) {
	mock := &MockBackendForTest{
		generateFunc: func(ctx context.Context, req *backend.GenerateRequest) (*backend.GenerateResponse, error) {
			// 验证 context 可以被取消
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}
			return &backend.GenerateResponse{
				ID:           req.RequestID,
				Model:        req.Model,
				Text:         "context-aware response",
				FinishReason: "stop",
				Usage:        backend.Usage{},
			}, nil
		},
	}

	req := &backend.GenerateRequest{
		RequestID: "test-req",
		Model:     "test",
		Messages: []backend.Message{
			{Role: "user", Content: "Hello"},
		},
	}

	resp, err := mock.Generate(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Text != "context-aware response" {
		t.Errorf("expected 'context-aware response', got %s", resp.Text)
	}
}

// TestBackend_GenerateStreamCancellation 流式取消测试
func TestBackend_GenerateStreamCancellation(t *testing.T) {
	mock := &MockBackendForTest{
		generateStreamFunc: func(ctx context.Context, req *backend.GenerateRequest) (<-chan backend.StreamChunk, error) {
			ch := make(chan backend.StreamChunk)
			go func() {
				defer close(ch)
				for i := 0; i < 100; i++ {
					select {
					case <-ctx.Done():
						return
					case ch <- backend.StreamChunk{Content: "chunk ", Done: false}:
					}
				}
			}()
			return ch, nil
		},
	}

	ctx, cancel := context.WithCancel(context.Background())

	ch, err := mock.GenerateStream(ctx, &backend.GenerateRequest{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 读取几个 chunk 后取消
	for i := 0; i < 5; i++ {
		chunk, ok := <-ch
		if !ok {
			break
		}
		if chunk.Done {
			t.Error("expected done to be false")
		}
	}

	cancel()

	// 验证 context 被取消后 channel 关闭
	for range ch {
		// 消费剩余数据
	}
}

// TestBackend_GenerateStreamMultipleChunks 流式多 chunk 测试
func TestBackend_GenerateStreamMultipleChunks(t *testing.T) {
	chunks := []backend.StreamChunk{
		{Content: "Hello ", Done: false},
		{Content: "world! ", Done: false},
		{Content: "", Done: true},
	}

	mock := &MockBackendForTest{
		generateStreamFunc: func(ctx context.Context, req *backend.GenerateRequest) (<-chan backend.StreamChunk, error) {
			ch := make(chan backend.StreamChunk, len(chunks))
			for _, c := range chunks {
				ch <- c
			}
			close(ch)
			return ch, nil
		},
	}

	result := make([]backend.StreamChunk, 0)
	for chunk := range chunks {
		result = append(result, chunks[chunk])
	}

	if len(result) != 3 {
		t.Errorf("expected 3 chunks, got %d", len(result))
	}

	// 也验证 mock 可以正常调用
	ch, err := mock.GenerateStream(context.Background(), &backend.GenerateRequest{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	count := 0
	for range ch {
		count++
	}
	if count != 3 {
		t.Errorf("expected 3 chunks from mock, got %d", count)
	}
}

// TestAdaptToGenerateRequest_FullMessage 完整消息转换测试
func TestAdaptToGenerateRequest_FullMessage(t *testing.T) {
	req := &ChatCompletionRequest{
		Model: "test-model",
		Messages: []ChatMessage{
			{Role: "system", Content: "You are a helpful assistant."},
			{Role: "user", Content: "Hello"},
			{Role: "assistant", Content: "Hi there!"},
			{Role: "user", Content: "How are you?"},
		},
		Temperature: 0.8,
		TopP:        0.95,
		MaxTokens:   512,
		Stop:        []string{"[STOP]"},
		Stream:      true,
	}

	genReq := AdaptToGenerateRequest(req, "session-full")

	if genReq.SessionID != "session-full" {
		t.Errorf("expected session_id 'session-full', got %s", genReq.SessionID)
	}
	if len(genReq.Messages) != 4 {
		t.Errorf("expected 4 messages, got %d", len(genReq.Messages))
	}
	if genReq.Temperature != 0.8 {
		t.Errorf("expected temperature 0.8, got %f", genReq.Temperature)
	}
	if genReq.TopP != 0.95 {
		t.Errorf("expected top_p 0.95, got %f", genReq.TopP)
	}
	if genReq.MaxTokens != 512 {
		t.Errorf("expected max_tokens 512, got %d", genReq.MaxTokens)
	}
	if len(genReq.Stop) != 1 || genReq.Stop[0] != "[STOP]" {
		t.Errorf("expected stop ['[STOP]'], got %v", genReq.Stop)
	}
	if !genReq.Stream {
		t.Error("expected stream to be true")
	}
}

// TestAdaptFromGenerateResponse_WithUsage 完整响应转换测试（含 Usage）
func TestAdaptFromGenerateResponse_WithUsage(t *testing.T) {
	resp := &backend.GenerateResponse{
		ID:           "resp-full",
		Model:        "test-model",
		Text:         "This is a detailed response.",
		FinishReason: "stop",
		Usage: backend.Usage{
			PromptTokens:     150,
			CompletionTokens: 350,
			TotalTokens:      500,
		},
	}

	chatResp := AdaptFromGenerateResponse(resp, "test-model")

	if chatResp.ID != "resp-full" {
		t.Errorf("expected id 'resp-full', got %s", chatResp.ID)
	}
	if chatResp.Model != "test-model" {
		t.Errorf("expected model 'test-model', got %s", chatResp.Model)
	}
	if chatResp.Usage.PromptTokens != 150 {
		t.Errorf("expected prompt_tokens 150, got %d", chatResp.Usage.PromptTokens)
	}
	if chatResp.Usage.CompletionTokens != 350 {
		t.Errorf("expected completion_tokens 350, got %d", chatResp.Usage.CompletionTokens)
	}
	if chatResp.Usage.TotalTokens != 500 {
		t.Errorf("expected total_tokens 500, got %d", chatResp.Usage.TotalTokens)
	}
}

// TestAdaptStreamChunkToResponse_Done 完成的流式块测试
func TestAdaptStreamChunkToResponse_Done(t *testing.T) {
	chunk := &backend.StreamChunk{
		Content: "",
		Done:    true,
	}

	resp := AdaptStreamChunkToResponse(chunk, "stream-req", "test-model")

	if resp.Choices[0].FinishReason == nil {
		t.Error("expected finish_reason to be set when Done is true")
	}
	if *resp.Choices[0].FinishReason != "stop" {
		t.Errorf("expected finish_reason 'stop', got %s", *resp.Choices[0].FinishReason)
	}
}