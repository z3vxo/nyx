package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Headers struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type Config struct {
	User         string    `json:"nyx_username"`
	Passwd       string    `json:"nyx_password"`
	JwtSecret    string    `json:"nyx_jwtsecret"`
	GetEndpoint  string    `json:"nyx_getendpoint"`
	PostEndpoint string    `json:"nyx_postendpoint"`
	GetHeaders   []Headers `json:"nyx_getheaders"`
	PostHeaders  []Headers `json:"nyx_posttheaders"`
}

var Cfg *Config

func LoadConfig() error {

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	path := filepath.Join(home, ".nyx", "config", "config.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	Cfg = &Config{}
	return json.Unmarshal(data, Cfg)
}
