$body = @{
    email = "admin@catter-cms.com"
    password = "admin123"
} | ConvertTo-Json

Write-Host "Testing /api/auth/login..."
try {
    $response = Invoke-RestMethod -Uri "http://127.0.0.1:3000/api/auth/login" -Method POST -ContentType "application/json" -Body $body
    Write-Host "Success: $($response | ConvertTo-Json)"
} catch {
    Write-Host "Failed: $_"
    Write-Host "Status: $($_.Exception.Response.StatusCode.value__)"
}
