# Goroutines & Channels — Ejemplos prácticos (Go)

Este paquete contiene **ejemplos autocontenidos** sobre goroutines, canales (unbuffered/buffered), `select`, cancelación con `context`, y patrones de concurrencia comunes (generador, fan-out/worker pool, fan-in, pipeline y backpressure).

> Probado con **Go 1.22+**. Ejecutá cada ejemplo con:  
> ```bash
> make list
> make run name=<archivo-sin-.go>
> # ej: make run name=01_goroutines_waitgroup
> ```

---

## ¿Qué es una goroutine?

Una **goroutine** es una ejecución concurrente *muy liviana* gestionada por el runtime de Go. El scheduler multiplexa **miles** de goroutines sobre **unos pocos** hilos del SO (modelo M:N). Cada goroutine tiene una **pila pequeña y dinámica**. Si se bloquea (I/O, canal, sleep…), el hilo puede correr otra goroutine.

Se lanzan con `go <func(...)>` y suelen coordinarse con `sync.WaitGroup`, canales o `context`.

---

## Canales (channels): tuberías tipadas

- `chan T` es un conducto para valores de tipo `T`.
- **Unbuffered**: `make(chan T)` capacidad 0 → comunicación **sincrónica** (handshake).
- **Buffered**: `make(chan T, N)` capacidad N → desacopla emisor/receptor *hasta* N.

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
  - **No elimina** la sincronización por completo: si el buffer está **lleno** → el `send` bloquea; si está **vacío** → el `receive` bloquea.
  - Útil para **suavizar picos**, colas de trabajo y mejorar throughput.

---

## `select`: multiplexar canales, timeouts y `default`

`select` permite esperar múltiples operaciones de canal y ejecutar la que esté lista (elección aleatoria si varias). Permite implementar timeouts (`time.After`) y operaciones no bloqueantes (`default`).

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
**Puntos clave:** devolver `<-chan T` (solo lectura para el cliente), cerrar `out` al terminar. Usar `context` para cancelar.

Relacionados: *pipeline stages* que transforman valores encadenando generadores.

## 2) Fan‑out / Worker Pool
**Qué es:** varios workers consumen de un único canal de trabajos y publican resultados.  
**Cuándo:** paralelizar CPU‑bound o I/O‑bound; controlar el grado de concurrencia fijando **N workers**.  
**Puntos clave:** canal de trabajos (a menudo buffered) + `WaitGroup` + canal de resultados. Cierre ordenado: cerrar **jobs** cuando no hay más, esperar `wg`, luego cerrar **results**.

**Backpressure:** el tamaño del buffer en **jobs** y **results** regula cuánto se desacopla productor/consumidores.

## 3) Fan‑in (merge)
**Qué es:** combinar múltiples canales de entrada en un único canal de salida.  
**Cuándo:** consolidar flujos provenientes de varias fuentes o etapas independientes.  
**Puntos clave:** lanzar una goroutine por entrada que reenvíe hacia `out`; usar `WaitGroup` y cerrar `out` cuando todas terminen.

## 4) Pipeline
**Qué es:** cadena de etapas, cada una toma un canal de entrada y expone un canal de salida.  
**Cuándo:** procesamientos en etapas (parsear → filtrar → mapear → agrupar).  
**Puntos clave:** cada etapa respeta cierre/cancelación; `context` o un **done channel** para cortar toda la línea.

## 5) Timeouts, `default` y operaciones no bloqueantes
**Qué es:** usar `select` con `time.After` o rama `default`.  
**Cuándo:** evitar bloqueos indeseados, timeouts de red/IO, o *try-send/try-receive*.  
**Puntos clave:** `default` hace que `select` no bloquee (útil para sondas o *heartbeats*).

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

## Cómo correr

```bash
make list
make run name=01_goroutines_waitgroup
make run name=07_worker_pool_fanout
# etc.
```

---

## Buenas prácticas & gotchas

- **Quien envía cierra** el canal; consumidores no cierran.
- Evitá enviar a canales cerrados → `panic`.
- Un canal `nil` bloquea siempre al enviar/recibir.
- Usá `context` para **cancelación cooperativa**.
- Regulá **backpressure** con tamaños de buffer razonables; buffers enormes pueden esconder problemas o gastar memoria.
- Preferí firmas que devuelvan `<-chan T` para encapsular el envío.
