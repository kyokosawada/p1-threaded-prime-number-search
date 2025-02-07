package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// Config holds the configuration values
type Config struct {
	NumThreads int `json:"num_threads"`
	MaxNumber  int `json:"max_number"`
}

// PrimeResult stores the prime number with thread ID
type PrimeResult struct {
	ThreadID int
	Prime    int
}

// loadConfig reads the config file, checks for validation errors, and returns a Config struct
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

// isPrime checks if a number is prime
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

// printEverything prints all prime results
func printEverything(pr []PrimeResult) {
	for _, result := range pr {
		fmt.Printf("[Thread-%d] Prime: %d\n", result.ThreadID, result.Prime)
	}
}

// Goroutine worker that checks for prime numbers in linear
func worker(threadID, maxNumber, numThreads int, wg *sync.WaitGroup, pr *[]PrimeResult, mu *sync.Mutex) {
	defer wg.Done()
	for i := threadID; i <= maxNumber; i += numThreads {
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

	// Variables for synchronization and results
	var wg sync.WaitGroup
	var mu sync.Mutex
	var pr []PrimeResult

	// Launch goroutines
	for i := 1; i <= config.NumThreads; i++ {
		wg.Add(1)
		go worker(i, config.MaxNumber, config.NumThreads, &wg, &pr, &mu)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Print all prime results
	printEverything(pr)
}