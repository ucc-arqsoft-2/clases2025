package main

import "fmt"

func main() {
	ch := make(chan int, 3)
	for i := 1; i <= 3; i++ {
		ch <- i
	}
	close(ch) // no habrá más envíos

	for v := range ch {
		fmt.Println(v)
	}

	// Leer de un canal cerrado: si aún hay datos en buffer, te da esos datos; cuando se agotan, te da el cero del tipo y ok=false.
	_, ok := <-ch
	fmt.Println("ok:", ok) // false
}
