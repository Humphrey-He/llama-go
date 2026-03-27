package backend

import (
	"testing"
)

func TestInferRequestStructure(t *testing.T) {
	req := InferRequest{
		Prompt:  "test prompt",
		KVCache: nil,
	}

	if req.Prompt != "test prompt" {
		t.Errorf("expected 'test prompt', got '%s'", req.Prompt)
	}
}

func TestInferResponseStructure(t *testing.T) {
	resp := InferResponse{
		Text:     "generated text",
		KV:       map[string]interface{}{"keys": []interface{}{}, "vals": []interface{}{}},
		TokenNum: 5,
	}

	if resp.TokenNum != 5 {
		t.Errorf("expected token_num 5, got %d", resp.TokenNum)
	}

	if resp.Text != "generated text" {
		t.Errorf("expected 'generated text', got '%s'", resp.Text)
	}
}
