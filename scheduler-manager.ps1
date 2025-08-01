# Script de PowerShell para gestionar el scheduler de sincronización

param(
    [Parameter(Mandatory=$false)]
    [ValidateSet("start", "stop", "status", "build")]
    [string]$Action = "status"
)

$SchedulerPath = Join-Path $PSScriptRoot "scheduler.exe"
$LogPath = Join-Path $PSScriptRoot "scheduler.log"

function Write-ColorOutput {
    param([string]$Message, [string]$Color = "White")
    Write-Host $Message -ForegroundColor $Color
}

function Test-SchedulerRunning {
    $process = Get-Process -Name "scheduler" -ErrorAction SilentlyContinue
    return $process -ne $null
}

function Start-SchedulerService {
    Write-ColorOutput "Iniciando scheduler..." "Green"
    
    if (Test-SchedulerRunning) {
        Write-ColorOutput "El scheduler ya está ejecutándose." "Yellow"
        return
    }
    
    if (-not (Test-Path $SchedulerPath)) {
        Write-ColorOutput "Error: No se encontró scheduler.exe" "Red"
        Write-ColorOutput "Ejecute primero: .\build_and_run.ps1" "Cyan"
        return
    }
    
    try {
        Start-Process -FilePath $SchedulerPath -WindowStyle Hidden -RedirectStandardOutput $LogPath -RedirectStandardError $LogPath
        Start-Sleep -Seconds 2
        
        if (Test-SchedulerRunning) {
            Write-ColorOutput "Scheduler iniciado exitosamente." "Green"
            Write-ColorOutput "Logs en: $LogPath" "Cyan"
        } else {
            Write-ColorOutput "Error al iniciar el scheduler." "Red"
        }
    } catch {
        Write-ColorOutput "Error: $_" "Red"
    }
}

function Stop-SchedulerService {
    Write-ColorOutput "Deteniendo scheduler..." "Yellow"
    
    $process = Get-Process -Name "scheduler" -ErrorAction SilentlyContinue
    if ($process) {
        try {
            $process | Stop-Process -Force
            Write-ColorOutput "Scheduler detenido." "Green"
        } catch {
            Write-ColorOutput "Error al detener: $_" "Red"
        }
    } else {
        Write-ColorOutput "El scheduler no está ejecutándose." "Yellow"
    }
}

function Get-SchedulerStatus {
    Write-ColorOutput "=== Estado del Scheduler ===" "Cyan"
    
    if (Test-SchedulerRunning) {
        Write-ColorOutput "Estado: EJECUTÁNDOSE" "Green"
        $process = Get-Process -Name "scheduler"
        Write-ColorOutput "PID: $($process.Id)" "White"
    } else {
        Write-ColorOutput "Estado: DETENIDO" "Red"
    }
    
    if (Test-Path $LogPath) {
        Write-ColorOutput "Últimas líneas del log:" "Yellow"
        Get-Content $LogPath -Tail 5 | ForEach-Object { Write-Host "  $_" }
    }
}

function Build-Applications {
    Write-ColorOutput "Compilando aplicaciones..." "Green"
    
    try {
        go build -o console.exe cmd/console/main.go
        go build -o scheduler.exe cmd/scheduler/main.go
        Write-ColorOutput "Compilación completada." "Green"
    } catch {
        Write-ColorOutput "Error en compilación: $_" "Red"
    }
}

# Ejecutar acción
switch ($Action.ToLower()) {
    "start" { Start-SchedulerService }
    "stop" { Stop-SchedulerService }
    "status" { Get-SchedulerStatus }
    "build" { Build-Applications }
    default { 
        Write-ColorOutput "Uso: .\scheduler-manager.ps1 [-Action start|stop|status|build]" "Cyan"
    }
} 