package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/z3vxo/nyx/internal/config"
	"github.com/z3vxo/nyx/internal/database"
	"github.com/z3vxo/nyx/internal/server"
)

func SetupNyx() error {
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

	configFile := filepath.Join(configDir, "config.json")
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
			User:         "nyx",
			Passwd:       "nyxpwd",
			JwtSecret:    "nyxtest123",
			GetEndpoint:  "/ms/download",
			PostEndpoint: "/ms/Upload",
			GetHeaders: []config.Headers{
				{
					Name:  "Server",
					Value: "apache 2.4",
				},
			},
			PostHeaders: []config.Headers{
				{
					Name:  "Server",
					Value: "apache 2.4",
				},
			},
		}

		data, err := json.MarshalIndent(defaultCfg, "", "  ")
		if err != nil {
			return err
		}
		if err := os.WriteFile(configFile, data, 0644); err != nil {
			return err
		}
	} else {
		err := config.LoadConfig()
		if err != nil {
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
	if err := SetupNyx(); err != nil {
		fmt.Println("failed")
	}
	if err := database.InitDB(); err != nil {
		fmt.Println("Failed Setting up DB")
	}

	if err := server.Setup(); err != nil {
		fmt.Println("Failed Setting up server")

	}
	return
}
