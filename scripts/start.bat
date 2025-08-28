@echo off
REM 🚀 Script de inicio automático para proyectos UCC
REM Compatible con Windows

setlocal enabledelayedexpansion

REM Verificar si se pidió ayuda
if "%1"=="-h" goto :show_help
if "%1"=="--help" goto :show_help
if "%1"=="/?" goto :show_help

echo ==================================
echo 🎓 UCC - Iniciador de Proyectos
echo ==================================
echo.

REM Verificar o navegar al directorio correcto
set "CLASS_DIR=%1"

if not "%CLASS_DIR%"=="" (
    echo [INFO] Navegando al directorio de clase: %CLASS_DIR%
    
    if not exist "%CLASS_DIR%" (
        echo [ERROR] El directorio '%CLASS_DIR%' no existe.
        echo [INFO] Directorios disponibles:
        dir /ad /b clase* 2>nul || echo No se encontraron directorios de clase
        pause
        exit /b 1
    )
    
    cd "%CLASS_DIR%" || (
        echo [ERROR] No se pudo navegar a '%CLASS_DIR%'
        pause
        exit /b 1
    )
    
    for %%i in (.) do echo [SUCCESS] Navegado a: %%~nxi
)

REM Verificar que estamos en el directorio correcto
echo [INFO] Verificando directorio de trabajo...

if not exist "docker-compose.yml" if not exist "go.mod" (
    if "%CLASS_DIR%"=="" (
        echo [ERROR] No se encontró docker-compose.yml o go.mod en el directorio actual.
        echo [ERROR] Opciones:
        echo [ERROR] 1. cd clase02-mongo ^&^& scripts\start.bat
        echo [ERROR] 2. scripts\start.bat clase02-mongo ^(desde el directorio raíz^)
        echo [INFO] Directorio actual: %CD%
        echo [INFO] Directorios disponibles:
        dir /ad /b clase* 2>nul || echo No se encontraron directorios de clase
    ) else (
        echo [ERROR] El directorio '%CLASS_DIR%' no contiene un proyecto válido.
        echo [ERROR] Verifica que contenga docker-compose.yml o go.mod
    )
    pause
    exit /b 1
)

for %%i in (.) do echo [SUCCESS] Directorio correcto confirmado: %%~nxi

goto :main

:show_help
echo 🎓 UCC - Iniciador de Proyectos
echo ==================================
echo.
echo Uso:
echo   %0                     # Ejecutar en el directorio de la clase
echo   %0 ^<nombre-clase^>      # Ejecutar desde directorio raíz
echo.
echo Ejemplos:
echo   cd clase02-mongo ^&^& %0
echo   %0 clase02-mongo
echo   %0 clase03-memcache
echo.
echo Opciones:
echo   -h, --help, /?         # Mostrar esta ayuda
echo.
pause
exit /b 0

:main

echo [INFO] Verificando dependencias...

REM Verificar Docker
docker --version >nul 2>&1
if %errorlevel% neq 0 (
    echo [ERROR] Docker no está instalado. Instalar desde: https://www.docker.com/products/docker-desktop/
    pause
    exit /b 1
)

REM Verificar Docker Compose
docker-compose --version >nul 2>&1 || docker compose version >nul 2>&1
if %errorlevel% neq 0 (
    echo [ERROR] Docker Compose no está disponible. Verificar instalación de Docker.
    pause
    exit /b 1
)

REM Verificar Go
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo [ERROR] Go no está instalado. Instalar desde: https://golang.org/dl/
    pause
    exit /b 1
)

REM Verificar que Docker esté ejecutándose
docker info >nul 2>&1
if %errorlevel% neq 0 (
    echo [ERROR] Docker no está ejecutándose. Iniciar Docker Desktop.
    pause
    exit /b 1
)

echo [SUCCESS] Todas las dependencias están disponibles

REM Configurar variables de entorno
echo [INFO] Configurando variables de entorno...

if not exist ".env" (
    if exist ".env.example" (
        copy ".env.example" ".env" >nul
        echo [SUCCESS] Archivo .env creado desde .env.example
    ) else (
        echo [WARNING] No se encontró .env.example, continuando sin variables específicas
    )
) else (
    echo [INFO] Archivo .env ya existe
)

REM Cargar variables de entorno si existe .env
if exist ".env" (
    for /f "usebackq tokens=*" %%a in (".env") do (
        echo %%a | findstr /v "^#" | findstr "=" >nul
        if !errorlevel! equ 0 (
            for /f "tokens=1,2 delims==" %%b in ("%%a") do (
                set "%%b=%%c"
            )
        )
    )
    echo [SUCCESS] Variables de entorno cargadas
)

REM Verificar puertos comunes
echo [INFO] Verificando puertos...
netstat -an | findstr ":8080 " | findstr "LISTENING" >nul && echo [WARNING] Puerto 8080 está ocupado
netstat -an | findstr ":27017 " | findstr "LISTENING" >nul && echo [WARNING] Puerto 27017 está ocupado
netstat -an | findstr ":11211 " | findstr "LISTENING" >nul && echo [WARNING] Puerto 11211 está ocupado

REM Detectar Docker Compose command
docker-compose --version >nul 2>&1
if %errorlevel% equ 0 (
    set "DOCKER_COMPOSE=docker-compose"
) else (
    set "DOCKER_COMPOSE=docker compose"
)

REM Levantar servicios Docker
echo [INFO] Levantando servicios Docker...
%DOCKER_COMPOSE% up -d --build
if %errorlevel% neq 0 (
    echo [ERROR] Error al levantar servicios Docker
    pause
    exit /b 1
)

echo [SUCCESS] Servicios Docker iniciados

REM Mostrar estado de contenedores
echo [INFO] Estado de contenedores:
%DOCKER_COMPOSE% ps

REM Esperar a que MongoDB esté listo
echo [INFO] Esperando a que los servicios estén listos...
%DOCKER_COMPOSE% ps | findstr mongo >nul
if %errorlevel% equ 0 (
    echo [INFO] Esperando MongoDB...
    set /a retries=30
    :wait_mongo
    if !retries! gtr 0 (
        %DOCKER_COMPOSE% exec -T mongo mongosh --quiet --eval "db.adminCommand('ping')" >nul 2>&1
        if !errorlevel! equ 0 (
            echo [SUCCESS] MongoDB está listo
            goto mongo_ready
        )
        set /a retries-=1
        timeout /t 2 /nobreak >nul
        goto wait_mongo
    ) else (
        echo [ERROR] Timeout esperando MongoDB
        pause
        exit /b 1
    )
    :mongo_ready
)

REM Preparar aplicación Go
echo [INFO] Preparando aplicación Go...

if not exist "go.mod" (
    echo [ERROR] No se encontró go.mod. ¿Estás en el directorio correcto?
    pause
    exit /b 1
)

echo [INFO] Descargando dependencias Go...
go mod download
if %errorlevel% neq 0 (
    echo [ERROR] Error descargando dependencias Go
    pause
    exit /b 1
)

go mod tidy
echo [SUCCESS] Aplicación Go preparada

REM Encontrar y ejecutar punto de entrada
echo [INFO] Iniciando aplicación Go...

set "MAIN_FILE="
if exist "cmd\api\main.go" (
    set "MAIN_FILE=.\cmd\api"
) else if exist "api\main.go" (
    set "MAIN_FILE=.\api"
) else if exist "main.go" (
    set "MAIN_FILE=.\main.go"
) else (
    echo [ERROR] No se encontró punto de entrada Go ^(main.go o cmd\api\main.go^)
    pause
    exit /b 1
)

echo [INFO] Ejecutando: go run %MAIN_FILE%
echo [SUCCESS] 🚀 Aplicación iniciada! Presiona Ctrl+C para detener

REM Configurar manejo de señales de interrupción
REM En Windows, cuando se presiona Ctrl+C, el batch se detiene automáticamente

REM Ejecutar aplicación Go
go run %MAIN_FILE%

REM Si llegamos aquí, la aplicación se cerró
echo.
echo [INFO] Deteniendo servicios...
%DOCKER_COMPOSE% down
echo [SUCCESS] Servicios detenidos

pause