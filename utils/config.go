package utils

import (
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

var Config struct {
	mu             sync.Mutex
	Authentication struct {
		API string `yaml:"api"`
	} `yaml:"auth"`

	Server struct {
		Protocol string `yaml:"protocol"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Version  string `yaml:"version"`
		Update   string `yaml:"update"`
		Retries  int    `yaml:"retries"`
	} `yaml:"server"`
}

func ImportConfig(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	log.Printf("found configuration file %v", path)

	decoder := yaml.NewDecoder(file)
	Config.mu.Lock()
	defer Config.mu.Unlock()
	err = decoder.Decode(&Config)
	if err != nil {
		return err
	}

	log.Printf("succsessfully imported configuration from %v", path)
	return nil
}

func GetAPI() string {
	return Config.Authentication.API
}

func GetBaseURL() string {
	return Config.Server.Protocol + "://" + Config.Server.Host + ":" + Config.Server.Port + "/" + Config.Server.Version
}

func GetUpdateTime() string {
	return Config.Server.Update
}

func GetRetries() int {
	return Config.Server.Retries
}
