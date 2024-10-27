package main

import (
	"fmt"
	"sync"
)

type Task struct {
	num                int
	isPrimeCalculation bool
}

var numOfWorker = 5

var tasks = []Task{
	{num: 2, isPrimeCalculation: true},
	{num: 3, isPrimeCalculation: true},
	{num: 5, isPrimeCalculation: true},
	{num: 7, isPrimeCalculation: true},
	{num: 10, isPrimeCalculation: false},
	{num: 15, isPrimeCalculation: false},
	{num: 20, isPrimeCalculation: false},
}

func isPrime(num int) bool {
	if num < 2 {
		return false
	}

	for i := 2; i*i <= num; i++ {
		if num%i == 0 {
			return false
		}
	}
	return true
}

func fibonacci(n int) int {
	if n < 0 {
		return -1
	}

	if n == 0 {
		return 0
	}

	if n == 1 {
		return 1
	}

	a, b := 0, 1

	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

func worker(taskChan <-chan Task, resultChan chan<- int, workerID int, wg *sync.WaitGroup) {
	defer wg.Done() // Notify that this worker is done
	for task := range taskChan {
		fmt.Printf("Worker %d: Processing task %d\n", workerID, task.num) // Print when worker starts processing

		if task.isPrimeCalculation {
			result := isPrime(task.num)
			if result {
				resultChan <- 1                                            // Prime
				fmt.Printf("Worker %d: %d is prime\n", workerID, task.num) // Print when result is produced
			} else {
				resultChan <- 0                                                // Not prime
				fmt.Printf("Worker %d: %d is not prime\n", workerID, task.num) // Print when result is produced
			}
		} else {
			result := fibonacci(task.num)
			resultChan <- result
			fmt.Printf("Worker %d: Fibonacci of %d is %d\n", workerID, task.num, result) // Print when result is produced
		}
	}
}

func main() {
	fmt.Printf("Hello from worker pool\n")
	var taskChan = make(chan Task)
	var resultChan = make(chan int, len(tasks))
	var wg sync.WaitGroup // WaitGroup to synchronize goroutines

	for i := 0; i < numOfWorker; i++ {
		wg.Add(1)                                 // Add to WaitGroup
		go worker(taskChan, resultChan, i+1, &wg) // Pass worker ID
	}

	// Sending task to the taskChan
	for _, task := range tasks {
		fmt.Printf("Main: Sending task %d\n", task.num) // Print when main sends a task
		taskChan <- task
	}

	close(taskChan) // Close the task channel after sending all tasks

	// Collect the result
	numResult := len(tasks)
	go func() {
		wg.Wait()         // Wait for all workers to finish
		close(resultChan) // Close resultChan after all results are processed
	}()

	for i := 0; i < numResult; i++ {
		result := <-resultChan
		fmt.Println("Main: Received result:", result)
	}
}
