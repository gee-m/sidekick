package config

import (
	"appgents/internal/core/email"
	"encoding/json"
	"os"
)

type Config struct {
	Environment string `json:"environment"`
	Server      struct {
		Port         string `json:"port"`
		ReadTimeout  int    `json:"read_timeout"`
		WriteTimeout int    `json:"write_timeout"`
		BaseURL      string `json:"base_url"`
	} `json:"server"`

	Database struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		DBName   string `json:"dbname"`
		SSLMode  string `json:"sslmode"`
	} `json:"database"`

	Email email.Config `json:"email"`
}

func Load(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
