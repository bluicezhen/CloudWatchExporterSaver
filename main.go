package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// Add global variables and mutex
var (
	metrics     string
	metricsLock sync.RWMutex
)

func getMetrics(cloudwatchExporterURL string) {
	logger := log.New(os.Stdout, "cloudwatch-exporter-saver: ", log.LstdFlags|log.Lshortfile)

	get, err := http.Get(cloudwatchExporterURL)
	if err != nil {
		logger.Printf("Failed to get metrics from cloudwatch exporter: %v", err)
		return
	}
	defer get.Body.Close()

	body, err := io.ReadAll(get.Body)
	if err != nil {
		logger.Printf("Failed to read response body: %v", err)
		return
	}

	logger.Printf("Get metrics from cloudwatch exporter success")

	metricsLock.Lock()
	metrics = string(body)
	metricsLock.Unlock()
}

// Add HTTP handler
func metricsHandler(w http.ResponseWriter, r *http.Request) {
	metricsLock.RLock()
	defer metricsLock.RUnlock()
	w.Write([]byte(metrics))
}

// Add health check handler
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	logger := log.New(os.Stdout, "cloudwatch-exporter-saver: ", log.LstdFlags|log.Lshortfile)

	cloudwatchExporterURL := os.Getenv("CLOUDWATCH_EXPORTER_URL")
	if cloudwatchExporterURL == "" {
		logger.Fatalf("CLOUDWATCH_EXPORTER_URL is not set")
	}

	// Initialize metrics for the first time
	getMetrics(cloudwatchExporterURL)

	// Task 1: Refresh metrics every 5 minutes
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		for range ticker.C {
			getMetrics(cloudwatchExporterURL)
		}
	}()

	// Task 2: Start HTTP server
	http.HandleFunc("/metrics", metricsHandler)
	http.HandleFunc("/health", healthHandler) // Add health check endpoint
	logger.Fatal(http.ListenAndServe(":9107", nil))
}
