package main

import (
	"fmt"
	"llama-go/pkg/llama"
)

func main() {
	// 创建推理配置
	config := &llama.InferenceConfig{
		ContextSize:   2048,
		MaxTokens:     512,
		Temperature:   0.7,
		TopP:          0.9,
		TopK:          40,
		RepeatPenalty: 1.1,
	}

	// 创建推理引擎
	inferencer := llama.NewInferencer(config)

	// 执行推理
	prompt := "Hello, how are you?"
	result, err := inferencer.Infer(prompt)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Prompt: %s\n", prompt)
	fmt.Printf("Result: %s\n", result)
}
