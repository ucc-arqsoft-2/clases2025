package main

import "fmt"

func main() {
	ch := make(chan int, 2) // capacidad 2

	ch <- 1 // no bloquea
	ch <- 2 // no bloquea

	// fmt.Println(<-ch) // 1

	// ch <- 3 // descomentar: bloquearía si nadie lee

	fmt.Println(<-ch) // 1
	fmt.Println(<-ch) // 2
}
