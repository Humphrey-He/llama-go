# Add Ollama to PATH

$ollamaPath = "$env:LOCALAPPDATA\Programs\Ollama"

# Check if already in PATH
$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($currentPath -notlike "*$ollamaPath*") {
    Write-Host "Adding Ollama to PATH..." -ForegroundColor Yellow
    [Environment]::SetEnvironmentVariable(
        "Path",
        "$currentPath;$ollamaPath",
        "User"
    )
    Write-Host "Done! Restart PowerShell to use 'ollama' command" -ForegroundColor Green
} else {
    Write-Host "Ollama already in PATH" -ForegroundColor Green
}
