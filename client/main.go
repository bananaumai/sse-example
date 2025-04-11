package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	fmt.Println("SSE Client starting...")

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "http://localhost:8080/event", nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		os.Exit(1)
	}

	// Set headers for SSE
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	// Create HTTP client
	client := &http.Client{}

	// Send the request
	fmt.Println("Connecting to SSE server...")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error connecting to server: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
		os.Exit(1)
	}

	fmt.Println("Connected to SSE server")

	// Create a scanner to read the response body line by line
	scanner := bufio.NewScanner(resp.Body)

	// Variables to track the current event
	var eventType string
	var eventData string
	var dataCount int
	var errorCount int

	// Start time for calculating duration
	startTime := time.Now()

	// Process the SSE stream
	for scanner.Scan() {
		line := scanner.Text()

		// Empty line indicates the end of an event
		if line == "" {
			if eventData != "" {
				// Process the complete event
				if eventType == "error" {
					errorCount++
					fmt.Printf("Received ERROR event: %s\n", eventData)
				} else {
					dataCount++
					fmt.Printf("Received DATA event: %s\n", eventData)
				}

				// Reset for the next event
				eventType = ""
				eventData = ""
			}
			continue
		}

		// Parse the line
		if strings.HasPrefix(line, "event:") {
			eventType = strings.TrimSpace(strings.TrimPrefix(line, "event:"))
		} else if strings.HasPrefix(line, "data:") {
			eventData = strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		}
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading SSE stream: %v\n", err)
	}

	// Calculate duration
	duration := time.Since(startTime)

	// Print summary
	fmt.Println("\nSSE stream completed")
	fmt.Printf("Received %d data events and %d error events\n", dataCount, errorCount)
	fmt.Printf("Total duration: %v\n", duration)
}
