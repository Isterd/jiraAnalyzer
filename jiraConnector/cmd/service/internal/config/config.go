package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	handler "jiraAnalyzer/jiraConnector/internal/handler/http"
	"jiraAnalyzer/jiraConnector/internal/jiraclient"
	"jiraAnalyzer/jiraConnector/internal/repository/database"
	"log"
	"os"
)

var (
	ErrOpenConfig  = errors.New("failed to open config file")
	ErrParseConfig = errors.New("error parsing config")
)

type Config struct {
	DB     database.DBConfig          `yaml:"DBSettings"`
	Jira   jiraclient.ProgramSettings `yaml:"ProgramSettings"`
	Server handler.ServerConfig       `yaml:"Server"`
}

func LoadConfig(ConfigPathFlag string) (Config, error) {
	var config Config

	configFile, err := os.Open(ConfigPathFlag)
	if err != nil {
		return Config{}, fmt.Errorf("%w: %w", ErrOpenConfig, err)
	}
	defer configFile.Close()

	yamlDecoder := yaml.NewDecoder(configFile)
	if err := yamlDecoder.Decode(&config); err != nil {
		return config, fmt.Errorf("%w: %w", ErrParseConfig, err)
	}

	log.Printf("Loaded configuration: %+v", config)
	return config, nil
}
