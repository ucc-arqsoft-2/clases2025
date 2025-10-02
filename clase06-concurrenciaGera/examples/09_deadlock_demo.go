package main

func main() {
	ch := make(chan int) // unbuffered

	// ch <- 1 // Descomentar para ver deadlock: no hay receptor

	// Arreglo: mover envío a una goroutine o usar buffer
	go func(){ ch <- 1 }()
	<-ch
}
