package main

import (
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
	_, err := net.LookupHost("google.com")
	if err != nil {
		log.Println("Error performing DNS lookup:", err)
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
