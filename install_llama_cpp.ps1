# llama-cpp-python 快速安装脚本

Write-Host "=== 安装 llama-cpp-python ===" -ForegroundColor Green

# 检查 Python
$pythonVersion = python --version 2>&1
Write-Host "Python: $pythonVersion" -ForegroundColor Cyan

# 检查虚拟环境
if (-not (Test-Path "venv")) {
    Write-Host "创建虚拟环境..." -ForegroundColor Yellow
    python -m venv venv
}

# 激活虚拟环境
Write-Host "激活虚拟环境..." -ForegroundColor Yellow
.\venv\Scripts\Activate.ps1

# 升级 pip
Write-Host "升级 pip..." -ForegroundColor Yellow
python -m pip install --upgrade pip -q

# 安装 llama-cpp-python (CPU 版本，快速安装)
Write-Host "安装 llama-cpp-python (CPU)..." -ForegroundColor Yellow
pip install llama-cpp-python[server] -q

Write-Host ""
Write-Host "=== 安装完成 ===" -ForegroundColor Green
Write-Host "启动服务: python python/start_llama_cpp.py" -ForegroundColor Cyan
