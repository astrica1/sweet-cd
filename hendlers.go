package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

func HandleConfigList(w http.ResponseWriter, r *http.Request) {
	response := struct {
		AppName  string    `json:"app_name"`
		Services []Service `json:"services"`
	}{
		AppName: config.AppName,
		Services: func() []Service {
			sanitizedServices := make([]Service, len(config.Services))
			for i, service := range config.Services {
				sanitizedServices[i] = Service{
					Name:     service.Name,
					Endpoint: service.Endpoint,
					Commands: service.Commands,
				}
			}
			return sanitizedServices
		}(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func runCommands(service Service) {
	for _, cmdStr := range service.Commands {
		cmd := exec.Command("bash", "-c", cmdStr)
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Error executing command %s: %v", cmdStr, err)
		}
		fmt.Printf("Output: %s\n", string(output))
	}
}

func handleRequest(service Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != service.APIToken {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		logRequest(service, r)
		fmt.Fprintf(w, "Executing commands for service: %s\n", service.Name)
		runCommands(service)
	}
}

func SetupEndpoints() {
	http.DefaultServeMux = http.NewServeMux()

	http.HandleFunc("/config", HandleConfigList)

	for _, service := range config.Services {
		http.HandleFunc("/"+service.Endpoint, handleRequest(service))
		fmt.Printf("Endpoint for %s registered at /%s\n", service.Name, service.Endpoint)
	}
}
