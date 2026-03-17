$token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJlbWFpbCI6ImFkbWluQGNhdHRlci1jbXMuY29tIiwicm9sZSI6ImFkbWluIn0.YP563V9YinUSViA0K56hkJSdwWJYAxtErYENQL7oxS8"

$headers = @{
    Authorization = "Bearer $token"
}

Write-Host "Testing /api/themes/catter-blog/config..."
try {
    $response = Invoke-RestMethod -Uri "http://127.0.0.1:3000/api/themes/catter-blog/config" -Method GET -Headers $headers
    Write-Host "Success: $($response | ConvertTo-Json -Depth 5)"
} catch {
    Write-Host "Failed: $_"
    $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
    $reader.BaseStream.Position = 0
    Write-Host "Response: $($reader.ReadToEnd())"
}
