package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// Config values
type Config struct {
	NumThreads int `json:"num_threads"`
	MaxNumber  int `json:"max_number"`
}

// Store prime number with thread ID
type PrimeResult struct {
	ThreadID int
	Prime    int
}

// Reads the config file, checks for validation errors , returns a Config struct
func loadConfig(filename string) (Config, error) {
	var config Config
	file, err := os.ReadFile(filename)
	if err != nil {
		return config, fmt.Errorf("failed to read config file: %w", err)
	}
	err = json.Unmarshal(file, &config)
	if err != nil {
		return config, fmt.Errorf("failed to parse config file: %w", err)
	}
	// Input validation
	if config.NumThreads <= 0 {
		return config, fmt.Errorf("invalid config: num_threads must be greater than 0")
	}
	if config.MaxNumber <= 0 {
		return config, fmt.Errorf("invalid config: max_number must be greater than 0")
	}
	return config, nil
}

// Checks if a number is prime
func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

// Immediately prints the result with thread ID and timestamp
func printEverything(pr []PrimeResult) {
	for _, result := range pr {
		fmt.Printf("[Thread-%v] Prime: %v\n", result.ThreadID, result.Prime)
	}
}

// Goroutine worker that searches for prime numbers
func worker(threadID, start, end int, wg *sync.WaitGroup, pr *[]PrimeResult, mu *sync.Mutex) {
	defer wg.Done()
	for i := start; i <= end; i++ {
		if isPrime(i) {
			mu.Lock()
			*pr = append(*pr, PrimeResult{ThreadID: threadID, Prime: i})
			mu.Unlock()
		}
	}
}

func main() {
	// Load configuration
	config, err := loadConfig("config.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}


	// Calculate the range for each thread
	rangeSize := config.MaxNumber / config.NumThreads
	var wg sync.WaitGroup
	var mu sync.Mutex
	var pr []PrimeResult

	
	// Launch goroutines
	for i := 0; i < config.NumThreads; i++ {
		wg.Add(1)
		start := i*rangeSize + 1
		end := (i + 1) * rangeSize
		if i == config.NumThreads-1 {
			end = config.MaxNumber // Ensure the last thread covers the remaining range
		}
		go worker(i+1, start, end, &wg, &pr, &mu)
	}

	wg.Wait()
	printEverything(pr)
}