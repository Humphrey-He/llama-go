package backend

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const PythonBackendURL = "http://localhost:8000/infer"

// InferRequest 推理请求
type InferRequest struct {
	Prompt  string      `json:"prompt"`
	KVCache interface{} `json:"kv_cache,omitempty"`
}

// InferResponse 推理响应
type InferResponse struct {
	Text     string      `json:"text"`
	KV       interface{} `json:"kv"`
	TokenNum int         `json:"token_num"`
}

// CallPythonBackend 调用 Python 推理服务
func CallPythonBackend(prompt string, kv interface{}) (*InferResponse, error) {
	reqBody := InferRequest{
		Prompt:  prompt,
		KVCache: kv,
	}
	jsonData, _ := json.Marshal(reqBody)

	resp, err := http.Post(PythonBackendURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result InferResponse
	json.NewDecoder(resp.Body).Decode(&result)
	return &result, nil
}
