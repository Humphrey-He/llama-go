package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	modelPath := flag.String("model", "", "Path to Llama model file")
	flag.Parse()

	if *modelPath == "" {
		log.Fatal("Model path is required")
	}

	fmt.Printf("Loading model from: %s\n", *modelPath)
	// TODO: 实现模型加载和推理逻辑
}
