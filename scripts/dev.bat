@echo off
REM 🔧 Script de desarrollo con hot reload para proyectos UCC
REM Compatible con Windows

setlocal enabledelayedexpansion

REM Verificar si se pidió ayuda
if "%1"=="-h" goto :show_help
if "%1"=="--help" goto :show_help
if "%1"=="/?" goto :show_help

echo ==================================
echo 🔧 UCC - Modo Desarrollo
echo ==================================
echo.

REM Verificar parámetro obligatorio si estamos en directorio raíz
set "CLASS_DIR=%1"

if "%CLASS_DIR%"=="" (
    REM Si estamos en directorio que contiene scripts\ y directorios clase*\
    if exist "scripts\" (
        dir /ad /b clase* >nul 2>&1
        if !errorlevel! equ 0 (
            echo [ERROR] ❌ Parámetro de clase es OBLIGATORIO cuando ejecutas desde el directorio raíz
            echo [ERROR] 
            echo [ERROR] Uso correcto:
            echo [ERROR]   %0 ^<nombre-clase^>
            echo [ERROR] 
            echo [ERROR] Ejemplos:
            echo [ERROR]   %0 clase02-mongo
            echo [ERROR]   %0 clase03-memcache
            echo [ERROR] 
            echo [INFO] Directorios de clases disponibles:
            for /f %%i in ('dir /ad /b clase* 2^>nul') do echo   %%i
            echo [ERROR] 
            echo [ERROR] Alternativa: navega manualmente al directorio
            echo [ERROR]   cd clase02-mongo ^&^& scripts\dev.bat
            pause
            exit /b 1
        )
    )
) else (
    echo [DEV] Navegando al directorio de clase: %CLASS_DIR%
    
    if not exist "%CLASS_DIR%" (
        echo [ERROR] El directorio '%CLASS_DIR%' no existe.
        echo [INFO] Directorios disponibles:
        for /f %%i in ('dir /ad /b clase* 2^>nul') do echo   %%i
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

echo [DEV] Configurando entorno de desarrollo...

goto :main

:show_help
echo 🔧 UCC - Modo Desarrollo
echo ==================================
echo.
echo Uso:
echo   %0                     # Ejecutar en el directorio de la clase
echo   %0 ^<nombre-clase^>      # Ejecutar desde directorio raíz
echo.
echo Ejemplos:
echo   cd clase02-mongo ^&^& %0
echo   %0 clase02-mongo
echo.
pause
exit /b 0

:main

REM Configurar variables de entorno
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

REM Verificar servicios Docker
if exist "docker-compose.yml" (
    echo [DEV] Verificando servicios Docker...
    
    REM Detectar comando Docker Compose
    docker-compose --version >nul 2>&1
    if %errorlevel% equ 0 (
        set "DOCKER_COMPOSE=docker-compose"
    ) else (
        set "DOCKER_COMPOSE=docker compose"
    )
    
    REM Verificar si los servicios están corriendo
    %DOCKER_COMPOSE% ps | findstr "Up" >nul
    if %errorlevel% neq 0 (
        echo [DEV] Iniciando servicios Docker...
        %DOCKER_COMPOSE% up -d
    )
)

REM Verificar si Air está instalado
air --version >nul 2>&1
if %errorlevel% equ 0 (
    echo [DEV] Usando Air para hot reload...
    
    REM Crear configuración de Air si no existe
    if not exist ".air.toml" (
        echo [DEV] Creando configuración de Air...
        
        REM Detectar directorio principal
        set "MAIN_DIR=."
        if exist "cmd\api" (
            set "MAIN_DIR=./cmd/api"
        ) else if exist "api" (
            set "MAIN_DIR=./api"
        )
        
        (
        echo root = "."
        echo testdata_dir = "testdata"
        echo tmp_dir = "tmp"
        echo.
        echo [build]
        echo   args_bin = []
        echo   bin = "./tmp/main.exe"
        echo   cmd = "go build -o ./tmp/main.exe !MAIN_DIR!"
        echo   delay = 1000
        echo   exclude_dir = ["assets", "tmp", "vendor", "testdata", "node_modules"]
        echo   exclude_file = []
        echo   exclude_regex = ["_test.go"]
        echo   exclude_unchanged = false
        echo   follow_symlink = false
        echo   full_bin = ""
        echo   include_dir = []
        echo   include_ext = ["go", "tpl", "tmpl", "html"]
        echo   kill_delay = "0s"
        echo   log = "build-errors.log"
        echo   send_interrupt = false
        echo   stop_on_root = false
        echo.
        echo [color]
        echo   app = ""
        echo   build = "yellow"
        echo   main = "magenta"
        echo   runner = "green"
        echo   watcher = "cyan"
        echo.
        echo [log]
        echo   time = false
        echo.
        echo [misc]
        echo   clean_on_exit = false
        echo.
        echo [screen]
        echo   clear_on_rebuild = false
        ) > .air.toml
        
        echo [SUCCESS] Configuración de Air creada
    )
    
    REM Ejecutar con Air
    air
) else (
    echo [DEV] Air no está instalado. ¿Quieres instalarlo para hot reload? ^(y/N^)
    set /p response="Respuesta: "
    
    if /i "!response!"=="y" (
        echo [DEV] Instalando Air...
        go install github.com/air-verse/air@latest
        if %errorlevel% equ 0 (
            echo [SUCCESS] Air instalado correctamente
            
            REM Agregar GOPATH\bin al PATH si no está
            for /f "tokens=*" %%i in ('go env GOPATH') do set GOPATH=%%i
            set "GOBIN=%GOPATH%\bin"
            echo %PATH% | findstr /C:"%GOBIN%" >nul
            if %errorlevel% neq 0 (
                echo [INFO] Agregando %GOBIN% al PATH...
                set "PATH=%PATH%;%GOBIN%"
            )
            
            call "%~f0"
            exit /b 0
        ) else (
            echo [ERROR] Error instalando Air
        )
    )
    
    REM Fallback a go run normal
    echo [DEV] Ejecutando sin hot reload...
    
    REM Buscar punto de entrada
    set "MAIN_FILE="
    if exist "cmd\api\main.go" (
        set "MAIN_FILE=.\cmd\api"
    ) else if exist "api\main.go" (
        set "MAIN_FILE=.\api"
    ) else if exist "main.go" (
        set "MAIN_FILE=.\main.go"
    ) else (
        echo [ERROR] No se encontró punto de entrada Go
        pause
        exit /b 1
    )
    
    go run %MAIN_FILE%
)

pause