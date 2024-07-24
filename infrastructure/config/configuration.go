package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type AppConfig struct {
	DatabaseConfig struct {
		Port     string `yaml:"port"`
		DBName   string `yaml:"name"`
		Domain   string `yaml:"domain"`
		Password string `yaml:"password"`
		Username string `yaml:"username"`
		Url      string `yaml:"url"`
	} `yaml:"database"`

	RedisConfig struct {
		Port           string `yaml:"port"`
		Password       string `yaml:"password"`
		DBName         string `yaml:"name"`
		ExpirationTime string `yaml:"expirationTime"`
	} `yaml:"redis"`

	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`

	Host string `yaml:"host"`
}

func LoadConfig() (*AppConfig, error) {
	config := &AppConfig{}

	file, err := os.Open("./infrastructure/config/application.yml")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}
	return config, nil
}