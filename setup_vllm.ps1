# vLLM 安装脚本 - Windows GPU 版本

Write-Host "=== vLLM 安装脚本 ===" -ForegroundColor Green

# 1. 创建虚拟环境
if (-not (Test-Path "venv")) {
    Write-Host "创建虚拟环境..." -ForegroundColor Yellow
    python -m venv venv
}

# 2. 激活虚拟环境
Write-Host "激活虚拟环境..." -ForegroundColor Yellow
.\venv\Scripts\Activate.ps1

# 3. 升级 pip
Write-Host "升级 pip..." -ForegroundColor Yellow
python -m pip install --upgrade pip setuptools wheel

# 4. 安装 PyTorch (CUDA 12.1)
Write-Host "安装 PyTorch (CUDA 12.1)..." -ForegroundColor Yellow
python -m pip install torch --index-url https://download.pytorch.org/whl/cu121

# 5. 验证 PyTorch
Write-Host "验证 PyTorch..." -ForegroundColor Yellow
python -c "import torch; print(f'PyTorch: {torch.__version__}'); print(f'CUDA: {torch.version.cuda}'); print(f'GPU可用: {torch.cuda.is_available()}'); print(f'GPU数量: {torch.cuda.device_count()}')"

# 6. 安装 llama-cpp-python (vLLM 替代方案，支持 GGUF)
Write-Host "安装 llama-cpp-python (支持 GPU)..." -ForegroundColor Yellow
$env:CMAKE_ARGS="-DGGML_CUDA=on"
python -m pip install llama-cpp-python --upgrade --force-reinstall --no-cache-dir

# 7. 验证 llama-cpp-python
Write-Host "验证 llama-cpp-python..." -ForegroundColor Yellow
python -c "from llama_cpp import Llama; print('llama-cpp-python 安装成功')"

Write-Host "`n=== 安装完成 ===" -ForegroundColor Green
Write-Host "运行服务: python python/start_llama_cpp.py" -ForegroundColor Cyan
