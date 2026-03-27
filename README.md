# Llama-Go

基于 Go 的 Llama 大模型本地化推理框架

## 项目定位

Llama-Go 是一个轻量级的 Go 语言框架，用于在本地部署和运行 Llama 大语言模型。支持模型加载、文本推理、API 服务等核心功能，适合学习、研究和小规模应用场景。

## 核心特性

- ✅ 本地化推理：无需云服务，完全离线运行
- ✅ 轻量级框架：基于 Go，编译快、资源占用低
- ✅ 灵活配置：支持推理参数自定义
- ✅ API 服务：提供 HTTP 接口
- ✅ 并发处理：高效的并发请求处理

## 环境要求

- Go 1.21+
- 内存：≥ 8GB（推荐 16GB+）
- 磁盘：≥ 10GB（用于模型文件）
- 可选：GPU（NVIDIA/AMD，用于加速推理）

## 快速开始

### 1. 克隆项目
```bash
git clone https://github.com/Humphrey-He/llama-go.git
cd llama-go
```

### 2. 下载模型
从 [Hugging Face](https://huggingface.co/models?search=llama) 下载 Llama 模型文件（GGUF 格式）

### 3. 编译运行
```bash
go build -o bin/llama cmd/llama/main.go
./bin/llama --model /path/to/model.gguf
```

## 项目结构

```
llama-go/
├── cmd/              # 可执行程序入口
│   └── llama/        # 主程序
├── pkg/              # 公共库
│   └── llama/        # 核心推理逻辑
├── internal/         # 内部代码
├── config/           # 配置文件
├── docs/             # 文档
├── test/             # 测试
├── examples/         # 示例代码
├── .gitignore
├── Dockerfile
├── go.mod
└── README.md
```

## 文档

- [CONTRIBUTING.md](CONTRIBUTING.md) - 贡献指南
- [docs/](docs/) - 详细文档

## 许可证

MIT License - 详见 [LICENSE](LICENSE)

## 免责声明

本项目仅供学习和研究使用。使用者需自行确保遵守 Llama 模型的官方授权协议。
