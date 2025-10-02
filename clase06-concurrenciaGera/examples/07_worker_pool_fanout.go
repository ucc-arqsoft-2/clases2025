package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Job struct{ N int }
type Result struct {
	N, Sq int
}

func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		time.Sleep(time.Duration(50+rand.Intn(50)) * time.Millisecond)
		results <- Result{N: j.N, Sq: j.N * j.N}
		fmt.Printf("[W%d] procesó %d\n", id, j.N)
	}
}

// El orden NO está garantizado

// jobs y results tienen buffer de 10 → el productor puede encolar todo sin bloquear, 
// y los workers pueden publicar resultados sin frenar al main.

func main() {
	rand.Seed(time.Now().UnixNano())

	const (
		numWorkers = 3
		numJobs    = 10
	)

	jobs := make(chan Job, numJobs)
	results := make(chan Result, numJobs)

	var wg sync.WaitGroup
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	for i := 1; i <= numJobs; i++ {
		jobs <- Job{N: i}
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		fmt.Printf("resultado: %d^2 = %d\n", r.N, r.Sq)
	}
}
