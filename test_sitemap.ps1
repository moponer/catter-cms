# Sitemap 404修复验证脚本
# 使用方法: .\test_sitemap.ps1

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Sitemap 404 修复验证脚本" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# 检查是否已编译
if (-Not (Test-Path "catter-cms-test.exe")) {
    Write-Host "❌ 错误: catter-cms-test.exe 不存在" -ForegroundColor Red
    Write-Host "请先运行: go build -o catter-cms-test.exe ./cmd/app/main.go" -ForegroundColor Yellow
    exit 1
}

Write-Host "✅ 找到 catter-cms-test.exe" -ForegroundColor Green
Write-Host ""

# 启动应用
Write-Host "正在启动应用..." -ForegroundColor Yellow
$process = Start-Process -FilePath ".\catter-cms-test.exe" -PassThru -WindowStyle Hidden

# 等待应用启动
Start-Sleep -Seconds 3

Write-Host "✅ 应用已启动 (PID: $($process.Id))" -ForegroundColor Green
Write-Host ""

# 测试sitemap.xml
Write-Host "正在测试 sitemap.xml..." -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "http://127.0.0.1:3001/sitemap.xml" -UseBasicParsing -TimeoutSec 5

    if ($response.StatusCode -eq 200) {
        Write-Host "✅ sitemap.xml 访问成功 (HTTP 200)" -ForegroundColor Green
        Write-Host "   Content-Type: $($response.Headers['Content-Type'])" -ForegroundColor Cyan

        # 检查XML格式
        if ($response.Content -like "*<urlset*") {
            Write-Host "   ✅ XML格式正确" -ForegroundColor Green
        } else {
            Write-Host "   ❌ XML格式错误" -ForegroundColor Red
        }

        # 显示前5行
        Write-Host ""
        Write-Host "Sitemap 内容预览:" -ForegroundColor Cyan
        $lines = $response.Content -split "`n"
        for ($i = 0; $i -lt [Math]::Min(5, $lines.Count); $i++) {
            Write-Host "   $($lines[$i])"
        }

        if ($lines.Count -gt 5) {
            Write-Host "   ... (共 $($lines.Count) 行)" -ForegroundColor Gray
        }
    } else {
        Write-Host "❌ sitemap.xml 返回 HTTP $($response.StatusCode)" -ForegroundColor Red
    }
} catch {
    Write-Host "❌ 访问 sitemap.xml 失败: $_" -ForegroundColor Red
}

Write-Host ""

# 测试Sitemap Info API
Write-Host "正在测试 Sitemap Info API..." -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "http://127.0.0.1:3001/api/admin/sitemap/info" -UseBasicParsing -TimeoutSec 5

    if ($response.StatusCode -eq 200) {
        Write-Host "✅ Sitemap Info API 访问成功 (HTTP 200)" -ForegroundColor Green
        Write-Host "   Response: $($response.Content)" -ForegroundColor Cyan
    } else {
        Write-Host "❌ Sitemap Info API 返回 HTTP $($response.StatusCode)" -ForegroundColor Red
    }
} catch {
    Write-Host "❌ 访问 Sitemap Info API 失败: $_" -ForegroundColor Red
}

Write-Host ""

# 测试Sitemap Rebuild API
Write-Host "正在测试 Sitemap Rebuild API..." -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "http://127.0.0.1:3001/api/admin/sitemap/rebuild" -Method POST -UseBasicParsing -TimeoutSec 5

    if ($response.StatusCode -eq 200) {
        Write-Host "✅ Sitemap Rebuild API 访问成功 (HTTP 200)" -ForegroundColor Green
        Write-Host "   Response: $($response.Content)" -ForegroundColor Cyan
    } else {
        Write-Host "❌ Sitemap Rebuild API 返回 HTTP $($response.StatusCode)" -ForegroundColor Red
    }
} catch {
    Write-Host "❌ 访问 Sitemap Rebuild API 失败: $_" -ForegroundColor Red
}

Write-Host ""

# 停止应用
Write-Host "正在停止应用..." -ForegroundColor Yellow
Stop-Process -Id $process.Id -Force
Write-Host "✅ 应用已停止" -ForegroundColor Green

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "验证完成！" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
