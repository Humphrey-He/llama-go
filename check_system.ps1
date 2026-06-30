# System Information Detection Script

[Console]::OutputEncoding = [System.Text.Encoding]::UTF8
$OutputEncoding = [System.Text.Encoding]::UTF8

Write-Host "`n=== System Information ===" -ForegroundColor Cyan

# CPU
Write-Host "`n[CPU]" -ForegroundColor Yellow
$cpu = Get-WmiObject Win32_Processor | Select-Object -First 1
Write-Host "Model: $($cpu.Name)"
Write-Host "Cores: $($cpu.NumberOfCores)"
Write-Host "Threads: $($cpu.NumberOfLogicalProcessors)"

# Memory
Write-Host "`n[Memory]" -ForegroundColor Yellow
$mem = Get-WmiObject Win32_ComputerSystem
$totalMem = [math]::Round($mem.TotalPhysicalMemory / 1GB, 2)
Write-Host "Total: $totalMem GB"

# GPU
Write-Host "`n[GPU]" -ForegroundColor Yellow
try {
    $nvidiaSmi = & nvidia-smi --query-gpu=name,driver_version,memory.total --format=csv,noheader 2>$null
    if ($LASTEXITCODE -eq 0) {
        $gpuInfo = $nvidiaSmi -split ','
        Write-Host "Model: $($gpuInfo[0].Trim())"
        Write-Host "Driver: $($gpuInfo[1].Trim())"
        Write-Host "VRAM: $($gpuInfo[2].Trim())"

        $cudaVersion = & nvidia-smi 2>$null | Select-String "CUDA Version: (\d+\.\d+)"
        if ($cudaVersion) {
            Write-Host "CUDA Driver: $($cudaVersion.Matches.Groups[1].Value)"
        }
    } else {
        Write-Host "No NVIDIA GPU detected" -ForegroundColor Red
    }
} catch {
    Write-Host "nvidia-smi not available" -ForegroundColor Red
}

# CUDA Toolkit
Write-Host "`n[CUDA Toolkit]" -ForegroundColor Yellow
try {
    $nvcc = & nvcc --version 2>$null | Select-String "release (\d+\.\d+)"
    if ($nvcc) {
        Write-Host "Installed: $($nvcc.Matches.Groups[1].Value)" -ForegroundColor Green
    } else {
        Write-Host "Not installed" -ForegroundColor Red
    }
} catch {
    Write-Host "Not installed" -ForegroundColor Red
}

# Python
Write-Host "`n[Python]" -ForegroundColor Yellow
$pythonVersion = & python --version 2>&1
Write-Host "Version: $pythonVersion"

Write-Host "`n=== Detection Complete ===" -ForegroundColor Cyan
