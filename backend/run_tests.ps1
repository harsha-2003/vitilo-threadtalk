Write-Host "========================================" -ForegroundColor Cyan
Write-Host "         STARTING BACKEND TEST SUITE" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

$env:GOCACHE = Join-Path (Get-Location) ".gocache"
$output = & go test -v ./internal/api/handlers -timeout 60s 2>&1
$exitCode = $LASTEXITCODE
$output | Write-Host

Write-Host ""
Write-Host "========================================" -ForegroundColor Green
Write-Host "         TEST EXECUTION SUMMARY" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green

if ($exitCode -eq 0) {
    Write-Host "All backend tests succeeded." -ForegroundColor Green
} else {
    Write-Host "Some backend tests failed." -ForegroundColor Red
}

Write-Host "Command: go test -v ./internal/api/handlers -timeout 60s"
Write-Host "========================================" -ForegroundColor Green
exit $exitCode
