param(
    [string]$BaseUrl = "http://127.0.0.1:3000",
    [string]$Email = "admin@example.com",
    [string]$Password = "Admin123!",
    [string]$ClientIp = "198.51.100.10"
)

$ErrorActionPreference = "Stop"

function Assert-Contract([bool]$Condition, [string]$Message) {
    if (-not $Condition) {
        throw "Contract check failed: $Message"
    }
}

$spec = Invoke-RestMethod -Uri "$BaseUrl/api/openapi.json"
Assert-Contract ($spec.openapi -eq "3.0.3") "OpenAPI version is 3.0.3"
foreach ($path in @("/auth/login", "/auth/refresh", "/auth/logout-all", "/auth/me", "/admin/registry", "/admin/search", "/admin/{resource}", "/admin/{resource}/{id}", "/admin/{resource}/export", "/admin/{resource}/actions/{action}", "/admin/{resource}/relations/{relation}/options", "/admin/{resource}/{id}/relations/{relation}", "/admin/overview", "/admin/audit-logs")) {
    Assert-Contract ($null -ne $spec.paths.$path) "contract path exists: $path"
}

$docs = & curl.exe -fsS "$BaseUrl/api/docs"
Assert-Contract ($LASTEXITCODE -eq 0) "Scalar docs page returns 200"
Assert-Contract (($docs -join "`n") -match "api-reference") "Scalar docs page contains the reference element"

$loginBody = @{ email = $Email; password = $Password } | ConvertTo-Json
$requestHeaders = @{ "X-Forwarded-For" = $ClientIp }
$login = Invoke-RestMethod -Uri "$BaseUrl/api/v1/auth/login" -Method Post -Headers $requestHeaders -ContentType "application/json" -Body $loginBody
$token = $login.data.access_token
Assert-Contract (-not [string]::IsNullOrWhiteSpace($token)) "login returns an access token"
$headers = @{ Authorization = "Bearer $token" }

$resources = Invoke-RestMethod -Uri "$BaseUrl/api/v1/admin/registry" -Headers $headers
Assert-Contract ($resources.data.Count -ge 1) "resource manifest returns data"
$users = Invoke-RestMethod -Uri "$BaseUrl/api/v1/admin/users?page=1&per_page=10&sort=id&dir=desc" -Headers $headers
Assert-Contract ($null -ne $users.meta) "user resource list returns pagination metadata"
Assert-Contract ($null -ne $users.data) "user resource list returns data"

$announcements = Invoke-RestMethod -Uri "$BaseUrl/api/v1/admin/announcements?page=1&per_page=10&sort=id&dir=desc" -Headers $headers
Assert-Contract ($null -ne $announcements.meta) "generated resource list returns pagination metadata"
Assert-Contract ($null -ne $announcements.data) "generated resource list returns data"

$search = Invoke-RestMethod -Uri "$BaseUrl/api/v1/admin/search?q=announcement&limit=10" -Headers $headers
Assert-Contract ($null -ne $search.data) "global search returns data"

Write-Output "Contract smoke passed: OpenAPI, Scalar, login, registry, users list, generated resource list, and global search."
