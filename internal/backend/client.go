package backend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"llama-go/internal/config"
	"net/http"
)

// InferRequest 推理请求
type InferRequest struct {
	Prompt  string      `json:"prompt"`
	KVCache interface{} `json:"kv_cache,omitempty"`
}

// InferResponse 推理响应
type InferResponse struct {
	Text     string                 `json:"text"`
	KV       map[string]interface{} `json:"kv"`
	TokenNum int                    `json:"token_num"`
}

// CallPythonBackend 调用 Python 推理服务
func CallPythonBackend(prompt string, kv interface{}) (*InferResponse, error) {
	reqBody := InferRequest{
		Prompt:  prompt,
		KVCache: kv,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %w", err)
	}

	resp, err := http.Post(config.PythonBackendURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("python backend returned status %d", resp.StatusCode)
	}

	var result InferResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	return &result, nil
}
