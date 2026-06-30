# Ollama 部署指南（推荐方案）

## 为什么选择 Ollama

llama-cpp-python 在 Windows 上编译困难，Ollama 是预编译的替代方案：
- 一键安装，无需编译
- 自动 GPU 加速
- OpenAI 兼容 API
- 支持多种模型

## 快速开始

### 1. 安装 Ollama

下载并安装：https://ollama.com/download

安装后会自动启动后台服务（端口 11434）

### 2. 下载模型

```powershell
# 下载 Llama 3.1 8B 模型（约 4.7GB）
ollama pull llama3.1:8b

# 或使用你的 GGUF 模型
ollama create dolphin -f Modelfile
```

### 3. 测试 Ollama

```powershell
# 测试对话
ollama run llama3.1:8b "Hello"

# 检查 API
curl http://localhost:11434/api/tags
```

### 4. 更新 Go 服务配置

修改 `cmd/inference/main.go`，后端 URL 改为 Ollama：

```go
pythonBackend := backend.NewPythonBackend("http://localhost:11434")
```

### 5. 启动服务

```bash
# 启动 Go 服务
go run cmd/inference/main.go

# 测试
curl http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{"model":"llama3.1:8b","messages":[{"role":"user","content":"Hello"}]}'
```

## 使用自定义 GGUF 模型

创建 Modelfile：

```dockerfile
FROM ./internal/model/Dolphin3.0-Llama3.1-8B-Q4_K_S.gguf
```

导入模型：

```powershell
ollama create dolphin -f Modelfile
ollama run dolphin "Hello"
```

## 优势

- 无需 Python 环境
- 自动使用 GPU（GTX 1660 Ti）
- 比 llama-cpp-python 更快
- 支持模型热切换

## 故障排查

检查 Ollama 服务：
```powershell
# 查看日志
ollama logs

# 重启服务
ollama serve
```
