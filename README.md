# Arz Sincronización de Reloj

Este proyecto sincroniza las marcaciones de un reloj con una base de datos MySQL.

## Estructura del Proyecto

- `cmd/console/main.go` - Programa principal que ejecuta la sincronización
- `cmd/scheduler/main.go` - Scheduler que ejecuta el programa todos los días a las 9:00 AM
- `internal/syncmarca/` - Lógica de sincronización de marcaciones

## Configuración

1. Crear un archivo `.env` en la raíz del proyecto con las siguientes variables:

```env
DB_USER=tu_usuario
DB_PASSWORD=tu_password
DB_HOST=tu_host
DB_PORT=3306
DB_NAME=tu_base_de_datos
```

## Compilación y Ejecución

### Linux/macOS (Bash)

#### Opción 1: Usar el script automático

```bash
./build_and_run.sh
```

#### Opción 2: Compilación manual

```bash
# Compilar el programa principal
go build -o console cmd/console/main.go

# Compilar el scheduler
go build -o scheduler cmd/scheduler/main.go
```

### Windows (PowerShell)

#### Opción 1: Usar el script automático

```powershell
.\build_and_run.ps1
```

#### Opción 2: Gestión avanzada con scheduler-manager

```powershell
# Ver estado del scheduler
.\scheduler-manager.ps1 -Action status

# Iniciar scheduler
.\scheduler-manager.ps1 -Action start

# Detener scheduler
.\scheduler-manager.ps1 -Action stop

# Compilar aplicaciones
.\scheduler-manager.ps1 -Action build
```

#### Opción 3: Compilación manual

```powershell
# Compilar el programa principal
go build -o console.exe cmd/console/main.go

# Compilar el scheduler
go build -o scheduler.exe cmd/scheduler/main.go
```

## Uso

### Ejecutar una sola vez

**Linux/macOS:**
```bash
./console
```

**Windows:**
```powershell
.\console.exe
```

### Ejecutar con scheduler (todos los días a las 9:00 AM)

**Linux/macOS:**
```bash
./scheduler
```

**Windows:**
```powershell
.\scheduler.exe
```

### Gestión del scheduler en Windows

El script `scheduler-manager.ps1` proporciona funciones avanzadas:

- **Ver estado**: `.\scheduler-manager.ps1 -Action status`
- **Iniciar**: `.\scheduler-manager.ps1 -Action start`
- **Detener**: `.\scheduler-manager.ps1 -Action stop`
- **Compilar**: `.\scheduler-manager.ps1 -Action build`

## Características del Scheduler

El scheduler:
- Se ejecuta en un bucle infinito
- Calcula automáticamente el próximo tiempo de ejecución (9:00 AM)
- Muestra el tiempo restante hasta la próxima ejecución
- Ejecuta el programa principal y muestra los resultados
- Maneja errores sin detener el scheduler

## Logs

El scheduler mostrará:
- Cuándo se inició
- La próxima hora de ejecución programada
- El tiempo restante hasta la próxima ejecución
- El resultado de cada ejecución
- Cualquier error que ocurra

**En Windows**, los logs se guardan en `scheduler.log` y se pueden ver con:
```powershell
Get-Content scheduler.log -Tail 10
```

## Notas

- El scheduler debe ejecutarse en segundo plano para funcionar continuamente
- **Linux/macOS**: Puedes usar `nohup` o `systemd` para mantenerlo ejecutándose
- **Windows**: Usa `Start-Process` con `-WindowStyle Hidden` para ejecutar en segundo plano
- El programa principal procesa las marcaciones del mes y año actual
- Se eliminan registros existentes antes de insertar nuevos para evitar duplicados

## Ejecución en Segundo Plano

### Linux/macOS
```bash
nohup ./scheduler > scheduler.log 2>&1 &
```

### Windows
```powershell
Start-Process -FilePath ".\scheduler.exe" -WindowStyle Hidden
```

## Configuración de Tareas Programadas (Windows)

Para una ejecución automática más robusta en Windows, puedes configurar el scheduler como una tarea programada del sistema:

### Instalar como Tarea Programada

```powershell
# Instalar la tarea programada
.\setup-windows-task.ps1 -Action install

# Verificar el estado
.\setup-windows-task.ps1 -Action status

# Desinstalar si es necesario
.\setup-windows-task.ps1 -Action uninstall
```

### Ventajas de la Tarea Programada

- **Ejecución automática**: Se ejecuta automáticamente al iniciar Windows
- **Gestión del sistema**: Windows maneja el ciclo de vida del proceso
- **Logs del sistema**: Los logs se integran con el Event Viewer de Windows
- **Configuración avanzada**: Permite configurar condiciones de ejecución
- **Recuperación automática**: Windows reinicia la tarea si falla

### Configuración de la Tarea

La tarea programada se configura con:
- **Trigger**: Diario a las 9:00 AM
- **Usuario**: SYSTEM (permisos elevados)
- **Condiciones**: Se ejecuta solo si hay conexión de red
- **Configuración**: Permite ejecución en batería y inicio automático 