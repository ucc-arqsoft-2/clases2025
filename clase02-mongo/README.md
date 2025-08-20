# 🧪 Proyecto Base – API Go (Gin) + MongoDB + Docker

Este proyecto base está pensado para usar en clase/lab. Incluye lo mínimo para arrancar:
- Docker Compose con MongoDB + init básico
- API Go con Gin (solo `GET /items` implementado)
- Archivos y estructura clave ya listos
- Actividades/consignas para completar el CRUD y features avanzados

## Requisitos
- Docker + Docker Compose
- Go 1.22+

## Levantar Mongo
```bash
cp .env.example .env
docker compose up -d
docker compose logs -f mongo
```

## Correr la API (local, fuera de Docker)
```bash
export MONGODB_URI="mongodb://appuser:apppass@localhost:27017/app?authSource=app&retryWrites=true&w=majority"
export MONGODB_DB=app
export PORT=8080
go run ./cmd/api
```

## Endpoints (base)
- GET `/healthz`
- GET `/items` ✅ Implementado
- POST `/items` ❌ TODO
- GET `/items/:id` ❌ TODO
- PUT `/items/:id` ❌ TODO
- DELETE `/items/:id` ❌ TODO

## Consignas
1. **Create**: implementar `POST /items` (validar `name` y `price >= 0`, timestamps).
2. **Get por ID**: implementar `GET /items/:id` (validar ObjectID).
3. **Update**: implementar `PUT /items/:id` (parcial + timestamps).
4. **Delete**: implementar `DELETE /items/:id`.

## Comandos útiles en mongosh
```javascript
use app
show collections
db.items.find().limit(5)
db.items.insertOne({ name: "Demo", price: 9.5, createdAt: new Date(), updatedAt: new Date() })
```
