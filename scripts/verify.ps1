param([switch]$Fresh)
$ErrorActionPreference = 'Stop'
$projectRoot = Split-Path -Parent $PSScriptRoot
Set-Location -LiteralPath $projectRoot
function Run-Check([string]$Name, [scriptblock]$Command) {
    Write-Host "Checking: $Name"
    & $Command
    if ($LASTEXITCODE -ne 0) { throw "$Name failed with exit code $LASTEXITCODE" }
}
# The script owns only the fixed syncframe Compose project. Fresh explicitly deletes its demo data.
if ($Fresh) { Run-Check 'reset own demo project' { docker compose -p syncframe --profile multi down -v --remove-orphans } }
Run-Check 'Compose build and startup' { docker compose -p syncframe --profile multi up -d --build --wait --wait-timeout 180 }
$env:TEST_MYSQL_DSN='syncframe:syncframe@tcp(127.0.0.1:33060)/syncframe?parseTime=true&timeout=3s&readTimeout=5s&writeTimeout=5s'
$env:TEST_ADMIN_DSN='root:local-root-only@tcp(127.0.0.1:33060)/syncframe?parseTime=true'
$env:TEST_ALLOW_FAULTS='1'
$env:TEST_BASE_URL='http://127.0.0.1:8080'
$env:TEST_SECOND_URL='http://127.0.0.1:8081'
$env:TEST_KAFKA_BROKER='127.0.0.1:19092'
$env:TEST_COMPOSE_PROJECT='syncframe'
Run-Check 'Go unit, MySQL integration and API checks' { go test -count=1 -v ./... }
Run-Check 'Go static analysis' { go vet ./... }
Run-Check 'Go dependency security' { go run golang.org/x/vuln/cmd/govulncheck@v1.9.0 ./... }
Push-Location web
try {
    Run-Check 'frontend locked install' { npm.cmd ci }
    Run-Check 'frontend tests' { npm.cmd test }
    Run-Check 'frontend build' { npm.cmd run build }
    Run-Check 'frontend dependency security' { npm.cmd audit }
    Run-Check 'Playwright browser install' { npx.cmd playwright install chromium }
    Run-Check 'browser E2E' { npm.cmd run e2e }
} finally { Pop-Location }
$env:TEST_DOCKER_FAULTS='1'
Run-Check 'real container outage scenarios' { go test -count=1 -timeout 5m -v ./tests -run TestMultiServerAndInfrastructureFaults }
Remove-Item Env:TEST_DOCKER_FAULTS
Run-Check 'WebSocket staged load' { go run ./cmd/load -users 1,5,20 -events 30 }
Write-Host 'Verification commands passed. Reports are under artifacts/. Read docs/testing.md for boundaries.'
