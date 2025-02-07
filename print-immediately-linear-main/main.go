package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// Config values
type Config struct {
	NumThreads int `json:"num_threads"`
	MaxNumber  int `json:"max_number"`
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
func printImmediately(threadID, number int, startTime time.Time) {
	elapsed := time.Since(startTime).Milliseconds()
	fmt.Printf("[Thread %v] %v ms Found prime number: %v\n", threadID, elapsed, number)
}

// Goroutine worker that searches for prime numbers
func worker(threadID, maxNumber, numThreads int, wg *sync.WaitGroup, startTime time.Time) {
	defer wg.Done()
	for i := threadID; i <= maxNumber; i += numThreads {
		if isPrime(i) {
			printImmediately(threadID, i, startTime)
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
	startTime := time.Now()

	// Calculate the range for each thread
	var wg sync.WaitGroup
	
	// Launch goroutines
	for i := 0; i < config.NumThreads; i++ {
		wg.Add(1)

		go worker(i+1, config.MaxNumber, config.NumThreads, &wg, startTime)
	}

	wg.Wait()

}