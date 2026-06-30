package vllm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"llama-go/internal/backend"
	"llama-go/internal/plugin"
)

type Plugin struct {
	baseURL string
	client  *http.Client
}

func New() plugin.Plugin {
	return &Plugin{
		client: &http.Client{},
	}
}

func (p *Plugin) Init(ctx context.Context, config map[string]interface{}) error {
	url, ok := config["base_url"].(string)
	if !ok {
		return fmt.Errorf("missing base_url")
	}
	p.baseURL = strings.TrimSuffix(url, "/")
	return nil
}

func (p *Plugin) Close() error                 { return nil }
func (p *Plugin) Name() string                 { return "vllm" }
func (p *Plugin) Version() string              { return "1.0.0" }
func (p *Plugin) Type() string                 { return "vllm" }

func (p *Plugin) Health(ctx context.Context) error {
	req, _ := http.NewRequestWithContext(ctx, "GET", p.baseURL+"/health", nil)
	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("unhealthy: %d", resp.StatusCode)
	}
	return nil
}

func (p *Plugin) Generate(ctx context.Context, req *backend.GenerateRequest) (*backend.GenerateResponse, error) {
	payload := map[string]interface{}{
		"model":       req.Model,
		"messages":    req.Messages,
		"temperature": req.Temperature,
		"max_tokens":  req.MaxTokens,
		"stream":      false,
	}

	body, _ := json.Marshal(payload)
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/v1/chat/completions", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	choices := result["choices"].([]interface{})
	message := choices[0].(map[string]interface{})["message"].(map[string]interface{})
	content := message["content"].(string)

	return &backend.GenerateResponse{
		ID:    result["id"].(string),
		Model: req.Model,
		Text:  content,
	}, nil
}

func (p *Plugin) GenerateStream(ctx context.Context, req *backend.GenerateRequest) (<-chan backend.StreamChunk, error) {
	ch := make(chan backend.StreamChunk)

	payload := map[string]interface{}{
		"model":    req.Model,
		"messages": req.Messages,
		"stream":   true,
	}

	body, _ := json.Marshal(payload)
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/v1/chat/completions", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	go func() {
		defer close(ch)

		resp, err := p.client.Do(httpReq)
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
					ch <- backend.StreamChunk{Done: true}
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
						ch <- backend.StreamChunk{Content: content}
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

func (p *Plugin) ClearSession(ctx context.Context, sessionID string) error {
	return nil
}

func (p *Plugin) Info() backend.BackendInfo {
	return backend.BackendInfo{
		Name:           "vllm",
		SupportsStream: true,
		MaxContextLen:  4096,
	}
}
