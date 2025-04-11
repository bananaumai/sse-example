package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	// Define the SSE handler
	http.HandleFunc("/event", sseHandler)

	// Start the server
	port := 8080
	fmt.Printf("SSE server starting on port %d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func sseHandler(w http.ResponseWriter, r *http.Request) {
	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Get request context to detect client disconnection
	ctx := r.Context()

	// Create a flusher to ensure data is sent immediately
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	// Send 10 data messages
	errorCount := 0
	for i := 1; i <= 10; i++ {
		select {
		case <-ctx.Done():
			// Client disconnected
			fmt.Println("Client disconnected")
			return
		default:
			// Send data message
			fmt.Fprintf(w, "data: #%d\n\n", i)
			flusher.Flush()

			// Every other message, send an error
			if i%2 == 0 {
				errorCount++
				fmt.Fprintf(w, "event: error\ndata: error #%d\n\n", errorCount)
				flusher.Flush()
			}

			// Simulate processing time
			time.Sleep(1 * time.Second)
		}
	}

	fmt.Println("SSE stream completed after 10 messages")
}
