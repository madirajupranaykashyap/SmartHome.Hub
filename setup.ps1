# SmartHome Hub - Single Executable Build Script (Windows)

param(
    [ValidateSet('current', 'all', 'linux', 'windows', 'macos')]
    [string]$Platform = 'current',
    
    [string]$FrontendPath,
    [switch]$SkipFrontend,
    [switch]$TestApp
)

$ErrorActionPreference = "Stop"

# Color functions
function Write-Step {
    Write-Host "▶ $args" -ForegroundColor Green
}

function Write-Error {
    Write-Host "✗ $args" -ForegroundColor Red
}

function Write-Success {
    Write-Host "✓ $args" -ForegroundColor Green
}

function Write-Info {
    Write-Host "ℹ $args" -ForegroundColor Yellow
}

# Header
Write-Host "`n╔════════════════════════════════════════╗" -ForegroundColor Blue
Write-Host "║   SmartHome Hub - Single Executable   ║" -ForegroundColor Blue
Write-Host "║      Build & Setup Script (Windows)   ║" -ForegroundColor Blue
Write-Host "╚════════════════════════════════════════╝`n" -ForegroundColor Blue

# Check dependencies
Write-Step "Checking dependencies..."

$goVersion = go version 2>$null
if ($LASTEXITCODE -ne 0) {
    Write-Error "Go is not installed"
    exit 1
}
Write-Success "Go installed: $($goVersion -split ' ' | Select-Object -Index 2)"

$nodeVersion = node --version 2>$null
if ($LASTEXITCODE -eq 0) {
    Write-Success "Node.js $nodeVersion"
} else {
    Write-Info "Node.js not found (optional, needed only if building frontend)"
}

$npmVersion = npm --version 2>$null
if ($LASTEXITCODE -eq 0) {
    Write-Success "npm $npmVersion"
}

Write-Host ""

# Prepare frontend
if (-not $SkipFrontend) {
    Write-Step "Preparing Angular frontend..."
    
    if (-not (Test-Path "frontend")) {
        if ($FrontendPath) {
            if (Test-Path $FrontendPath) {
                Copy-Item -Path $FrontendPath -Destination "frontend" -Recurse
                Write-Success "Frontend copied from $FrontendPath"
            } else {
                Write-Error "Frontend path not found: $FrontendPath"
                exit 1
            }
        } else {
            Write-Info "Frontend directory not found. Skipping frontend build."
            Write-Info "To set up later: copy your Angular dist/ to backend/pkg/hub/frontend/"
        }
    }
    
    if (Test-Path "frontend\package.json") {
        Write-Step "Building Angular frontend..."
        
        Push-Location frontend
        
        if (-not (Test-Path "node_modules")) {
            npm install
        }
        
        npm run build
        Write-Success "Frontend built successfully"
        
        Pop-Location
    }
}

Write-Host ""

# Embed frontend
if (Test-Path "frontend\dist") {
    Write-Step "Embedding frontend in Go binary..."
    
    Push-Location backend
    
    $frontendDir = "pkg\hub\frontend"
    if (-not (Test-Path $frontendDir)) {
        New-Item -ItemType Directory -Path $frontendDir -Force | Out-Null
    }
    
    # Copy frontend files
    Get-ChildItem -Path "..\frontend\dist" -Recurse | Copy-Item -Destination $frontendDir -Force -Recurse
    
    if (Test-Path "$frontendDir\index.html") {
        Write-Success "Frontend embedded successfully"
        $fileCount = (Get-ChildItem -Path $frontendDir -Recurse).Count
        Write-Info "Files: $fileCount items"
    } else {
        Write-Error "Frontend files not embedded correctly"
    }
    
    Pop-Location
}

Write-Host ""

# Build executable
Write-Step "Building executable..."

Push-Location backend

switch ($Platform) {
    'current' {
        Write-Step "Building for current platform..."
        & make build
        Write-Success "Executable created: .\build\smarthome-hub.exe (or without .exe on non-Windows)"
    }
    'all' {
        Write-Step "Building for all platforms..."
        & make build-all
        Write-Success "Executables created in .\build\"
    }
    'linux' {
        Write-Step "Building for Linux..."
        & make build-linux
        Write-Success "Executables created in .\build\linux\"
    }
    'windows' {
        Write-Step "Building for Windows..."
        & make build-windows
        Write-Success "Executables created in .\build\windows\"
    }
    'macos' {
        Write-Step "Building for macOS..."
        & make build-macos
        Write-Success "Executables created in .\build\macos\"
    }
}

Pop-Location

Write-Host ""

# Test application
if ($TestApp) {
    Write-Step "Testing application..."
    
    Push-Location backend
    
    if (Test-Path "build\smarthome-hub.exe") {
        Write-Info "Starting application (will stop after 5 seconds)..."
        
        $process = Start-Process -FilePath ".\build\smarthome-hub.exe" `
                                 -NoNewWindow `
                                 -PassThru
        
        Start-Sleep -Seconds 5
        
        if (-not $process.HasExited) {
            Stop-Process -InputObject $process -Force
        }
        
        if (Test-Path "config.json") {
            Write-Success "Configuration created"
            Write-Info "Config file: config.json"
        }
    }
    
    Pop-Location
}

Write-Host ""
Write-Host "╔════════════════════════════════════════╗" -ForegroundColor Blue
Write-Success "Setup complete!"
Write-Host "╚════════════════════════════════════════╝" -ForegroundColor Blue
Write-Host ""

Write-Host "Next steps:" -ForegroundColor Green
Write-Host "  1. Test locally: .\backend\build\smarthome-hub.exe" -ForegroundColor Yellow
Write-Host "  2. Access UI: http://localhost:8080" -ForegroundColor Yellow
Write-Host "  3. Configure: Edit config.json" -ForegroundColor Yellow
Write-Host "  4. Deploy: Upload executable to your server" -ForegroundColor Yellow
Write-Host ""

Write-Host "Documentation:" -ForegroundColor Green
Write-Host "  - Setup guide: SETUP.md" -ForegroundColor Yellow
Write-Host "  - Deployment: DEPLOYMENT.md" -ForegroundColor Yellow
Write-Host "  - Build help: cd backend && make help" -ForegroundColor Yellow
