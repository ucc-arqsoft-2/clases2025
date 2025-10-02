package main

import (
	"fmt"
	"sync"
	"time"
)

func trabajo(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("[G%d] arrancó\n", id)
	time.Sleep(200 * time.Millisecond)
	fmt.Printf("[G%d] terminó\n", id)
}

func main() {
	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go trabajo(i, &wg)
	}
	fmt.Println("[main] esperando…")
	wg.Wait()
	fmt.Println("[main] listo")
}
