package main

import (
	"fmt"
	"time"
)

func main() {
	a := make(chan string, 1)
	b := make(chan string, 1)

	go func() {
		//le doy 100 ms para que no este listo
		time.Sleep(100 * time.Millisecond) 
		a <- "A listo"
	}()
	go func() {
		time.Sleep(200 * time.Millisecond)
		b <- "B listo"
	}()

	timeout := time.After(150 * time.Millisecond)

	select {
	case x := <-a:
		fmt.Println(x) // probablemente A
	case y := <-b:
		fmt.Println(y)
	case <-timeout:
		fmt.Println("timeout")
	default:
		// rama no bloqueante
	}
}
