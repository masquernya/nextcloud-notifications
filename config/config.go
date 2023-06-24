package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	CloudUrl       string
	DiscordWebhook string
	DiscordPingId  string
}

var config *Config

func Get() *Config {
	if config == nil {
		loadConfig()
	}
	return config
}

func loadConfig() {
	bits, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bits, &config)
	if err != nil {
		panic(err)
	}

}
