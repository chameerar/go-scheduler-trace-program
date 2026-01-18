package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Task represents a unit of work
type Task struct {
	ID int
}

// Worker processes CPU-intensive tasks from a channel
func Worker(id int, tasks <-chan Task, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range tasks {
		computation := 0
		for i := 0; i < 100000000; i++ {
			computation += i % 97
		}
		fmt.Printf("[Worker %d] Task %d done\n", id, task.ID)
	}
}

// DataProducer generates and sends data periodically
func DataProducer(ctx context.Context, dataChan chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(dataChan)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for count := 0; ; count++ {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			select {
			case dataChan <- count:
			case <-ctx.Done():
				return
			}
		}
	}
}

// DataConsumer processes data with CPU-intensive work
func DataConsumer(id int, dataChan <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for data := range dataChan {
		computation := 0
		for i := 0; i < 150000000; i++ {
			computation += (data * i) % 97
		}
		fmt.Printf("[Consumer %d] Processed data: %d\n", id, data)
	}
}

// Monitor periodically prints status
func Monitor(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			fmt.Printf("[Monitor] Running...\n")
		}
	}
}

// ComputeWorker performs CPU-intensive computation
func ComputeWorker(id int, iterations int, wg *sync.WaitGroup) {
	defer wg.Done()
	result := 0
	for i := 0; i < iterations; i++ {
		result += rand.Intn(100)
	}
	fmt.Printf("[ComputeWorker %d] Done\n", id)
}

func main() {
	fmt.Println("Starting multi-threaded Go program for scheduler tracing...")
	fmt.Printf("Start: %s\n", time.Now().Format("15:04:05"))

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var wg sync.WaitGroup

	// Worker pool - CPU-intensive tasks
	tasksChan := make(chan Task, 10)
	wg.Add(6)
	for i := 1; i <= 6; i++ {
		go Worker(i, tasksChan, &wg)
	}
	go func() {
		for i := 1; i <= 20; i++ {
			tasksChan <- Task{ID: i}
			time.Sleep(100 * time.Millisecond)
		}
		close(tasksChan)
	}()

	// Producer-Consumer - Continuous data processing
	dataChan := make(chan int, 5)
	wg.Add(1)
	go DataProducer(ctx, dataChan, &wg)
	wg.Add(4)
	for i := 1; i <= 4; i++ {
		go DataConsumer(i, dataChan, &wg)
	}

	// Monitor - Periodic status checks
	wg.Add(1)
	go Monitor(ctx, &wg)

	// Compute workers - Sustained CPU load
	wg.Add(8)
	for i := 1; i <= 8; i++ {
		go ComputeWorker(i, 5000000, &wg)
	}

	wg.Wait()
	fmt.Printf("End: %s\n", time.Now().Format("15:04:05"))
	fmt.Println("Program complete")
}
