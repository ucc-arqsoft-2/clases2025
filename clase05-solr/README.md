# Clase 05 — Motores de búsqueda SOLR (Completar)

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

## Ver la cache desde tu PC
Cuando completes el punto 4, podrás:
```bash
curl -s http://localhost:8080/__cache/keys | jq .
curl -s "http://localhost:8080/__cache/get?key=items:all" | jq .
```

## Notas
- Memcached está expuesto en el puerto 11211 del host para que puedas probar herramientas externas.
- Mongo se inicializa con `mongo-init/seed.js`.
