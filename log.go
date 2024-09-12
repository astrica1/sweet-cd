package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var logger *log.Logger

func createLogFile() {
	today := time.Now().Format("2006-01-02")
	logDir := "./logs"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.Mkdir(logDir, 0755)
	}

	logFile := filepath.Join(logDir, fmt.Sprintf("%s.log", today))
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	logger = log.New(file, "", log.LstdFlags)
}

func logRequest(service Service, r *http.Request) {
	createLogFile()
	logger.Printf("Request to %s | IP: %s | User-Agent: %s\n", service.Endpoint, r.RemoteAddr, r.UserAgent())
}
