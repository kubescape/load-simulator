package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/time/rate"
)

func main() {
	// Create a context that can be used to stop the goroutines.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Go routine to generate file opens
	go runAtRate(ctx, "os.Open", 100, func() {
		file, err := os.Open("/dev/null")
		if err != nil {
			log.Println("Error opening file:", err)
			return
		}
		file.Close()
	})

	// Go routine to generate execs
	go runAtRate(ctx, "cmd.Run", 5, func() {
		cmd := exec.CommandContext(ctx, "true")
		if err := cmd.Run(); err != nil {
			log.Println("Error executing command:", err)
		}
	})

	// Start a simple TCP server on port 8080 for testing
	go func() {
		ln, err := net.Listen("tcp", ":8080")
		if err != nil {
			log.Println("Error starting test server:", err)
			return
		}
		defer ln.Close()
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()

	// Go routine to generate network connections
	go runAtRate(ctx, "net.DialTimeout", 100, func() {
		address := "127.0.0.1:8080"
		conn, err := net.DialTimeout("tcp", address, time.Millisecond)
		if err != nil {
			log.Println("Error connecting to", address, err)
			return
		}
		_ = conn.Close()
	})

	log.Println("Started system activity generator...")

	// Wait for interrupt (Ctrl-C) or SIGTERM
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
	log.Println("Received interrupt. Shutting down...")
	cancel()
	// Give goroutines a moment to exit gracefully
	time.Sleep(1 * time.Second)
}

func runAtRate(ctx context.Context, name string, timesPerSecond int, f func()) {
	limiter := rate.NewLimiter(rate.Limit(timesPerSecond), timesPerSecond)
	log.Println("Running", name, "at rate:", timesPerSecond, "times per second")
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := limiter.Wait(ctx); err != nil {
				return
			}
			f()
		}
	}
}
