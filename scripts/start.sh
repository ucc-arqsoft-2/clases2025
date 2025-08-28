#!/bin/bash

# 🚀 Script de inicio automático para proyectos UCC
# Compatible con Mac y Linux

set -e  # Salir si hay errores

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Funciones de utilidad
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Verificar o navegar al directorio correcto
check_and_navigate_directory() {
    local class_dir="$1"
    
    if [ -n "$class_dir" ]; then
        # Si se pasó un parámetro, intentar navegar a ese directorio
        log_info "Navegando al directorio de clase: $class_dir"
        
        if [ ! -d "$class_dir" ]; then
            log_error "El directorio '$class_dir' no existe."
            log_info "Directorios disponibles:"
            ls -d clase*/ 2>/dev/null || echo "No se encontraron directorios de clase"
            exit 1
        fi
        
        cd "$class_dir" || {
            log_error "No se pudo navegar a '$class_dir'"
            exit 1
        }
        
        log_success "Navegado a: $(basename $(pwd))"
    fi
    
    # Verificar que estamos en el directorio correcto
    log_info "Verificando directorio de trabajo..."
    
    if [ ! -f "docker-compose.yml" ] && [ ! -f "go.mod" ]; then
        if [ -z "$class_dir" ]; then
            log_error "No se encontró docker-compose.yml o go.mod en el directorio actual."
            log_error "Opciones:"
            log_error "1. cd clase02-mongo && ./scripts/start.sh"
            log_error "2. ./scripts/start.sh clase02-mongo (desde el directorio raíz)"
            log_info "Directorio actual: $(pwd)"
            log_info "Directorios disponibles:"
            ls -d */clase*/ 2>/dev/null || ls -d clase*/ 2>/dev/null || echo "No se encontraron directorios de clase"
        else
            log_error "El directorio '$class_dir' no contiene un proyecto válido."
            log_error "Verifica que contenga docker-compose.yml o go.mod"
        fi
        exit 1
    fi
    
    log_success "Directorio correcto confirmado: $(basename $(pwd))"
}

# Verificar dependencias
check_dependencies() {
    log_info "Verificando dependencias..."
    
    # Verificar Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker no está instalado. Instalar desde: https://www.docker.com/products/docker-desktop/"
        exit 1
    fi
    
    # Verificar Docker Compose
    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        log_error "Docker Compose no está disponible. Verificar instalación de Docker."
        exit 1
    fi
    
    # Verificar Go
    if ! command -v go &> /dev/null; then
        log_error "Go no está instalado. Instalar desde: https://golang.org/dl/"
        exit 1
    fi
    
    # Verificar que Docker esté ejecutándose
    if ! docker info &> /dev/null; then
        log_error "Docker no está ejecutándose. Iniciar Docker Desktop o servicio Docker."
        exit 1
    fi
    
    log_success "Todas las dependencias están disponibles"
}

# Configurar variables de entorno
setup_env() {
    log_info "Configurando variables de entorno..."
    
    if [ ! -f ".env" ]; then
        if [ -f ".env.example" ]; then
            cp .env.example .env
            log_success "Archivo .env creado desde .env.example"
        else
            log_warning "No se encontró .env.example, continuando sin variables de entorno específicas"
        fi
    else
        log_info "Archivo .env ya existe"
    fi
    
    # Cargar variables de entorno si existe .env
    if [ -f ".env" ]; then
        export $(grep -v '^#' .env | xargs)
        log_success "Variables de entorno cargadas"
    fi
}

# Verificar puertos disponibles
check_ports() {
    local ports=(8080 27017 11211)
    
    for port in "${ports[@]}"; do
        if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null 2>&1; then
            log_warning "Puerto $port está ocupado, puede haber conflictos"
        fi
    done
}

# Levantar servicios Docker
start_docker_services() {
    log_info "Levantando servicios Docker..."
    
    # Usar docker-compose o docker compose según disponibilidad
    if command -v docker-compose &> /dev/null; then
        DOCKER_COMPOSE="docker-compose"
    else
        DOCKER_COMPOSE="docker compose"
    fi
    
    # Construir y levantar servicios
    if ! $DOCKER_COMPOSE up -d --build; then
        log_error "Error al levantar servicios Docker"
        
        # Verificar si el error es por falta de git
        if $DOCKER_COMPOSE logs | grep -q "git.*executable file not found"; then
            log_error "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
            log_error "PROBLEMA DETECTADO: Falta Git en el contenedor Docker"
            log_error "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
            log_info "SOLUCIÓN:"
            log_info "1. Edita el Dockerfile y agrega esta línea después de FROM:"
            log_info "   RUN apk add --no-cache git"
            log_info ""
            log_info "2. Reconstruye la imagen:"
            log_info "   $DOCKER_COMPOSE build --no-cache"
            log_info ""
            log_info "3. Vuelve a ejecutar este script"
            log_error "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        else
            log_info "Ver logs completos con: $DOCKER_COMPOSE logs"
        fi
        exit 1
    fi
    
    log_success "Servicios Docker iniciados"
    
    # Mostrar estado de contenedores
    log_info "Estado de contenedores:"
    $DOCKER_COMPOSE ps
}

# Esperar a que los servicios estén listos
wait_for_services() {
    log_info "Esperando a que los servicios estén listos..."
    
    # Esperar por MongoDB (puerto 27017)
    if docker-compose ps | grep -q mongo; then
        log_info "Esperando MongoDB..."
        local retries=30
        while [ $retries -gt 0 ]; do
            if docker-compose exec -T mongo mongosh --quiet --eval "db.adminCommand('ping')" &> /dev/null; then
                log_success "MongoDB está listo"
                break
            fi
            retries=$((retries - 1))
            sleep 2
        done
        
        if [ $retries -eq 0 ]; then
            log_error "Timeout esperando MongoDB"
            exit 1
        fi
    fi
    
    # Esperar por Memcached (puerto 11211) si existe
    if docker-compose ps | grep -q memcached; then
        log_info "Esperando Memcached..."
        local retries=15
        while [ $retries -gt 0 ]; do
            if nc -z localhost 11211 2>/dev/null || timeout 1 bash -c 'cat < /dev/null > /dev/tcp/localhost/11211' 2>/dev/null; then
                log_success "Memcached está listo"
                break
            fi
            retries=$((retries - 1))
            sleep 1
        done
    fi
}

# Preparar aplicación Go
prepare_go_app() {
    log_info "Preparando aplicación Go..."
    
    # Verificar que go.mod existe
    if [ ! -f "go.mod" ]; then
        log_error "No se encontró go.mod. ¿Estás en el directorio correcto?"
        exit 1
    fi
    
    # Descargar dependencias
    log_info "Descargando dependencias Go..."
    if ! go mod download; then
        log_error "Error descargando dependencias Go"
        exit 1
    fi
    
    # Limpiar y actualizar módulos
    go mod tidy
    
    log_success "Aplicación Go preparada"
}

# Encontrar y ejecutar el punto de entrada de Go
start_go_app() {
    log_info "Iniciando aplicación Go..."
    
    # Buscar punto de entrada
    local main_file=""
    
    if [ -f "cmd/api/main.go" ]; then
        main_file="./cmd/api"
    elif [ -f "api/main.go" ]; then
        main_file="./api"
    elif [ -f "main.go" ]; then
        main_file="./main.go"
    else
        log_error "No se encontró punto de entrada Go (main.go o cmd/api/main.go)"
        exit 1
    fi
    
    log_info "Ejecutando: go run $main_file"
    log_success "🚀 Aplicación iniciada! Presiona Ctrl+C para detener"
    
    # Función para manejar señales de interrupción
    cleanup() {
        log_info "\nDeteniendo aplicación..."
        if command -v docker-compose &> /dev/null; then
            docker-compose down
        else
            docker compose down
        fi
        log_success "Servicios detenidos"
        exit 0
    }
    
    # Configurar trap para limpieza
    trap cleanup SIGINT SIGTERM
    
    # Ejecutar aplicación Go
    go run $main_file
}

# Mostrar ayuda
show_help() {
    echo "🎓 UCC - Iniciador de Proyectos"
    echo "=================================="
    echo
    echo "Uso:"
    echo "  $0                    # Ejecutar en el directorio de la clase"
    echo "  $0 <nombre-clase>     # Ejecutar desde directorio raíz"
    echo
    echo "Ejemplos:"
    echo "  cd clase02-mongo && $0"
    echo "  $0 clase02-mongo"
    echo "  $0 clase03-memcache"
    echo
    echo "Opciones:"
    echo "  -h, --help           # Mostrar esta ayuda"
    echo
}

# Función principal
main() {
    # Verificar si se pidió ayuda
    if [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
        show_help
        exit 0
    fi
    
    echo "=================================="
    echo "🎓 UCC - Iniciador de Proyectos"
    echo "=================================="
    echo
    
    check_and_navigate_directory "$1"
    check_dependencies
    setup_env
    check_ports
    start_docker_services
    wait_for_services
    prepare_go_app
    start_go_app
}

# Ejecutar función principal
main "$@"