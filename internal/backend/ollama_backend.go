package backend

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

// OllamaBackend Ollama 后端适配器
type OllamaBackend struct {
	baseURL string
	client  *http.Client
}

// NewOllamaBackend 创建 Ollama 后端
func NewOllamaBackend(baseURL string) *OllamaBackend {
	return &OllamaBackend{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

// Generate 生成响应
func (ob *OllamaBackend) Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error) {
	// 转换为 Ollama chat 格式
	messages := make([]map[string]string, 0)
	for _, msg := range req.Messages {
		messages = append(messages, map[string]string{
			"role":    msg.Role,
			"content": msg.Content,
		})
	}

	ollamaReq := map[string]interface{}{
		"model":    req.Model,
		"messages": messages,
		"stream":   false,
	}

	body, _ := json.Marshal(ollamaReq)
	resp, err := ob.client.Post(ob.baseURL+"/api/chat", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	// 提取响应内容
	content := ""
	if message, ok := result["message"].(map[string]interface{}); ok {
		if c, ok := message["content"].(string); ok {
			content = c
		}
	}

	return &GenerateResponse{
		ID:   req.RequestID,
		Text: content,
	}, nil
}

// GenerateStream 流式生成
func (ob *OllamaBackend) GenerateStream(ctx context.Context, req *GenerateRequest) (<-chan StreamChunk, error) {
	chunkCh := make(chan StreamChunk)

	go func() {
		defer close(chunkCh)

		ollamaReq := map[string]interface{}{
			"model":  req.Model,
			"prompt": req.Prompt,
			"stream": true,
		}

		body, _ := json.Marshal(ollamaReq)
		resp, err := ob.client.Post(ob.baseURL+"/api/generate", "application/json", bytes.NewReader(body))
		if err != nil {
			return
		}
		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)
		for {
			var chunk map[string]interface{}
			if err := decoder.Decode(&chunk); err == io.EOF {
				break
			} else if err != nil {
				return
			}

			if text, ok := chunk["response"].(string); ok && text != "" {
				chunkCh <- StreamChunk{
					Content: text,
					Done:    false,
				}
			}

			if done, ok := chunk["done"].(bool); ok && done {
				chunkCh <- StreamChunk{Done: true}
				break
			}
		}
	}()

	return chunkCh, nil
}

// ClearSession 清除会话
func (ob *OllamaBackend) ClearSession(ctx context.Context, sessionID string) error {
	return nil
}

// Info 返回后端信息
func (ob *OllamaBackend) Info() BackendInfo {
	return BackendInfo{
		Name:            "ollama",
		MaxContextLen:   8192,
		SupportsStream:  true,
		SupportedModels: []string{"llama3.1:8b"},
	}
}

