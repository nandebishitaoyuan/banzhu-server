package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Addr string `yaml:"addr"`
	} `yaml:"server"`

	Log struct {
		Level string `yaml:"level"`
	} `yaml:"log"`

	Database struct {
		DSN          string `yaml:"dsn"`
		MaxOpenConns int    `yaml:"max_open_conns"`
		MaxIdleConns int    `yaml:"max_idle_conns"`
	} `yaml:"database"`

	JWT struct {
		AccessSecret       string `yaml:"access_secret"`
		AccessExpireHours  int    `yaml:"access_expire_hours"`
		RefreshSecret      string `yaml:"refresh_secret"`
		RefreshExpireHours int    `yaml:"refresh_expire_hours"`
	} `yaml:"jwt"`

	Path struct {
		Book    string `yaml:"book"`
		Chapter string `yaml:"chapter"`
		Images  string `yaml:"images"`
	} `yaml:"path"`
}

func Load(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}
	if cfg.Server.Addr == "" {
		cfg.Server.Addr = ":8080"
	}
	return &cfg, nil
}
