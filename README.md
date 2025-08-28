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

Tienes **dos opciones** para ejecutar el proyecto:

### Opción 1: Automática con parámetro (Más fácil) 🌟

**Desde el directorio raíz del repo**, pasar el nombre de la clase:

**Linux/Mac:**
```bash
cd clases2025                    # Directorio raíz del repo
chmod +x scripts/start.sh       # Solo la primera vez
./scripts/start.sh clase02-mongo
```

**Windows:**
```cmd
cd clases2025                    REM Directorio raíz del repo
scripts\start.bat clase02-mongo
```

### Opción 2: Comandos manuales (Para aprender)

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

**Todos los scripts soportan ambos modos:**
```bash
# Opción 1: Con parámetro (desde directorio raíz)
./scripts/start.sh clase02-mongo
./scripts/dev.sh clase02-mongo

# Opción 2: Sin parámetro (desde directorio de clase)
cd clase02-mongo && ./scripts/start.sh
cd clase02-mongo && ./scripts/dev.sh
```

**Scripts disponibles:**
- `scripts/start.sh` / `scripts/start.bat` - Inicia todo el proyecto automáticamente
- `scripts/dev.sh` / `scripts/dev.bat` - Modo desarrollo con hot reload
- `scripts/stop.sh` / `scripts/stop.bat` - Detiene todos los servicios
- `scripts/clean.sh` / `scripts/clean.bat` - Limpia contenedores y datos

**Ayuda:**
```bash
./scripts/start.sh --help    # Ver opciones disponibles
./scripts/dev.sh --help      # Ver opciones de desarrollo
```

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
├── docker-compose.yml       # Definición de servicios
├── .env.example             # Variables de entorno template
├── .env                     # Variables de entorno (no commitear)
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
└── init/                  # Scripts de inicialización DB
```

## 🎯 Flujo de trabajo recomendado

1. **Clonar el repositorio**
   ```bash
   git clone [URL_DEL_REPO]
   cd [NOMBRE_PROYECTO]
   ```

2. **Ejecutar proyecto (elige tu opción preferida)**
   
   **Opción A - Automática (Recomendada):**
   ```bash
   ./scripts/start.sh clase02-mongo    # Linux/Mac
   scripts\start.bat clase02-mongo     # Windows
   ```
   
   **Opción B - Manual:**
   ```bash
   cd clase02-mongo                    # Navegar a la clase
   ./scripts/start.sh                  # Linux/Mac
   scripts\start.bat                   # Windows
   ```

3. **Desarrollar y probar**
   - Código en `internal/`
   - Probar endpoints con curl o Postman
   - Ver logs: `docker-compose logs -f`

4. **Detener servicios al terminar**
   ```bash
   docker-compose down
   ```

## 💡 Tips para estudiantes

- **⚠️ SIEMPRE hacer `cd` al directorio de la clase primero** - Es el error más común
- **Usa los scripts automatizados** - Evitan errores comunes
- **Lee los logs** - `docker-compose logs` te dice qué está pasando
- **Variables de entorno** - Siempre copia `.env.example` a `.env`
- **Hot reload** - Usa `./scripts/dev.sh` para development
- **Limpieza** - Ejecuta `docker-compose down -v` para limpiar datos de prueba

## 🆘 ¿Algo no funciona?

1. **⚠️ Verifica que estés en el directorio correcto** - `pwd` debe mostrar `/ruta/clases2025/claseXX-nombre`
2. **Verifica requisitos** - Docker y Go instalados correctamente
3. **Usa los scripts** - Están diseñados para manejar errores comunes  
4. **Lee los logs** - `docker-compose logs` muestra errores detallados
5. **Pregunta al profesor** - Con el error completo y pasos que siguiste

---

**¡Listo para programar! 🚀**