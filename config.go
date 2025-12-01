package main

import (
	"encoding/json"
	"os"
	"log"
	"flag"
)

type Config struct {
	Port string `json:"port"`
	LogPath string `json:"logPath"`
	Debug bool `json:"debug"`
	Commands map[string]map[string][]string `json:"commands"`
}

var ConfigValue Config = getConfig()

func getConfig() Config {
	configPath := flag.String("config", "", "Path to the config file")

	flag.Parse()

	if *configPath == "" {
		log.Fatal("No config path provided. Use --config=/path/to/config")
	}

	file, err := os.ReadFile(*configPath)
	if err != nil {
		log.Fatalf("%v", err)
	}

	var configValue Config
	err = json.Unmarshal(file, &configValue)
	if err != nil {
		log.Fatalf("%v", err)
	}

	validateConfig(configValue)

	return configValue
}

func validateConfig(config Config) {
    if config.Port == "" {
        log.Fatal("Config file error. Key 'port' is missing")
    }

    if config.LogPath == "" {
        log.Fatal("Config file error. Key 'logPath' is missing")
    }

    if !JsonFieldExists("debug", config) {
        log.Fatal("Config file error. Key 'debug' is missing")
    }

    if len(config.Commands) == 0 {
        log.Fatal("Configfile error. Key 'commands' is missing or empty")
    }
}