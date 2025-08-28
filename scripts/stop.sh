#!/bin/bash

# 🛑 Script para detener servicios UCC
# Compatible con Mac y Linux

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

main() {
    echo "=================================="
    echo "🛑 UCC - Detener Servicios"
    echo "=================================="
    echo
    
    if [ -f "docker-compose.yml" ]; then
        log_info "Deteniendo servicios Docker..."
        
        if command -v docker-compose &> /dev/null; then
            docker-compose down
        else
            docker compose down
        fi
        
        log_success "Servicios detenidos correctamente"
    else
        log_info "No se encontró docker-compose.yml"
    fi
}

main "$@"