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

main() {
    echo "=================================="
    echo "🧹 UCC - Limpieza de Proyecto"
    echo "=================================="
    echo
    
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