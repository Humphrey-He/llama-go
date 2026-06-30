#!/usr/bin/env pwsh
# 模型调用测试脚本
# 用法: .\test_model_call.ps1 [-UseMock] [-Model tinyllama-chat]

param(
    [switch]$UseMock,
    [string]$Model = "tinyllama-chat",
    [string]$Question = "Hello, how are you?"
)

$ErrorActionPreference = "Continue"

Write-Host "=== 模型调用测试 ===" -ForegroundColor Cyan
Write-Host "模型: $Model" -ForegroundColor White
Write-Host "问题: $Question" -ForegroundColor White
Write-Host ""

if ($UseMock) {
    Write-Host "[模式] 使用 Mock Backend (无需启动服务)" -ForegroundColor Yellow
    Write-Host ""
    # Mock 测试已通过单元测试验证
    Write-Host "提示: 运行以下命令进行 Mock 测试:" -ForegroundColor Gray
    Write-Host "  go test -v ./internal/api/... -run TestChat" -ForegroundColor Gray
} else {
    Write-Host "[模式] 测试真实 API 调用" -ForegroundColor Yellow
    Write-Host ""

    $payload = @{
        model = $Model
        messages = @(
            @{role = "user"; content = $Question}
        )
        max_tokens = 200
        temperature = 0.7
    } | ConvertTo-Json -Depth 10

    Write-Host "发送请求到 http://localhost:8080/v1/chat/completions ..." -ForegroundColor Gray
    Write-Host ""

    try {
        $startTime = Get-Date
        $response = Invoke-RestMethod -Uri "http://localhost:8080/v1/chat/completions" `
            -Method Post `
            -ContentType "application/json" `
            -Body $payload `
            -TimeoutSec 120

        $duration = (Get-Date) - $startTime

        Write-Host "=== 响应 ===" -ForegroundColor Green
        Write-Host "模型: $($response.model)" -ForegroundColor Gray
        Write-Host "回复: $($response.choices[0].message.content)" -ForegroundColor White
        Write-Host ""
        Write-Host "=== 统计 ===" -ForegroundColor Green
        Write-Host "耗时: $($duration.TotalSeconds)s" -ForegroundColor Gray
        Write-Host "Prompt Tokens: $($response.usage.prompt_tokens)" -ForegroundColor Gray
        Write-Host "Completion Tokens: $($response.usage.completion_tokens)" -ForegroundColor Gray
        Write-Host "Total Tokens: $($response.usage.total_tokens)" -ForegroundColor Gray
    } catch {
        Write-Host "调用失败: $_" -ForegroundColor Red
        Write-Host ""
        Write-Host "请确保服务已启动:" -ForegroundColor Yellow
        Write-Host "  1. 启动 Python 后端: python python/start_llama_cpp.py" -ForegroundColor Gray
        Write-Host "  2. 启动 Go API: go run cmd/inference/main.go" -ForegroundColor Gray
    }
}
