package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Service struct {
	Name     string   `json:"name" yaml:"name"`
	Endpoint string   `json:"endpoint" yaml:"endpoint"`
	APIToken string   `json:"-" yaml:"-"`
	Commands []string `json:"commands" yaml:"commands"`
}

type Config struct {
	AppName  string    `json:"app_name" yaml:"app_name"`
	Services []Service `json:"services" yaml:"services"`
}

var config *Config

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func loadYAMLConfig(path string) (*Config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func loadJSONConfig(path string) (*Config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func LoadConfig(yamlPath, jsonPath string) error {
	var err error
	if fileExists(yamlPath) {
		config, err = loadYAMLConfig(yamlPath)
	} else if fileExists(jsonPath) {
		config, err = loadJSONConfig(jsonPath)
	}

	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	return nil
}
