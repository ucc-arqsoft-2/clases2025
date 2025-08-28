#!/bin/bash

# 🔧 Script de desarrollo con hot reload para proyectos UCC
# Compatible con Mac y Linux

set -e

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${BLUE}[DEV]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Verificar si Air está instalado (para hot reload)
install_air() {
    if ! command -v air &> /dev/null; then
        log_info "Instalando Air para hot reload..."
        if ! go install github.com/cosmtrek/air@latest; then
            log_error "Error instalando Air"
            exit 1
        fi
        log_success "Air instalado correctamente"
    fi
}

# Crear configuración de Air si no existe
create_air_config() {
    if [ ! -f ".air.toml" ]; then
        log_info "Creando configuración de Air..."
        
        # Detectar directorio principal
        local main_dir="."
        if [ -d "cmd/api" ]; then
            main_dir="./cmd/api"
        elif [ -d "api" ]; then
            main_dir="./api"
        fi
        
        cat > .air.toml << EOF
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main $main_dir"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata", "node_modules"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  kill_delay = "0s"
  log = "build-errors.log"
  send_interrupt = false
  stop_on_root = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
EOF
        log_success "Configuración de Air creada"
    fi
}

# Función principal de desarrollo
main() {
    # Verificar si se pidió ayuda
    if [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
        echo "🔧 UCC - Modo Desarrollo"
        echo "=================================="
        echo
        echo "Uso:"
        echo "  $0                    # Ejecutar en el directorio de la clase"
        echo "  $0 <nombre-clase>     # Ejecutar desde directorio raíz"
        echo
        echo "Ejemplos:"
        echo "  cd clase02-mongo && $0"
        echo "  $0 clase02-mongo"
        echo
        exit 0
    fi
    
    echo "=================================="
    echo "🔧 UCC - Modo Desarrollo"
    echo "=================================="
    echo
    
    # Navegar al directorio si se especificó
    if [ -n "$1" ]; then
        log_info "Navegando al directorio de clase: $1"
        if [ ! -d "$1" ]; then
            log_error "El directorio '$1' no existe."
            exit 1
        fi
        cd "$1" || {
            log_error "No se pudo navegar a '$1'"
            exit 1
        }
        log_success "Navegado a: $(basename $(pwd))"
    fi
    
    log_info "Configurando entorno de desarrollo..."
    
    # Configurar variables de entorno
    if [ -f ".env" ]; then
        export $(grep -v '^#' .env | xargs)
        log_success "Variables de entorno cargadas"
    fi
    
    # Verificar servicios Docker
    if [ -f "docker-compose.yml" ]; then
        log_info "Verificando servicios Docker..."
        if command -v docker-compose &> /dev/null; then
            DOCKER_COMPOSE="docker-compose"
        else
            DOCKER_COMPOSE="docker compose"
        fi
        
        # Levantar servicios si no están corriendo
        if ! $DOCKER_COMPOSE ps | grep -q "Up"; then
            log_info "Iniciando servicios Docker..."
            $DOCKER_COMPOSE up -d
        fi
    fi
    
    # Configurar hot reload
    if command -v air &> /dev/null; then
        log_info "Usando Air para hot reload..."
        create_air_config
        air
    else
        log_info "Air no está instalado. ¿Quieres instalarlo para hot reload? (y/N)"
        read -r response
        if [[ "$response" =~ ^[Yy]$ ]]; then
            install_air
            create_air_config
            air
        else
            # Fallback a go run normal
            log_info "Ejecutando sin hot reload..."
            
            # Buscar punto de entrada
            local main_file=""
            if [ -f "cmd/api/main.go" ]; then
                main_file="./cmd/api"
            elif [ -f "api/main.go" ]; then
                main_file="./api"
            elif [ -f "main.go" ]; then
                main_file="./main.go"
            else
                log_error "No se encontró punto de entrada Go"
                exit 1
            fi
            
            go run $main_file
        fi
    fi
}

main "$@"