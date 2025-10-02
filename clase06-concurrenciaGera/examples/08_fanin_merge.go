package main

import (
	"fmt"
	"sync"
	"time"
)

func productor(tag string, n int, delay time.Duration) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for i := 1; i <= n; i++ {
			time.Sleep(delay)
			out <- fmt.Sprintf("%s-%d", tag, i)
		}
	}()
	return out
}

func fanIn[T any](chs ...<-chan T) <-chan T {
	out := make(chan T)
	var wg sync.WaitGroup
	wg.Add(len(chs))
	for _, ch := range chs {
		go func(c <-chan T) {
			defer wg.Done()
			for v := range c {
				out <- v
			}
		}(ch)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// varios canales de entrada en un solo canal de salida y consume todo desde ahí.
// El orden no está garantizado entre A-* y B-* porque llegan con ritmos distintos
func main() {
	a := productor("A", 3, 120*time.Millisecond)
	b := productor("B", 5, 80*time.Millisecond)

	for v := range fanIn(a, b) {
		fmt.Println(v)
	}
}
