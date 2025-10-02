# Goroutines & Channels — Ejemplos prácticos (Go)

Este paquete contiene **ejemplos autocontenidos** sobre goroutines, canales (unbuffered/buffered), `select`, cancelación con `context`, y patrones de concurrencia comunes (**generador**, **fan‑out/worker pool**, **fan‑in**, **pipeline** y **backpressure**).

> Probado con **Go 1.22+**. Podés ejecutar cada ejemplo con **Make** o directamente con `go run` (ideal en Windows si no tenés Make).

---

## Requisitos

- Go 1.22 o superior (`go version`).
- (Opcional) **Make** para usar los atajos:
  - macOS (Homebrew): `brew install make`
  - Linux (Debian/Ubuntu): `sudo apt-get install make`
  - Windows (Chocolatey): `choco install make` (o usá los comandos sin Make).

---

## ¿Qué es una goroutine?

Una **goroutine** es una ejecución concurrente *muy liviana* gestionada por el runtime de Go. El scheduler multiplexa **miles** de goroutines sobre **unos pocos** hilos del SO (modelo M:N). Cada goroutine tiene una **pila pequeña y dinámica**. Si se bloquea (I/O, canal, `Sleep`…), el hilo puede correr otra goroutine.

Se lanzan con `go <func(...)>` y suelen coordinarse con `sync.WaitGroup`, canales o `context`.

---

## Canales (channels): tuberías tipadas

- `chan T` es un conducto para valores de tipo `T`.
- **Unbuffered**: `make(chan T)` (capacidad 0) → comunicación **sincrónica** (handshake).
- **Buffered**: `make(chan T, N)` (capacidad N) → desacopla emisor/receptor *hasta* N.

Operaciones:
- Enviar: `ch <- v` (bloquea si **lleno** en buffered / si **no hay receptor** en unbuffered).
- Recibir: `x := <-ch` (bloquea si **vacío**).
- Cerrar: `close(ch)` (lo hace el **productor** cuando no enviará más).

Garantías clave:
- Orden **FIFO** por emisor dentro de un canal.
- `send` **happens-before** el `receive` correspondiente (garantía de memoria).

---

## Sincronización (unbuffered) vs buffering (buffered)

- **Unbuffered**:
  - El *sender* se bloquea hasta que exista un *receiver* listo.
  - El *receiver* se bloquea hasta que exista un *sender* listo.
  - Útil para **backpressure fuerte** y handshake exacto.

- **Buffered**:
  - Desacopla temporalmente; el *sender* no bloquea hasta llenar la capacidad.
  - **No elimina** la sincronización por completo: si el buffer está **lleno** → `send` bloquea; si está **vacío** → `receive` bloquea.
  - Útil para **suavizar picos**, colas de trabajo y mejorar throughput.

---

## `select`: multiplexar canales, timeouts y `default`

`select` permite esperar múltiples operaciones de canal y ejecutar la que esté lista (elección aleatoria si varias). Permite implementar **timeouts** (`time.After`) y operaciones **no bloqueantes** (`default`).

---

## Cierre de canales

- Cierra el **productor** con `close(ch)` cuando no habrá más envíos.
- `for range ch` itera hasta que el canal se cierre y drene.
- Lectura con coma-ok: `v, ok := <-ch` → `ok=false` si está cerrado y vacío.
- **Nunca** envíes a un canal cerrado (causa `panic`). No existe `isClosed(ch)`.

---

# Patrones usados y cuándo elegirlos

## 1) Generador (+ cancelación con `context`)
**Qué es:** una función que devuelve `<-chan T` y produce valores en una goroutine interna.  
**Cuándo:** producir flujos (ticks, secuencias, lecturas) sin exponer el canal de envío.  
**Puntos clave:** devolver `<-chan T` (solo lectura), cerrar `out` al terminar y respetar `ctx.Done()`.

*Relacionado:* **Pipeline** (encadenar generadores: producir → transformar → consumir).

## 2) Fan‑out / Worker Pool
**Qué es:** varios workers consumen de un único canal de trabajos y publican resultados.  
**Cuándo:** paralelizar CPU‑bound o I/O‑bound; controlar el grado de concurrencia con **N workers**.  
**Puntos clave:** canal de trabajos (a menudo buffered) + `WaitGroup` + canal de resultados. Cierre ordenado: cerrar **jobs** cuando no hay más, esperar `wg`, luego cerrar **results**.

**Backpressure:** el tamaño del buffer en **jobs**/**results** regula el desacople entre productor/consumidores.

## 3) Fan‑in (merge)
**Qué es:** combinar múltiples canales de entrada en un único canal de salida.  
**Cuándo:** consolidar flujos de varias fuentes o etapas independientes.  
**Puntos clave:** una goroutine por entrada que reenvíe hacia `out`; `WaitGroup`; cerrar `out` al final.

## 4) Pipeline
**Qué es:** cadena de etapas; cada etapa toma un canal de entrada y expone un canal de salida.  
**Cuándo:** procesamientos por etapas (parsear → filtrar → mapear → agrupar).  
**Puntos clave:** cada etapa debe respetar cierre/cancelación (usar `context` o un **done channel**).

## 5) Timeouts, `default` y operaciones no bloqueantes
**Qué es:** usar `select` con `time.After` o rama `default`.  
**Cuándo:** evitar bloqueos indeseados, timeouts de red/IO, *try-send/try-receive*.  
**Puntos clave:** `default` hace que `select` no bloquee (útil para *heartbeats* o sondas).

---

## (Extra) Canales vs. `sync.Mutex`/`WaitGroup`

- **Canales**: para pasar **datos** entre goroutines y coordinar etapas (pipeline).  
- **Mutex/RWMutex**: para proteger **estado compartido** en memoria.  
- Regla práctica: *“No comunique compartiendo memoria; comparta memoria comunicando”*. Si solo protegés una estructura, un **mutex** suele ser más simple y rápido.

## (Extra) `GOMAXPROCS` y scheduling

- `runtime.GOMAXPROCS` controla cuántos **threads de CPU** pueden ejecutar goroutines simultáneamente.  
- Por defecto ≈ # de cores. Subirlo no da más paralelismo si la carga es I/O‑bound.

## (Extra) Evitar fugas (goroutine leaks)

- Cerrá canales cuando no habrá más envíos para que `range` termine.  
- Usá `context` para **cancelación cooperativa** (`<-ctx.Done()`).  
- Evitá `send/receive` sin consumidor/productor garantizado (o agregá timeouts).

## (Extra) Detección de condiciones de carrera

- Ejecutá tests o binarios con `-race`:
  - Linux/macOS: `go test -race ./...` o `go run -race ./examples/xx.go`
  - Windows: `go test -race ./...` (requiere toolchain soportado)

---

# Ejemplos incluidos

- `01_goroutines_waitgroup.go` → lanzar goroutines y esperar con `sync.WaitGroup`.
- `02_unbuffered_sync.go` → sincronía estricta con canal sin buffer.
- `03_buffered_queue.go` → canal con buffer y desacoplamiento.
- `04_directional_context_generator.go` → generador que devuelve `<-chan T` + cancelación.
- `05_select_timeout.go` → `select` con timeouts y ramas múltiples.
- `06_close_and_range.go` → cerrar canales, `for range`, y coma‑ok.
- `07_worker_pool_fanout.go` → pool de workers con jobs/results buffered.
- `08_fanin_merge.go` → fan‑in: merge de múltiples canales en uno.
- `09_deadlock_demo.go` → ejemplo mínimo de deadlock por envío sin receptor.

---

## Cómo correr (Linux/macOS)

Con **Make**:
```bash
make list
make run name=01_goroutines_waitgroup
make run name=07_worker_pool_fanout
```

**Sin Make** (directo con Go):
```bash
ls examples
go run ./examples/01_goroutines_waitgroup.go
go run ./examples/07_worker_pool_fanout.go
```

---

## Cómo correr (Windows)

### PowerShell
Con **Make** (si lo instalaste con Chocolatey):
```powershell
make list
make run name=01_goroutines_waitgroup
```

**Sin Make** (directo con Go):
```powershell
Get-ChildItem .\examples
go run .\examples\01_goroutines_waitgroup.go
go run .\examples\07_worker_pool_fanout.go
```

### CMD clásico
```cmd
dir examples
go run examples\01_goroutines_waitgroup.go
go run examples\07_worker_pool_fanout.go
```

> Si Windows te bloquea la ejecución de scripts, abrí PowerShell **como Administrador** y corré: `Set-ExecutionPolicy RemoteSigned` (podés revertirlo luego).

---

## Buenas prácticas & gotchas

- **Quien envía cierra** el canal; consumidores no cierran.
- Evitá enviar a canales cerrados → `panic`.
- Un canal `nil` bloquea siempre al enviar/recibir.
- Usá `context` para **cancelación cooperativa**.
- Regulá **backpressure** con tamaños de buffer razonables; buffers enormes pueden esconder problemas o gastar memoria.
- Preferí firmas que devuelvan `<-chan T` para encapsular el envío.
