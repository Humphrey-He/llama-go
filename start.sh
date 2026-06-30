#!/bin/bash
# llama-go 快速启动脚本

echo "=== llama-go 插件化推理网关 ==="
echo ""

# 检查配置文件
if [ ! -f "config.yaml" ]; then
    echo "错误: config.yaml 不存在"
    exit 1
fi

# 检查后端服务
echo "检查后端服务..."
if curl -s http://localhost:8000/health > /dev/null 2>&1; then
    echo "✓ llama-cpp-python 后端运行中"
else
    echo "✗ llama-cpp-python 后端未启动"
    echo "  请先运行: python python/start_llama_cpp.py"
fi

echo ""
echo "启动网关..."
go run cmd/inference/main.go
