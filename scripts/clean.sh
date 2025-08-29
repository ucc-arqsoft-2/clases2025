#!/bin/bash

# 🧹 Script para limpiar contenedores y datos UCC
# Compatible con Mac y Linux

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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
    
    # Si no se pasó parámetro, verificar si estamos en directorio raíz del proyecto
    if [ -z "$class_dir" ]; then
        # Si estamos en directorio que contiene scripts/ y directorios clase*/
        if [ -d "scripts" ] && ls -d clase*/ >/dev/null 2>&1; then
            log_error "❌ Parámetro de clase es OBLIGATORIO cuando ejecutas desde el directorio raíz"
            log_error ""
            log_error "Uso correcto:"
            log_error "  $0 <nombre-clase>"
            log_error ""
            log_error "Ejemplos:"
            log_error "  $0 clase02-mongo"
            log_error "  $0 clase03-memcache"
            log_error ""
            log_info "Directorios de clases disponibles:"
            ls -d clase*/ 2>/dev/null | sed 's|/||g' | sed 's/^/  /'
            log_error ""
            log_error "Alternativa: navega manualmente al directorio"
            log_error "  cd clase02-mongo && ./scripts/clean.sh"
            exit 1
        fi
    else
        # Si se pasó un parámetro, intentar navegar a ese directorio
        log_info "Navegando al directorio de clase: $class_dir"
        
        if [ ! -d "$class_dir" ]; then
            log_error "El directorio '$class_dir' no existe."
            log_info "Directorios disponibles:"
            ls -d clase*/ 2>/dev/null | sed 's|/||g' | sed 's/^/  /' || echo "  No se encontraron directorios de clase"
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
        log_error "El directorio actual no contiene un proyecto válido."
        log_error "Debe contener docker-compose.yml o go.mod"
        log_info "Directorio actual: $(pwd)"
        exit 1
    fi
    
    log_success "Directorio correcto confirmado: $(basename $(pwd))"
}

# Mostrar ayuda
show_help() {
    echo "🧹 UCC - Limpieza de Proyecto"
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

main() {
    # Verificar si se pidió ayuda
    if [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
        show_help
        exit 0
    fi
    
    echo "=================================="
    echo "🧹 UCC - Limpieza de Proyecto"
    echo "=================================="
    echo
    
    check_and_navigate_directory "$1"
    
    if [ -f "docker-compose.yml" ]; then
        log_warning "Esto eliminará TODOS los contenedores y datos del proyecto."
        echo "¿Estás seguro? (y/N)"
        read -r response
        
        if [[ "$response" =~ ^[Yy]$ ]]; then
            log_info "Deteniendo y eliminando contenedores..."
            
            if command -v docker-compose &> /dev/null; then
                docker-compose down -v --remove-orphans
            else
                docker compose down -v --remove-orphans
            fi
            
            log_info "Eliminando imágenes del proyecto..."
            docker image prune -f --filter label=com.docker.compose.project
            
            log_info "Limpiando archivos temporales..."
            rm -rf tmp/
            rm -f build-errors.log
            
            log_success "Proyecto limpiado correctamente"
        else
            log_info "Operación cancelada"
        fi
    else
        log_info "No se encontró docker-compose.yml"
    fi
}

main "$@"