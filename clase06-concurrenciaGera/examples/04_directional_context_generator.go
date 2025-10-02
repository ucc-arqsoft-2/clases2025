package main

import (
	"context"
	"fmt"
	"time"
)

// generador crea y devuelve un canal **solo lectura** (<-chan int).
// Internamente lanza una goroutine que:
//  - Emite enteros consecutivos comenzando en `start`.
//  - Respeta cancelación/timeout a través de ctx.Done().
//  - Cierra el canal cuando termina (defer close(out)).
func generador(ctx context.Context, start int) <-chan int {
	out := make(chan int) // canal sin buffer (unbuffered) para handshakes send/receive

	go func() {
		defer close(out) // garantizamos cierre del canal al salir de la goroutine
		for i := start; ; i++ { // ciclo infinito produciendo i, i+1, i+2, ...
			select {
			case <-ctx.Done():
				// Si el contexto se cancela o vence el timeout,
				// dejamos de producir y salimos (lo que dispara el close(out) del defer).
				return
			case out <- i:
				// Si hay un receiver listo en el otro extremo, se envía i.
				// Luego simulamos trabajo/ritmo de emisión con un sleep.
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	// Devolvemos el canal de salida como **solo lectura** para el consumidor.
	return out
}

func main() {
	// Creamos un contexto con timeout de ~450 ms.
	// Pasado ese tiempo, ctx.Done() se cierra y la goroutine productora se cancela.
	ctx, cancel := context.WithTimeout(context.Background(), 450*time.Millisecond)
	defer cancel() // buena práctica: liberar recursos del contexto

	// Consumimos los valores del generador usando "range" sobre el canal.
	// El for-range termina automáticamente cuando el canal se **cierra**.
	for v := range generador(ctx, 10) {
		fmt.Println("recibí:", v)
	}

	// Al terminar el range, sabemos que el canal se cerró (porque el generador retornó).
	fmt.Println("fin")
}
