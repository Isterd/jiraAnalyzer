package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	handler "jiraAnalyzer/jiraConnector/internal/handler/http"
	"jiraAnalyzer/jiraConnector/internal/repository/database"
	"jiraAnalyzer/jiraConnector/internal/repository/jira"
	"log"
	"os"
)

var (
	ErrOpenConfig  = errors.New("failed to open config file")
	ErrParseConfig = errors.New("error parsing config")
)

type Config struct {
	DB            database.DBConfig           `yaml:"DBSettings"`
	ClientConfig  jira.ClientConfig           `yaml:"JiraClient"`
	JiraConnector handler.JiraConnectorConfig `yaml:"JiraConnector"`
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
