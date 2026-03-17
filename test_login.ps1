# Test login API
$body = @{
    email = "admin@catter-cms.com"
    password = "admin123"
} | ConvertTo-Json

Write-Host "Testing login API..."
Write-Host "Request body: $body"

try {
    $response = Invoke-RestMethod -Uri "http://127.0.0.1:3000/api/auth/login" -Method POST -ContentType "application/json" -Body $body
    Write-Host "Login successful!"
    Write-Host "Token: $($response.token)"
    Write-Host "User: $($response.user | ConvertTo-Json)"
} catch {
    Write-Host "Login failed: $_"
    Write-Host "Response: $($_.Exception.Response.StatusCode.value__)"
}
