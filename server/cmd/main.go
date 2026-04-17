package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

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

	basePath := filepath.Join(home, ".kronos")
	configDir := filepath.Join(basePath, "config")
	dbDir := filepath.Join(basePath, "database")

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return err
	}

	configFile := filepath.Join(configDir, "config.yaml")
	dbFile := filepath.Join(dbDir, "kronos_db.sql")
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
					Username:          "kronos",
					Password:          "kronospwd",
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
				NotFoundFile: "~/.kronos/404.html",
			},
		}

		data, err := yaml.Marshal(defaultCfg)
		if err != nil {
			return err
		}
		if err := os.WriteFile(configFile, data, 0644); err != nil {
			return err
		}
	}

	return config.LoadConfig()
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
	go func() {
		if err := ts.Start(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Failed setting up server:", err)
		}
	}()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			ts.SSE.Broadcast("test")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutting down...")
	ts.Stop()
}
