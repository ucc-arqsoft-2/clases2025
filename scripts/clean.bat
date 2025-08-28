@echo off
REM 🧹 Script para limpiar contenedores y datos UCC
REM Compatible con Windows

echo ==================================
echo 🧹 UCC - Limpieza de Proyecto
echo ==================================
echo.

if exist "docker-compose.yml" (
    echo [WARNING] Esto eliminará TODOS los contenedores y datos del proyecto.
    set /p response="¿Estás seguro? (y/N): "
    
    if /i "!response!"=="y" (
        echo [INFO] Deteniendo y eliminando contenedores...
        
        REM Detectar comando Docker Compose
        docker-compose --version >nul 2>&1
        if %errorlevel% equ 0 (
            docker-compose down -v --remove-orphans
        ) else (
            docker compose down -v --remove-orphans
        )
        
        echo [INFO] Eliminando imágenes del proyecto...
        docker image prune -f --filter label=com.docker.compose.project
        
        echo [INFO] Limpiando archivos temporales...
        if exist "tmp\" rmdir /s /q "tmp\"
        if exist "build-errors.log" del "build-errors.log"
        
        echo [SUCCESS] Proyecto limpiado correctamente
    ) else (
        echo [INFO] Operación cancelada
    )
) else (
    echo [INFO] No se encontró docker-compose.yml
)

pause