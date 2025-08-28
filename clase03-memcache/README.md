# 🧪 Proyecto Base (Starter): Gin + MongoDB (+ Memcached a integrar)

Este es el **proyecto base** para la clase. Ya trae:
- API en Gin
- Conexión a MongoDB
- Endpoints de `items` funcionando **sin cache**
- Estructura de archivos para integrar **Memcached** (con TODOs)

## 🎯 Objetivo del laboratorio
Integrar **Memcached** con patrón **cache-first** y **invalidación**:

- Clave de lista: `items:all`
- Clave de item: `item:<id>`
- Endpoints:
  - `GET /v1/items` → cache-first sobre `items:all`
  - `POST /v1/items` → cachear `item:<id>` e **invalidar** `items:all`
  - `GET /v1/items/:id` → cache-first sobre `item:<id>`
  - `/healthz` → debe realizar una escritura/lectura efímera en Memcached

## ▶️ Levantar el proyecto base
```bash
docker compose up --build
```

El API queda expuesto en `http://localhost:8081` (para coexistir con la solución).

Probar (sin cache aún):
```bash
curl -s localhost:8081/healthz | jq
curl -s -X POST localhost:8081/v1/items -H "Content-Type: application/json" -d '{"name":"Demo","price":10}' | jq
curl -s localhost:8081/v1/items | jq
```

## 🧩 Archivos clave con TODOs
- `internal/cache/memcached.go` → implementar wrapper de Memcached (GetJSON, SetJSON, Delete, SelfTest).
- `cmd/api/main.go` → leer `MEMCACHED_ADDR`, `CACHE_TTL_SECONDS`, crear cliente y pasar al router.
- `internal/server/server.go` → usar `c.SelfTest` en `/healthz`.
- `internal/handlers/items.go` → implementar:
  - cache-first en `List` y `GetByID`
  - set cache en `Create` e invalidar lista

## 📝 Consignas y actividades
1. **Configurar cliente Memcached**:
   - Crear `cache.New(memcachedAddr, ttlDur)` con TTL leído de `CACHE_TTL_SECONDS`.
   - Inyectarlo en el router y en el `ItemHandler`.

2. **Implementar cache-first**:
   - `List`: probar `GetJSON("items:all", &items)`; si miss, leer store y `SetJSON`.
   - `GetByID`: probar `GetJSON("item:"+id, &item)`; si miss, leer store y `SetJSON`.

3. **Invalidación**:
   - En `Create`, además de `SetJSON("item:<id>", item)`, invocar `Delete("items:all")`.

4. **Healthcheck real**:
   - `/healthz`: usar `c.SelfTest(ctx)` (set/get efímero).

5. **Experimentos**:
   - Bajar `CACHE_TTL_SECONDS` a 5 y observar expiración.
   - Agregar `DELETE /v1/items/:id` y ajustar invalidación.

## 💡 Sugerencia de pruebas
```bash
# Crear 2 items
curl -s -X POST localhost:8081/v1/items -H "Content-Type: application/json" -d '{"name":"Coca-Cola 350ml","price":123.45}' | jq
curl -s -X POST localhost:8081/v1/items -H "Content-Type: application/json" -d '{"name":"Sprite 500ml","price":99.99}' | jq

# Listar dos veces (esperado: store → cache)
curl -s localhost:8081/v1/items | jq
curl -s localhost:8081/v1/items | jq

# Obtener por ID dos veces (esperado: store → cache)
ID="<reemplazar>"
curl -s localhost:8081/v1/items/$ID | jq
curl -s localhost:8081/v1/items/$ID | jq
```

## 🧹 Limpieza
```bash
docker compose down -v
```
