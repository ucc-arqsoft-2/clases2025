package main

import "fmt"

func main() {
	ch := make(chan string) // unbuffered

	go func() {
		fmt.Println("prod: preparando dato…")
		ch <- "hola" // bloquea hasta que alguien lea
		// ch <- "hola2" // bloquea hasta que alguien lea
		fmt.Println("prod: ya se envió y sigo")
	}()

	msg := <-ch // bloquea hasta que haya un send
	fmt.Println("cons: recibido =", msg)

	// msg2 := <-ch // bloquea hasta que haya un send
	// fmt.Println("cons: recibido =", msg2)
}
