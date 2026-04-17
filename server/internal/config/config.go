package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type AuthConf struct {
	Username          string `yaml:"user"`
	Password          string `yaml:"pass"`
	JwtSecret         string `yaml:"jwtSecret"`
	TokenHours        int    `yaml:"token_hours"`
	TokenRefreshHours int    `yaml:"token_refresh_hours"`
}

type TeamServer struct {
	ListenInterface string   `yaml:"interface"`
	Port            int      `yaml:"port"`
	Auth            AuthConf `yaml:"auth"`
	Cert            string   `yaml:"cert"`
	Key             string   `yaml:"key"`
}

type HttpServer struct {
	GetEndpoint  string            `yaml:"getEndpoint"`
	PostEndpoint string            `yaml:"postEndpoint"`
	GetHeaders   map[string]string `yaml:"GetResponseHeaders"`
	PostHeaders  map[string]string `yaml:"PostResponseHeaders"`
}

type Config struct {
	TS     TeamServer `yaml:"teamserver"`
	Server HttpServer `yaml:"server"`
}

var Cfg *Config

func LoadConfig() error {

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	path := filepath.Join(home, ".nyx", "config", "config.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	Cfg = &Config{}
	return yaml.Unmarshal(data, Cfg)
}
