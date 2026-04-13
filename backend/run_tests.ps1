Write-Host "========================================" -ForegroundColor Cyan
Write-Host "         STARTING TEST SUITE" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

$output = & go test -v ./internal/api/handlers -timeout 30s 2>&1
$output | Write-Host

Write-Host ""
Write-Host "========================================" -ForegroundColor Green
Write-Host "         TEST EXECUTION SUMMARY" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green

$lastLine = $output | Select-Object -Last 1

if ($lastLine -like "ok*") {
    Write-Host "✅ All tests succeeded!" -ForegroundColor Green
} else {
    Write-Host "❌ Some tests failed!" -ForegroundColor Red
}

Write-Host "========================================" -ForegroundColor Green
