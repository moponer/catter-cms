$token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJlbWFpbCI6ImFkbWluQGNhdHRlci1jbXMuY29tIiwicm9sZSI6ImFkbWluIn0.YP563V9YinUSViA0K56hkJSdwWJYAxtErYENQL7oxS8"

$headers = @{
    Authorization = "Bearer $token"
}

# 测试仪表盘API（简单测试数据库连接）
Write-Host "Testing /api/dashboard/stats..."
try {
    $response = Invoke-RestMethod -Uri "http://127.0.0.1:3000/api/dashboard/stats" -Method GET -Headers $headers
    Write-Host "Success: $($response | ConvertTo-Json)"
} catch {
    Write-Host "Failed: $_"
    $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
    $reader.BaseStream.Position = 0
    $responseBody = $reader.ReadToEnd()
    Write-Host "Response: $responseBody"
}

# 测试文章API
Write-Host "`nTesting /api/posts..."
try {
    $response = Invoke-RestMethod -Uri "http://127.0.0.1:3000/api/posts" -Method GET -Headers $headers
    Write-Host "Success"
} catch {
    Write-Host "Failed: $_"
}
