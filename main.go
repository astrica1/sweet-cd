package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fsnotify/fsnotify"
)

func main() {
	yamlPath := "config.yml"
	jsonPath := "config.json"

	err := LoadConfig(yamlPath, jsonPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	SetupEndpoints()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				fmt.Printf("Event: %s\n", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println("Config file updated, reloading...")
					LoadConfig(yamlPath, jsonPath)
					SetupEndpoints()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("Error:", err)
			}
		}
	}()

	if fileExists(yamlPath) {
		err = watcher.Add(yamlPath)
	} else if fileExists(jsonPath) {
		err = watcher.Add(jsonPath)
	}

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting server on :80")
	log.Fatal(http.ListenAndServe(":80", nil))
}
