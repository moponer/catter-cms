# Catter CMS Build Script

$ErrorActionPreference = "Stop"
$projectRoot = $PSScriptRoot
$deployDir = Join-Path $projectRoot "deploy"

Write-Host "Building Catter CMS..."

# Create directories
$dirs = @($deployDir, "$deployDir/uploads", "$deployDir/templates/themes", "$deployDir/templates/fontawesome", "$deployDir/public")
foreach ($dir in $dirs) { if (-not (Test-Path $dir)) { New-Item -ItemType Directory -Path $dir -Force | Out-Null } }

# Build backend
Set-Location "$projectRoot/gin-app"
go build -ldflags="-s -w" -o "$deployDir/catter-cms.exe" ./cmd

# Build frontend
Set-Location "$projectRoot/admin"
npm run build 2>&1 | Out-Null

# Copy files
if (Test-Path "$deployDir/admin") { Remove-Item -Recurse -Force "$deployDir/admin" }
Copy-Item -Recurse -Force "$projectRoot/admin/dist" "$deployDir/admin"
Copy-Item -Force "$projectRoot/config.yaml" $deployDir
Copy-Item -Recurse -Force "$projectRoot/templates/themes/*" "$deployDir/templates/themes/"
Copy-Item -Recurse -Force "$projectRoot/templates/fontawesome/*" "$deployDir/templates/fontawesome/"
if (Test-Path "$projectRoot/catter.db") { Copy-Item -Force "$projectRoot/catter.db" $deployDir }

Set-Location $projectRoot
Write-Host "Build complete: $deployDir"