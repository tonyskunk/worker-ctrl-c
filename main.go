package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func worker(id int, wg *sync.WaitGroup, jobs <-chan int) {
	defer wg.Done()
	for job := range jobs {
		fmt.Printf("Worker %d received job %d\n", id, job)
	}
}
func main() {
	jobs := make(chan int)

	var numWorkers int
	fmt.Print("enter numbers of workers: ")
	fmt.Scan(&numWorkers)

	var wg sync.WaitGroup

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, &wg, jobs)
	}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		i := 0
		for {
			select {
			case <-sigs:
				close(jobs)
				return
			default:
				jobs <- i
				i++
				time.Sleep(1 * time.Second)
			}
		}
	}()
	wg.Wait()
	fmt.Println("All workers have finished.")
}
