package config

import (
	_ "embed"
	"gopkg.in/yaml.v3"
	"log"
)

//go:embed config.yaml
var config []byte
var parsedConfig map[string]Config

type Config struct {
	Broker   string `yaml:"broker"`
	Topic    string `yaml:"topic"`
	ClientID string `yaml:"client_id"`
}

func init() {
	err := yaml.Unmarshal(config, &parsedConfig)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func GetConfig(id string) (Config, bool) {
	config, ok := parsedConfig[id]
	return config, ok
}
