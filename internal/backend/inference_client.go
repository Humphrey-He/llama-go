package backend

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// InferenceClient Python 推理客户端
type InferenceClient struct {
	baseURL string
	client  *http.Client
}

// GenerateRequest 生成请求
type GenerateRequest struct {
	SessionID    string  `json:"session_id"`
	Prompt       string  `json:"prompt"`
	MaxNewTokens int     `json:"max_new_tokens"`
	Temperature  float64 `json:"temperature"`
	TopP         float64 `json:"top_p"`
	TopK         int     `json:"top_k"`
	Stream       bool    `json:"stream"`
}

// GenerateResponse 生成响应
type GenerateResponse struct {
	Success bool   `json:"success"`
	Data    Result `json:"data"`
}

// Result 推理结果
type Result struct {
	SessionID       string  `json:"session_id"`
	GeneratedText   string  `json:"generated_text"`
	PromptTokens    int     `json:"prompt_tokens"`
	GeneratedTokens int     `json:"generated_tokens"`
	CacheHit        bool    `json:"cache_hit"`
	Mode            string  `json:"mode"`
	TTFTMS          float64 `json:"ttft_ms"`
	TotalLatencyMS  float64 `json:"total_latency_ms"`
}

// NewInferenceClient 创建推理客户端
func NewInferenceClient(baseURL string) *InferenceClient {
	return &InferenceClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Generate 生成文本
func (ic *InferenceClient) Generate(ctx context.Context, req GenerateRequest) (*Result, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST",
		ic.baseURL+"/api/generate", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := ic.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("inference failed: %d", resp.StatusCode)
	}

	var result GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if !result.Success {
		return nil, fmt.Errorf("inference error")
	}

	return &result.Data, nil
}

// ClearSession 清空会话
func (ic *InferenceClient) ClearSession(ctx context.Context, sessionID string) error {
	httpReq, err := http.NewRequestWithContext(ctx, "POST",
		ic.baseURL+"/api/sessions/"+sessionID+"/clear", nil)
	if err != nil {
		return err
	}

	resp, err := ic.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("clear session failed: %d", resp.StatusCode)
	}

	return nil
}

// Health 健康检查
func (ic *InferenceClient) Health(ctx context.Context) error {
	httpReq, err := http.NewRequestWithContext(ctx, "GET",
		ic.baseURL+"/healthz", nil)
	if err != nil {
		return err
	}

	resp, err := ic.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed: %d", resp.StatusCode)
	}

	return nil
}
