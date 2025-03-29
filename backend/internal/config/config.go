package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"jiraAnalyzer/backend/internal/repository/database"
	"time"
)

type Config struct {
	DBSettings database.DBSettings `yaml:"DBSettings"`
	Backend    Backend             `yaml:"Backend"`
	Logging    Logging             `yaml:"Logging"`
}

type Backend struct {
	BaseUrl          string        `yaml:"baseUrl"`
	Host             string        `yaml:"host"`
	Port             string        `yaml:"port"`
	AnalyticsTimeout time.Duration `yaml:"analyticsTimeout"`
	ResourceTimeout  time.Duration `yaml:"resourceTimeout"`
}

type Logging struct {
	LogFile      string `yaml:"logFile"`
	ErrorLogFile string `yaml:"errorLogFile"`
}

func LoadConfig(filePath string) (*Config, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
