# Script de PowerShell para compilar y ejecutar el scheduler de sincronización

Write-Host "Compilando programas..." -ForegroundColor Green

# Compilar el programa principal (console)
Write-Host "Compilando programa principal..." -ForegroundColor Yellow
try {
    go build -o console.exe cmd/console/main.go
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Error compilando el programa principal" -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "Error compilando el programa principal: $_" -ForegroundColor Red
    exit 1
}

# Compilar el scheduler
Write-Host "Compilando scheduler..." -ForegroundColor Yellow
try {
    go build -o scheduler.exe cmd/scheduler/main.go
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Error compilando el scheduler" -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "Error compilando el scheduler: $_" -ForegroundColor Red
    exit 1
}

Write-Host "Compilación completada exitosamente" -ForegroundColor Green
Write-Host ""
Write-Host "Para ejecutar el scheduler que correrá todos los días a las 9am:" -ForegroundColor Cyan
Write-Host ".\scheduler.exe" -ForegroundColor White
Write-Host ""
Write-Host "Para ejecutar el programa principal una sola vez:" -ForegroundColor Cyan
Write-Host ".\console.exe" -ForegroundColor White
Write-Host ""
Write-Host "Para ejecutar el scheduler en segundo plano:" -ForegroundColor Cyan
Write-Host "Start-Process -FilePath '.\scheduler.exe' -WindowStyle Hidden" -ForegroundColor White 