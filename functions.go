package main

import (
	"context"
	"crypto/sha256"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/google/uuid"
)

var httpServerAddr = "localhost:8080"

func httpServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello, World!"))
	})
	err := http.ListenAndServe(httpServerAddr, mux)
	if err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}

func dnsFunc() {
	// Generate a unique subdomain to ensure no caching
	// This forces a fresh DNS lookup every time
	uniqueID := uuid.New().String()[:8] // Use first 8 chars of UUID
	hostname := uniqueID + ".nip.io"

	// Create a custom resolver with no caching
	resolver := &net.Resolver{
		PreferGo: true, // Use Go's DNS resolver instead of system resolver
	}

	// Perform DNS lookup with unique hostname
	_, err := resolver.LookupHost(context.Background(), hostname)
	if err != nil {
		log.Println("Error performing DNS lookup for", hostname, ":", err)
	}
}

func execFunc() {
	cmd := exec.Command("true")
	if err := cmd.Run(); err != nil {
		log.Println("Error executing command:", err)
	}
}

func hardlinkFunc() {
	if _, err := os.Stat("/tmp/source_file"); os.IsNotExist(err) {
		_ = os.WriteFile("/tmp/source_file", []byte("test"), 0644)
	}
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return
	}
	dest := "/tmp/hardlink_" + newUUID.String()
	err = os.Link("/tmp/source_file", dest)
	if err != nil {
		log.Println("Error creating hard link:", err)
	}
	_ = os.Remove(dest)
}

func httpFunc() {
	resp, err := http.Get("http://" + httpServerAddr)
	if err != nil {
		log.Println("Error making HTTP request:", err)
		return
	}
	_ = resp.Body.Close()
}

func networkFunc() {
	conn, err := net.DialTimeout("tcp", httpServerAddr, time.Second)
	if err != nil {
		log.Println("Error connecting to", httpServerAddr, err)
		return
	}
	_ = conn.Close()
}

func openFunc() {
	file, err := os.Open("/dev/null")
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	_ = file.Close()
}

func symlinkFunc() {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return
	}
	dest := "/tmp/synlink_" + newUUID.String()
	err = os.Symlink("/dev/null", dest)
	if err != nil {
		log.Println("Error creating hard link:", err)
	}
	_ = os.Remove(dest)
}

func loadSingleCPU(ctx context.Context, loadPercent int) {
	// Ensure loadPercent is within valid range (0-1000 for 0-100%)
	if loadPercent < 0 {
		loadPercent = 0
	} else if loadPercent > 1000 {
		loadPercent = 1000
	}

	// Calculate sleep time to achieve desired CPU load
	// For 25% load, we need to sleep 75% of the time
	sleepTime := time.Duration(1000-loadPercent) * time.Microsecond
	workTime := time.Duration(loadPercent) * time.Microsecond

	log.Printf("Starting CPU load generator at %d%% (sleep: %v, work: %v)",
		loadPercent/10, sleepTime, workTime)

	for {
		select {
		case <-ctx.Done():
			log.Println("CPU load generator stopped")
			return
		default:
			// Do some CPU-intensive work for the specified duration
			start := time.Now()
			for time.Since(start) < workTime {
				// CPU-intensive SHA256 hash calculation
				data := []byte("load-simulator-cpu-work")
				hash := sha256.Sum256(data)
				_ = hash
			}

			// Sleep for the remaining time to achieve desired load percentage
			time.Sleep(sleepTime)
		}
	}
}
