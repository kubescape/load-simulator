package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/time/rate"
)

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Config file not found: %v\n", err)
	} else {
		log.Printf("Using config file: %s\n", viper.ConfigFileUsed())
	}
}

func main() {
	initConfig()

	// Start a simple HTTP server for testing
	go httpServer()

	// Create a context that can be used to stop the goroutines.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Go routine to generate DNS lookups
	if viper.IsSet("dnsRate") {
		go runAtRate(ctx, "dns", viper.GetInt("dnsRate"), dnsFunc)
	}

	// Go routine to generate execs
	if viper.IsSet("execRate") {
		go runAtRate(ctx, "exec", viper.GetInt("execRate"), execFunc)
	}

	// Go routine to generate hard links
	if viper.IsSet("hardlinkRate") {
		go runAtRate(ctx, "hardlink", viper.GetInt("hardlinkRate"), hardlinkFunc)
	}

	// Go routine to generate HTTP requests
	if viper.IsSet("httpRate") {
		go runAtRate(ctx, "http", viper.GetInt("httpRate"), httpFunc)
	}

	if viper.IsSet("cpuLoadMs") {
		numberParallelCPUs := 1
		if viper.IsSet("numberParallelCPUs") {
			numberParallelCPUs = viper.GetInt("numberParallelCPUs")
			maxCPUs := runtime.NumCPU()
			if numberParallelCPUs > maxCPUs {
				log.Printf("Requested numberParallelCPUs (%d) is greater than available CPUs (%d). Using %d.", numberParallelCPUs, maxCPUs, maxCPUs)
				numberParallelCPUs = maxCPUs
			}
		}
		for i := 0; i < numberParallelCPUs; i++ {
			go loadSingleCPU(ctx, viper.GetInt("cpuLoadMs"))
		}
	}

	// Go routine to generate network connections
	if viper.IsSet("networkRate") {
		go runAtRate(ctx, "network", viper.GetInt("networkRate"), networkFunc)
	}

	// Go routine to generate file opens
	if viper.IsSet("openRate") {
		go runAtRate(ctx, "open", viper.GetInt("openRate"), openFunc)
	}

	// Go routine to generate symlinks
	if viper.IsSet("symlinkRate") {
		go runAtRate(ctx, "symlink", viper.GetInt("symlinkRate"), symlinkFunc)
	}

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
	if timesPerSecond <= 0 {
		return
	}
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
