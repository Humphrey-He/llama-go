# llama-cpp-python Installation Script - Python 3.11 Required

[Console]::OutputEncoding = [System.Text.Encoding]::UTF8
$OutputEncoding = [System.Text.Encoding]::UTF8

Write-Host "`n=== llama-cpp-python Installation (CPU) ===" -ForegroundColor Green

# Check Python version
Write-Host "`nChecking Python version..." -ForegroundColor Yellow
$pythonVersion = & python --version 2>&1
Write-Host "Current: $pythonVersion"

if ($pythonVersion -match "3\.13") {
    Write-Host "`nERROR: Python 3.13 is NOT compatible with llama-cpp-python!" -ForegroundColor Red
    Write-Host "Please install Python 3.11 from:" -ForegroundColor Yellow
    Write-Host "https://www.python.org/downloads/release/python-31110/" -ForegroundColor Cyan
    Write-Host "`nOr use: py -3.11 -m venv venv (if multiple versions installed)" -ForegroundColor Yellow
    exit 1
}

if ($pythonVersion -notmatch "3\.11") {
    Write-Host "WARNING: Python 3.11 is recommended for best compatibility" -ForegroundColor Yellow
    $continue = Read-Host "Continue anyway? (y/n)"
    if ($continue -ne "y") { exit 0 }
}

# Remove old venv
if (Test-Path "venv") {
    Write-Host "`nRemoving old virtual environment..." -ForegroundColor Yellow
    Remove-Item -Recurse -Force venv
}

# Create venv
Write-Host "Creating virtual environment..." -ForegroundColor Yellow
python -m venv venv

# Activate venv
Write-Host "Activating virtual environment..." -ForegroundColor Yellow
.\venv\Scripts\Activate.ps1

# Upgrade pip
Write-Host "`nUpgrading pip..." -ForegroundColor Yellow
python -m pip install --upgrade pip

# Install dependencies
Write-Host "Installing dependencies..." -ForegroundColor Yellow
python -m pip install numpy setuptools==70.2.0 torch --index-url https://download.pytorch.org/whl/cpu

# Install llama-cpp-python
Write-Host "`nInstalling llama-cpp-python (this may take 5-10 minutes)..." -ForegroundColor Yellow
$env:CMAKE_ARGS = "-DLLAMA_BLAS=ON -DLLAMA_BLAS_VENDOR=OpenBLAS"
python -m pip install llama-cpp-python

# Install uvicorn
Write-Host "`nInstalling uvicorn..." -ForegroundColor Yellow
python -m pip install "uvicorn[standard]"

# Verify installation
Write-Host "`nVerifying installation..." -ForegroundColor Yellow
python -c "from llama_cpp import Llama; print('Installation successful!')"

Write-Host "`n=== Installation Complete ===" -ForegroundColor Green
Write-Host "Run server: python python/start_llama_cpp.py" -ForegroundColor Cyan
Write-Host "Note: CPU mode, inference will be slower" -ForegroundColor Yellow
