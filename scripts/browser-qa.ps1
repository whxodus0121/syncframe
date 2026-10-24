param([string]$RunId = ([DateTime]::UtcNow.ToString('yyyyMMdd-HHmmssfff') + '-' + [Guid]::NewGuid().ToString('N')))
$ErrorActionPreference = 'Stop'
if ($RunId -notmatch '^[a-zA-Z0-9_-]+$') { throw 'RunId must contain only letters, digits, hyphens or underscores' }
$qaProjectRoot = Split-Path -Parent $PSScriptRoot
Set-Location -LiteralPath $qaProjectRoot
$qaOutput = Join-Path $qaProjectRoot "artifacts\browser-qa\runs\$RunId"
if (Test-Path -LiteralPath $qaOutput) { throw 'This run already exists. Use a new RunId to preserve earlier evidence.' }
New-Item -ItemType Directory -Path $qaOutput -Force | Out-Null
function Check([string]$Name,[scriptblock]$Command) {
    Write-Host "Checking $Name"
    & $Command *> (Join-Path $qaOutput "$Name.txt")
    $checkExit = $LASTEXITCODE
    Get-Content -LiteralPath (Join-Path $qaOutput "$Name.txt") -Tail 12
    if ($checkExit -ne 0) { throw "$Name failed (exit $checkExit). See $qaOutput" }
}
Check 'compose' { docker compose -p syncframe --profile multi up -d --build --wait --wait-timeout 180 }
$env:TEST_MYSQL_DSN='syncframe:syncframe@tcp(127.0.0.1:33060)/syncframe?parseTime=true&timeout=3s&readTimeout=5s&writeTimeout=5s'
$env:TEST_ADMIN_DSN='root:local-root-only@tcp(127.0.0.1:33060)/syncframe?parseTime=true'
$env:TEST_ALLOW_FAULTS='1'
$env:TEST_BASE_URL='http://localhost:8080'
$env:TEST_SECOND_URL='http://localhost:8081'
$env:TEST_KAFKA_BROKER='127.0.0.1:19092'
$env:TEST_COMPOSE_PROJECT='syncframe'
Remove-Item Env:TEST_DOCKER_FAULTS -ErrorAction SilentlyContinue
Check 'go-tests' { go test -count=1 -v ./... }
Check 'go-vet' { go vet ./... }
Check 'go-security' { go run golang.org/x/vuln/cmd/govulncheck@v1.9.0 ./... }
$env:QA_RUN_ID=$RunId
$env:BASE_URL='http://localhost:8080'
Push-Location web
try {
    Check 'frontend-install' { npm.cmd ci }
    Check 'frontend-tests' { npm.cmd test }
    Check 'frontend-build' { npm.cmd run build }
    Check 'frontend-security' { npm.cmd audit }
    Check 'browser-install' { npx.cmd playwright install chromium }
    try { Check 'browser-tests' { npm.cmd run browser:qa } }
    finally { & (Join-Path $qaProjectRoot 'scripts\sanitize-browser-traces.ps1') }
    node -e 'const fs=require("fs"); const {chromium}=require("playwright"); (async()=>{const b=await chromium.launch(); const report={at:new Date().toISOString(),node:process.version,platform:process.platform,playwright:require("playwright/package.json").version,browser:b.version(),baseURL:process.env.BASE_URL}; fs.writeFileSync("../artifacts/browser-qa/runs/"+process.env.QA_RUN_ID+"/environment.json",JSON.stringify(report,null,2)); await b.close()})().catch(e=>{console.error(e);process.exit(1)})'
    if ($LASTEXITCODE -ne 0) { throw 'Browser environment capture failed' }
} finally { Pop-Location }
$env:TEST_DOCKER_FAULTS='1'
try { Check 'infrastructure-regression' { go test -count=1 -timeout 5m -v ./tests -run TestMultiServerAndInfrastructureFaults } }
finally { Remove-Item Env:TEST_DOCKER_FAULTS -ErrorAction SilentlyContinue }
docker compose -p syncframe --profile multi ps --format json | Set-Content -LiteralPath (Join-Path $qaOutput 'services.json') -Encoding utf8
if ($LASTEXITCODE -ne 0) { throw 'Final Compose status capture failed' }
Write-Host "PASS: browser QA and required regressions. Evidence: $qaOutput"
