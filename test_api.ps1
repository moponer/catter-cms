# Test protected API with token
$token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJlbWFpbCI6ImFkbWluQGNhdHRlci1jbXMuY29tIiwicm9sZSI6ImFkbWluIn0.YP563V9YinUSViA0K56hkJSdwWJYAxtErYENQL7oxS8"

$headers = @{
    Authorization = "Bearer $token"
}

Write-Host "Testing /api/auth/me..."
try {
    $response = Invoke-RestMethod -Uri "http://127.0.0.1:3000/api/auth/me" -Method GET -Headers $headers
    Write-Host "Success: $($response | ConvertTo-Json)"
} catch {
    Write-Host "Failed: $_"
}

Write-Host "`nTesting /api/posts..."
try {
    $response = Invoke-RestMethod -Uri "http://127.0.0.1:3000/api/posts" -Method GET -Headers $headers
    Write-Host "Success: $($response | ConvertTo-Json)"
} catch {
    Write-Host "Failed: $_"
}

Write-Host "`nTesting /api/categories..."
try {
    $response = Invoke-RestMethod -Uri "http://127.0.0.1:3000/api/categories" -Method GET -Headers $headers
    Write-Host "Success: $($response | ConvertTo-Json)"
} catch {
    Write-Host "Failed: $_"
}
