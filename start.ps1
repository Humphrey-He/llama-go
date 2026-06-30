# llama-go 快速启动脚本 (Windows)

Write-Host "=== llama-go 插件化推理网关 ===" -ForegroundColor Green
Write-Host ""

# 检查配置文件
if (-not (Test-Path "config.yaml")) {
    Write-Host "错误: config.yaml 不存在" -ForegroundColor Red
    exit 1
}

# 检查后端服务
Write-Host "检查后端服务..." -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "http://localhost:8000/health" -TimeoutSec 2 -ErrorAction Stop
    Write-Host "✓ llama-cpp-python 后端运行中" -ForegroundColor Green
} catch {
    Write-Host "✗ llama-cpp-python 后端未启动" -ForegroundColor Red
    Write-Host "  请先运行: python python/start_llama_cpp.py" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "启动网关..." -ForegroundColor Yellow
go run cmd/inference/main.go
