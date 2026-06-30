package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	p.baseURL = url
	return nil
}

func (p *Plugin) Close() error                 { return nil }
func (p *Plugin) Name() string                 { return "ollama" }
func (p *Plugin) Version() string              { return "1.0.0" }
func (p *Plugin) Type() string                 { return "ollama" }

func (p *Plugin) Health(ctx context.Context) error {
	req, _ := http.NewRequestWithContext(ctx, "GET", p.baseURL+"/api/tags", nil)
	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (p *Plugin) Generate(ctx context.Context, req *backend.GenerateRequest) (*backend.GenerateResponse, error) {
	messages := make([]map[string]string, 0)
	for _, msg := range req.Messages {
		messages = append(messages, map[string]string{
			"role":    msg.Role,
			"content": msg.Content,
		})
	}

	payload := map[string]interface{}{
		"model":    req.Model,
		"messages": messages,
		"stream":   false,
	}

	body, _ := json.Marshal(payload)
	resp, err := p.client.Post(p.baseURL+"/api/chat", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	content := ""
	if message, ok := result["message"].(map[string]interface{}); ok {
		if c, ok := message["content"].(string); ok {
			content = c
		}
	}

	return &backend.GenerateResponse{
		ID:   req.RequestID,
		Text: content,
	}, nil
}

func (p *Plugin) GenerateStream(ctx context.Context, req *backend.GenerateRequest) (<-chan backend.StreamChunk, error) {
	ch := make(chan backend.StreamChunk)

	go func() {
		defer close(ch)

		messages := make([]map[string]string, 0)
		for _, msg := range req.Messages {
			messages = append(messages, map[string]string{
				"role":    msg.Role,
				"content": msg.Content,
			})
		}

		payload := map[string]interface{}{
			"model":    req.Model,
			"messages": messages,
			"stream":   true,
		}

		body, _ := json.Marshal(payload)
		resp, err := p.client.Post(p.baseURL+"/api/chat", "application/json", bytes.NewReader(body))
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

			if message, ok := chunk["message"].(map[string]interface{}); ok {
				if text, ok := message["content"].(string); ok && text != "" {
					ch <- backend.StreamChunk{Content: text}
				}
			}

			if done, ok := chunk["done"].(bool); ok && done {
				ch <- backend.StreamChunk{Done: true}
				break
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
		Name:           "ollama",
		SupportsStream: true,
		MaxContextLen:  8192,
	}
}
