# Clase 04 — RabbitMQ (Completar)

Estructura igual a la solución (respetando el screenshot). Este repo está
preparado para que completes los ejercicios de cacheo con Memcached.

## Requisitos
- Docker Desktop (Windows / macOS)
- `git`

## Levantar
1. Copiá `.env.example` a `.env`
2. `docker compose up --build`
3. Probar:
   - macOS/Linux:
     ```bash
     curl -s http://localhost:8080/healthz | jq .
     curl -s http://localhost:8080/items | jq .
     ```
   - Windows PowerShell:
     ```powershell
     Invoke-RestMethod http://localhost:8080/healthz
     Invoke-RestMethod http://localhost:8080/items
     ```

## Ejercicios
1. **Listado cacheado**: en `internal/service/items.go` almacenar el resultado de `List()`
   en Memcached con clave `items:all` y TTL 60s. Leer desde cache si existe.
2. **Detalle cacheado**: en `Get()` cachear bajo `item:<id>`.
3. **Invalidación**: al `Create`, `Update`, `Delete` invalidar `items:all` y `item:<id>`.
4. **Índice de claves**: en `internal/cache/memcached.go` mantener un índice de claves
   (por ejemplo, clave `cache:index` con un JSON de strings) para poder listarlas desde
   `/__cache/keys`. Remover del índice en `Delete`.
5. **Endpoints de inspección**: hacer que `/__cache/keys` y `/__cache/get?key=` funcionen
   usando tu implementación.

## Ver la cache desde tu PC
Cuando completes el punto 4, podrás:
```bash
curl -s http://localhost:8080/__cache/keys | jq .
curl -s "http://localhost:8080/__cache/get?key=items:all" | jq .
```

## Notas
- Memcached está expuesto en el puerto 11211 del host para que puedas probar herramientas externas.
- Mongo se inicializa con `mongo-init/seed.js`.

## Testear CREATE
```bash
curl -X POST http://localhost:8080/items \
  -H "Content-Type: application/json" \
  -d '{"id":"item-123","name":"Coca Cola","price":100}'
```

##
docker compose down --remove-orphans
docker compose build --no-cache api
docker compose up -d
docker compose logs -f api
