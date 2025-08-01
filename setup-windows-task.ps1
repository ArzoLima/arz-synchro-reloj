# Script para configurar el scheduler como tarea programada de Windows

param(
    [Parameter(Mandatory=$false)]
    [ValidateSet("install", "uninstall", "status")]
    [string]$Action = "status"
)

$TaskName = "ArzScheduler"
$TaskDescription = "Scheduler de sincronización de marcaciones - Ejecuta todos los días a las 9:00 AM"
$SchedulerPath = Join-Path $PSScriptRoot "scheduler.exe"

function Write-ColorOutput {
    param([string]$Message, [string]$Color = "White")
    Write-Host $Message -ForegroundColor $Color
}

function Test-TaskExists {
    $task = Get-ScheduledTask -TaskName $TaskName -ErrorAction SilentlyContinue
    return $task -ne $null
}

function Install-WindowsTask {
    Write-ColorOutput "Instalando tarea programada..." "Green"
    
    if (-not (Test-Path $SchedulerPath)) {
        Write-ColorOutput "Error: No se encontró scheduler.exe" "Red"
        Write-ColorOutput "Ejecute primero: .\build_and_run.ps1" "Cyan"
        return
    }
    
    if (Test-TaskExists) {
        Write-ColorOutput "La tarea ya existe. Desinstalando primero..." "Yellow"
        Uninstall-WindowsTask
    }
    
    try {
        # Crear la acción (ejecutar el scheduler)
        $action = New-ScheduledTaskAction -Execute $SchedulerPath -WorkingDirectory $PSScriptRoot
        
        # Crear el trigger (todos los días a las 9:00 AM)
        $trigger = New-ScheduledTaskTrigger -Daily -At "9:00 AM"
        
        # Configurar el settings
        $settings = New-ScheduledTaskSettingsSet -AllowStartIfOnBatteries -DontStopIfGoingOnBatteries -StartWhenAvailable -RunOnlyIfNetworkAvailable
        
        # Crear la tarea
        $task = New-ScheduledTask -Action $action -Trigger $trigger -Settings $settings -Description $TaskDescription
        
        # Registrar la tarea
        Register-ScheduledTask -TaskName $TaskName -InputObject $task -User "SYSTEM" -RunLevel Highest
        
        Write-ColorOutput "Tarea programada instalada exitosamente." "Green"
        Write-ColorOutput "La tarea se ejecutará todos los días a las 9:00 AM" "Cyan"
        Write-ColorOutput "Para ver la tarea: Get-ScheduledTask -TaskName '$TaskName'" "Cyan"
        
    } catch {
        Write-ColorOutput "Error al instalar la tarea: $_" "Red"
    }
}

function Uninstall-WindowsTask {
    Write-ColorOutput "Desinstalando tarea programada..." "Yellow"
    
    try {
        Unregister-ScheduledTask -TaskName $TaskName -Confirm:$false
        Write-ColorOutput "Tarea programada desinstalada." "Green"
    } catch {
        Write-ColorOutput "Error al desinstalar la tarea: $_" "Red"
    }
}

function Get-TaskStatus {
    Write-ColorOutput "=== Estado de la Tarea Programada ===" "Cyan"
    
    if (Test-TaskExists) {
        $task = Get-ScheduledTask -TaskName $TaskName
        Write-ColorOutput "Estado: INSTALADA" "Green"
        Write-ColorOutput "Nombre: $($task.TaskName)" "White"
        Write-ColorOutput "Estado actual: $($task.State)" "White"
        Write-ColorOutput "Última ejecución: $($task.LastRunTime)" "White"
        Write-ColorOutput "Próxima ejecución: $($task.NextRunTime)" "White"
        
        # Mostrar triggers
        Write-ColorOutput "Triggers:" "Yellow"
        foreach ($trigger in $task.Triggers) {
            Write-Host "  - $($trigger.CimClass.CimClassName)" "White"
        }
    } else {
        Write-ColorOutput "Estado: NO INSTALADA" "Red"
    }
    
    Write-ColorOutput "=====================================" "Cyan"
}

# Ejecutar acción
switch ($Action.ToLower()) {
    "install" { Install-WindowsTask }
    "uninstall" { Uninstall-WindowsTask }
    "status" { Get-TaskStatus }
    default { 
        Write-ColorOutput "Uso: .\setup-windows-task.ps1 [-Action install|uninstall|status]" "Cyan"
        Write-ColorOutput "Acciones:" "Yellow"
        Write-ColorOutput "  install   - Instalar como tarea programada de Windows" "White"
        Write-ColorOutput "  uninstall - Desinstalar tarea programada" "White"
        Write-ColorOutput "  status    - Ver estado de la tarea" "White"
    }
} 