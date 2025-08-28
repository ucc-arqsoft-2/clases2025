@echo off
REM 🛑 Script para detener servicios UCC
REM Compatible con Windows

echo ==================================
echo 🛑 UCC - Detener Servicios
echo ==================================
echo.

if exist "docker-compose.yml" (
    echo [INFO] Deteniendo servicios Docker...
    
    REM Detectar comando Docker Compose
    docker-compose --version >nul 2>&1
    if %errorlevel% equ 0 (
        docker-compose down
    ) else (
        docker compose down
    )
    
    echo [SUCCESS] Servicios detenidos correctamente
) else (
    echo [INFO] No se encontró docker-compose.yml
)

pause