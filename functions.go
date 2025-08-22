package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"time"
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

func openFunc() {
	file, err := os.Open("/dev/null")
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	_ = file.Close()
}

func execFunc() {
	cmd := exec.Command("true")
	if err := cmd.Run(); err != nil {
		log.Println("Error executing command:", err)
	}
}

func networkFunc() {
	conn, err := net.DialTimeout("tcp", httpServerAddr, time.Second)
	if err != nil {
		log.Println("Error connecting to", httpServerAddr, err)
		return
	}
	_ = conn.Close()
}

func httpFunc() {
	resp, err := http.Get("http://" + httpServerAddr)
	if err != nil {
		log.Println("Error making HTTP request:", err)
		return
	}
	_ = resp.Body.Close()
}
