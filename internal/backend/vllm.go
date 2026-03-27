package backend

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// VLLMBackend vLLM 后端实现
type VLLMBackend struct {
	baseURL string
	client  *http.Client
}

// NewVLLMBackend 创建 vLLM 后端
func NewVLLMBackend(baseURL string) *VLLMBackend {
	return &VLLMBackend{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		client:  &http.Client{},
	}
}

// Generate 非流式生成
func (v *VLLMBackend) Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error) {
	payload := map[string]interface{}{
		"model":       req.Model,
		"messages":    req.Messages,
		"temperature": req.Temperature,
		"top_p":       req.TopP,
		"max_tokens":  req.MaxTokens,
		"stream":      false,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", v.baseURL+"/v1/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := v.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	choices := result["choices"].([]interface{})
	message := choices[0].(map[string]interface{})["message"].(map[string]interface{})
	content := message["content"].(string)

	usage := result["usage"].(map[string]interface{})
	return &GenerateResponse{
		ID:           result["id"].(string),
		Model:        req.Model,
		Text:         content,
		FinishReason: "stop",
		Usage: Usage{
			PromptTokens:     int(usage["prompt_tokens"].(float64)),
			CompletionTokens: int(usage["completion_tokens"].(float64)),
			TotalTokens:      int(usage["total_tokens"].(float64)),
		},
	}, nil
}

// GenerateStream 流式生成
func (v *VLLMBackend) GenerateStream(ctx context.Context, req *GenerateRequest) (<-chan StreamChunk, error) {
	ch := make(chan StreamChunk)

	payload := map[string]interface{}{
		"model":       req.Model,
		"messages":    req.Messages,
		"temperature": req.Temperature,
		"top_p":       req.TopP,
		"max_tokens":  req.MaxTokens,
		"stream":      true,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		close(ch)
		return ch, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", v.baseURL+"/v1/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		close(ch)
		return ch, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	go func() {
		defer close(ch)

		resp, err := v.client.Do(httpReq)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadString('\n')
			if err != nil && err != io.EOF {
				return
			}

			line = strings.TrimSpace(line)
			if line == "" || line == "[DONE]" {
				if line == "[DONE]" {
					ch <- StreamChunk{Done: true}
				}
				if err == io.EOF {
					return
				}
				continue
			}

			if strings.HasPrefix(line, "data: ") {
				line = strings.TrimPrefix(line, "data: ")
				var data map[string]interface{}
				if err := json.Unmarshal([]byte(line), &data); err != nil {
					continue
				}

				choices := data["choices"].([]interface{})
				if len(choices) > 0 {
					delta := choices[0].(map[string]interface{})["delta"].(map[string]interface{})
					if content, ok := delta["content"].(string); ok {
						ch <- StreamChunk{Content: content}
					}
				}
			}

			if err == io.EOF {
				return
			}
		}
	}()

	return ch, nil
}

// Health 健康检查
func (v *VLLMBackend) Health(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", v.baseURL+"/health", nil)
	if err != nil {
		return err
	}

	resp, err := v.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed: %d", resp.StatusCode)
	}

	return nil
}
