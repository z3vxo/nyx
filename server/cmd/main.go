package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/z3vxo/kronos/internal/config"
	"github.com/z3vxo/kronos/internal/database"
	"github.com/z3vxo/kronos/internal/teamserver"
	"gopkg.in/yaml.v3"
)

func SetupKronos() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	basePath := filepath.Join(home, ".nyx")
	configDir := filepath.Join(basePath, "config")
	dbDir := filepath.Join(basePath, "database")

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return err
	}

	configFile := filepath.Join(configDir, "config.yaml")
	dbFile := filepath.Join(dbDir, "nyx_db.sql")
	logFile := filepath.Join(configDir, "nyx.log")

	if err := ensureFile(logFile); err != nil {
		return err
	}
	if err := ensureFile(dbFile); err != nil {
		return err
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		defaultCfg := config.Config{
			TS: config.TeamServer{
				ListenInterface: "127.0.0.1",
				Port:            50050,
				Auth: config.AuthConf{
					Username:          "nyx",
					Password:          "nyxpwd",
					JwtSecret:         "nyxtest123",
					TokenHours:        24,
					TokenRefreshHours: 168,
				},
				Cert: "~/.nyx/certs/cert.crt",
				Key:  "~/.nyx/certs/server.key",
			},
			Server: config.HttpServer{
				GetEndpoint:  "/ms/download",
				PostEndpoint: "/ms/upload",
				GetHeaders:   map[string]string{"Server": "apache"},
				PostHeaders:  map[string]string{"Server": "nginx"},
			},
		}

		data, err := yaml.Marshal(defaultCfg)
		if err != nil {
			return err
		}
		if err := os.WriteFile(configFile, data, 0644); err != nil {
			return err
		}
	} else {
		if err := config.LoadConfig(); err != nil {
			return err
		}
	}
	return nil
}

func ensureFile(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		return file.Close()
	}
	return nil
}

func main() {
	fmt.Println("starting!")
	if err := SetupKronos(); err != nil {
		fmt.Println("failed")
	}
	if err := database.InitDB(); err != nil {
		fmt.Println("Failed Setting up DB")
	}

	ts := teamserver.NewTeamServer()
	if err := ts.Start(); err != nil && err != http.ErrServerClosed {
		fmt.Println("Failed Setting up server")
	}
	return
}
