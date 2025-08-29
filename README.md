# 🎓 UCC - Clases de Programación

> **Template estándar para todas las clases** - Compatible con Windows, Mac y Linux

## 📋 Requisitos del sistema

### Obligatorios
- **Docker** + **Docker Compose** ([Descargar Docker Desktop](https://www.docker.com/products/docker-desktop/))
- **Go 1.22+** ([Descargar Go](https://golang.org/dl/))
- **Git** ([Descargar Git](https://git-scm.com/downloads))

### Verificar instalación
```bash
docker --version
docker-compose --version
go version
git --version
```

## 🚀 Inicio rápido

⚠️ Primero asegúrate de estar en el directorio de la clase. 
Ej.: cd clase02-mongo
Ej.: cd clase03-memcache

**1. Levantar servicios (MongoDB, Memcached, etc.)**
```bash
docker-compose up -d
```

**2. Verificar que los servicios estén listos**
```bash
docker-compose ps
docker-compose logs
```

**3. Configurar variables de entorno**
```bash
# Linux/Mac
cp .env.example .env
export $(grep -v '^#' .env | xargs)

# Windows (PowerShell)
Copy-Item .env.example .env
Get-Content .env | ForEach-Object { if ($_ -match '^([^#].*)=(.*)') { Set-Item -Path "env:$($matches[1])" -Value $matches[2] } }
```

**4. Ejecutar la aplicación Go**
```bash
go run ./cmd/api
```

## 🔧 Scripts disponibles

### Desarrollo diario
```bash
./scripts/start.sh clase02-mongo
./scripts/dev.sh clase02-mongo
```

**Ayuda:**
```bash
./scripts/start.sh --help    # Ver opciones disponibles
./scripts/dev.sh --help      # Ver opciones de desarrollo
```

**💡 Recomendación:** Usa `dev.sh` cuando estés programando y `start.sh` solo para probar rápidamente.
**Archivos observados por Air:**
- ✅ Todos los `.go` en `cmd/`, `internal/`
- ✅ Templates (`.html`, `.tmpl`)
- ❌ Archivos de test (`_test.go`) - ignorados
- ❌ Directorio `tmp/` - ignorados

## 🌐 Endpoints comunes

- **Health Check**: `GET /healthz`
- **Items**: `GET /items`, `POST /items`, `GET /items/:id`
- **API Base**: `http://localhost:8080` (puede variar por clase)

### Ejemplos de uso
```bash
# Verificar salud del servicio
curl http://localhost:8080/healthz

# Listar items
curl http://localhost:8080/items

# Crear nuevo item
curl -X POST http://localhost:8080/items \
  -H "Content-Type: application/json" \
  -d '{"name":"Ejemplo","price":10.99}'
```

## 🐛 Solución de problemas comunes

### Docker no responde
```bash
# Verificar que Docker esté ejecutándose
docker info

# Reiniciar Docker Desktop si es necesario
# Windows/Mac: Reiniciar Docker Desktop desde el systray
# Linux: sudo systemctl restart docker
```

### Puerto ocupado
```bash
# Ver qué proceso usa el puerto
# Linux/Mac:
lsof -i :8080

# Windows:
netstat -ano | findstr :8080

# Cambiar puerto en .env o detener proceso
```

### Base de datos no conecta
```bash
# Verificar logs de la base de datos
docker-compose logs mongo
docker-compose logs memcached

# Reiniciar servicios específicos
docker-compose restart mongo
```

### Errores de Go modules
```bash
# Limpiar cache de módulos
go clean -modcache
go mod download

# Actualizar dependencias
go mod tidy
```

### Errores con Air (hot reload)
```bash
# Error: "module declares its path as: github.com/air-verse/air but was required as: github.com/cosmtrek/air"
# Solución: Air cambió su repositorio

# Instalar manualmente con el nuevo path:
go install github.com/air-verse/air@latest

# Error: "air: command not found" después de la instalación
# Solución: Agregar GOPATH/bin al PATH

# Linux/Mac:
export PATH=$PATH:$(go env GOPATH)/bin

# Windows (PowerShell):
$env:PATH += ";$(go env GOPATH)\bin"

# Windows (CMD):
set PATH=%PATH%;%GOPATH%\bin

# Los scripts ya manejan esto automáticamente
```

### Errores de Docker Build
```bash
# Error: "git": executable file not found in $PATH
# Solución: El Dockerfile necesita instalar git para go mod download

# En el Dockerfile, agregar antes de COPY go.mod:
# RUN apk add --no-cache git

# Reconstruir imagen sin cache
docker-compose build --no-cache

# Ver logs detallados del build
docker-compose build --progress=plain
```

### Permisos en Linux/Mac
```bash
# Dar permisos de ejecución a scripts
chmod +x scripts/*.sh

# Si hay problemas con Docker sin sudo
sudo usermod -aG docker $USER
# Luego reiniciar sesión
```

## 📁 Estructura típica del proyecto

```
proyecto-clase/
├── README.md                 # Este archivo
├── .gitignore               # Archivos a ignorar en Git ⚠️
├── docker-compose.yml       # Definición de servicios
├── .env.example             # Variables de entorno template
├── .env                     # Variables de entorno (no commitear) ⚠️
├── go.mod                   # Dependencias Go
├── scripts/                 # Scripts de automatización
│   ├── start.sh            # Linux/Mac - Iniciar proyecto
│   ├── start.bat           # Windows - Iniciar proyecto
│   ├── dev.sh              # Linux/Mac - Modo desarrollo
│   └── dev.bat             # Windows - Modo desarrollo
├── cmd/api/main.go         # Entrada principal API
├── internal/               # Código interno de la aplicación
│   ├── config/            # Configuración
│   ├── controllers/       # Controladores HTTP
│   ├── services/          # Lógica de negocio
│   ├── repository/        # Acceso a datos
│   └── models/            # Estructuras de datos
├── tmp/                    # Archivos temporales (ignorado) ⚠️
└── init/                  # Scripts de inicialización DB
```

## 💡 Tips para estudiantes

- **🔥 Para DESARROLLO: Usa `./scripts/dev.sh`** - Hot reload automático, cambios instantáneos
- **⚡ Para PRUEBAS: Usa `./scripts/start.sh`** - Ejecución simple una sola vez
- **⚠️ Parámetro de clase es OBLIGATORIO** - `./scripts/dev.sh clase02-mongo`
- **Variables de entorno** - El script copia `.env.example` a `.env` automáticamente
- **Preserva datos** - El modo `dev.sh` mantiene datos en MongoDB entre reinicios
- **Lee los logs** - `docker-compose logs -f` muestra logs en tiempo real
- **Limpieza** - Usa `./scripts/clean.sh` cuando quieras empezar desde cero

## 🆘 ¿Algo no funciona?

1. **⚠️ Verifica que estés en el directorio correcto** - `pwd` debe mostrar `/ruta/clases2025/claseXX-nombre`
2. **Verifica requisitos** - Docker y Go instalados correctamente
3. **Usa los scripts** - Están diseñados para manejar errores comunes  
4. **Lee los logs** - `docker-compose logs` muestra errores detallados
5. **Pregunta al profesor** - Con el error completo y pasos que siguiste

---

**¡Listo para programar! 🚀**
