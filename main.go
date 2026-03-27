package main

import (
	"github.com/gin-gonic/gin"
	//"llama-go/api"
	//"llama-go/internal/kvcache"
)

func main() {
	r := gin.Default()
	println("Go服务运行 : 8080 ")
	err := r.Run(":8080")
	if err != nil {
		return
	}
}
