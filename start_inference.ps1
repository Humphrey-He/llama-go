# llama-go 一键启动脚本
# 功能：安装 llama-cpp-python、启动 Python 后端服务、启动 Go API 服务

param(
    [switch]$SkipInstall,
    [switch]$SkipGoBuild
)

$ErrorActionPreference = "Stop"
$ProjectRoot = Split-Path -Parent $PSScriptRoot

Write-Host "=== Llama-Go 一键启动 ===" -ForegroundColor Cyan

# 1. 安装 llama-cpp-python（如果需要）
if (-not $SkipInstall) {
    Write-Host "`n[1/3] 安装 llama-cpp-python..." -ForegroundColor Yellow
    Write-Host "注意：首次安装需要编译，时间较长（约10-30分钟）..." -ForegroundColor Gray

    python -m pip install llama-cpp-python[server] --quiet
    if ($LASTEXITCODE -ne 0) {
        Write-Host "安装失败，尝试不包含server选项..." -ForegroundColor Red
        python -m pip install llama-cpp-python --quiet
    }
    Write-Host "安装完成" -ForegroundColor Green
}

# 2. 启动 Python 后端服务
Write-Host "`n[2/3] 启动 Python 后端服务 (llama-cpp)..." -ForegroundColor Yellow

# 检查模型文件
$modelPath = Join-Path $ProjectRoot "models\Dolphin3.0-Llama3.1-8B-Q4_K_S.gguf"
if (-not (Test-Path $modelPath)) {
    Write-Host "错误：模型文件不存在: $modelPath" -ForegroundColor Red
    Write-Host "请先下载 GGUF 模型文件" -ForegroundColor Yellow
    exit 1
}

$pythonScript = Join-Path $ProjectRoot "python\start_llama_cpp.py"
$pythonJob = Start-Job -ScriptBlock {
    param($script, $modelPath)
    Set-Location (Split-Path $script)
    & python $script
} -ArgumentList $pythonScript, $modelPath

Write-Host "Python 服务已在后台启动 (Job ID: $($pythonJob.Id))" -ForegroundColor Gray
Write-Host "等待服务就绪..." -ForegroundColor Gray

# 等待 Python 服务就绪（最多60秒）
$maxWait = 60
$ready = $false
for ($i = 0; $i -lt $maxWait; $i++) {
    Start-Sleep -Seconds 1
    try {
        $response = Invoke-WebRequest -Uri "http://localhost:8000/health" -TimeoutSec 2 -ErrorAction SilentlyContinue
        if ($response.StatusCode -eq 200) {
            $ready = $true
            Write-Host "Python 服务已就绪" -ForegroundColor Green
            break
        }
    } catch {
        Write-Host -NoNewline "`r等待中... ($($i+1)/$maxWait)"
    }
}

if (-not $ready) {
    Write-Host "`n警告：Python 服务可能未就绪，继续启动 Go 服务..." -ForegroundColor Yellow
}

# 3. 构建并启动 Go 服务
Write-Host "`n[3/3] 启动 Go API 服务..." -ForegroundColor Yellow

if (-not $SkipGoBuild) {
    Write-Host "编译 Go 服务..." -ForegroundColor Gray
    Set-Location $ProjectRoot
    go build -o bin/inference.exe ./cmd/inference/main.go
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Go 编译失败" -ForegroundColor Red
        exit 1
    }
}

$goBin = Join-Path $ProjectRoot "bin\inference.exe"
$goJob = Start-Job -ScriptBlock {
    param($bin)
    & $bin
} -ArgumentList $goBin

Write-Host "Go 服务已在后台启动 (Job ID: $($goJob.Id))" -ForegroundColor Gray
Start-Sleep -Seconds 2

# 4. 验证服务
Write-Host "`n=== 验证服务 ===" -ForegroundColor Cyan

# 检查 Go API
try {
    $health = Invoke-WebRequest -Uri "http://localhost:8080/healthz" -TimeoutSec 5 -ErrorAction SilentlyContinue
    if ($health.StatusCode -eq 200) {
        Write-Host "[OK] Go API 服务正常: http://localhost:8080" -ForegroundColor Green
    }
} catch {
    Write-Host "[FAIL] Go API 服务异常" -ForegroundColor Red
}

# 检查模型列表
try {
    $models = Invoke-RestMethod -Uri "http://localhost:8080/v1/models" -TimeoutSec 5 -ErrorAction SilentlyContinue
    Write-Host "[OK] 可用模型:" -ForegroundColor Green
    $models.data | ForEach-Object { Write-Host "  - $($_.id)" -ForegroundColor White }
} catch {
    Write-Host "[INFO] 无法获取模型列表（Python服务可能未就绪）" -ForegroundColor Yellow
}

# 5. 测试对话
Write-Host "`n=== 测试对话 ===" -ForegroundColor Cyan

$testPayload = @{
    model = "tinyllama-chat"
    messages = @(
        @{role = "user"; content = "Hello, how are you?"}
    )
    max_tokens = 100
    temperature = 0.7
} | ConvertTo-Json

Write-Host "发送测试请求..." -ForegroundColor Gray
Write-Host "模型: tinyllama-chat" -ForegroundColor Gray
Write-Host "问题: Hello, how are you?" -ForegroundColor Gray
Write-Host ""

try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/v1/chat/completions" `
        -Method Post `
        -ContentType "application/json" `
        -Body $testPayload `
        -TimeoutSec 120 `
        -ErrorAction SilentlyContinue

    Write-Host "回答:" -ForegroundColor Green
    Write-Host $response.choices[0].message.content -ForegroundColor White
} catch {
    Write-Host "调用失败: $_" -ForegroundColor Red
    Write-Host "提示：检查 Python 服务是否完全就绪（首次启动需要更长时间加载模型）" -ForegroundColor Yellow
}

Write-Host "`n=== 启动完成 ===" -ForegroundColor Cyan
Write-Host "Go API: http://localhost:8080" -ForegroundColor White
Write-Host "Python Backend: http://localhost:8000" -ForegroundColor White
Write-Host ""
Write-Host "停止服务:" -ForegroundColor Gray
Write-Host "  Stop-Job -Id $goJob.Id" -ForegroundColor Gray
Write-Host "  Stop-Job -Id $pythonJob.Id" -ForegroundColor Gray
