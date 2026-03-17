$token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJlbWFpbCI6ImFkbWluQGNhdHRlci1jbXMuY29tIiwicm9sZSI6ImFkbWluIn0.YP563V9YinUSViA0K56hkJSdwWJYAxtErYENQL7oxS8"

$headers = @{
    Authorization = "Bearer $token"
    "Content-Type" = "application/json"
}

# 测试创建标签
Write-Host "Testing POST /api/tags (Create)..."
try {
    $body = @{ name = "Test Tag"; slug = "test-tag" } | ConvertTo-Json
    $response = Invoke-RestMethod -Uri "http://127.0.0.1:3000/api/tags" -Method POST -Headers $headers -Body $body
    Write-Host "Create Success: $($response | ConvertTo-Json)"
    $newTagId = $response.data.id
    Write-Host "New Tag ID: $newTagId"
} catch {
    Write-Host "Create Failed: $_"
    $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
    $reader.BaseStream.Position = 0
    Write-Host "Response: $($reader.ReadToEnd())"
}
