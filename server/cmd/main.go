package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/z3vxo/kronos/internal/config"
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
		log.Fatalf("Failed bootstrap: %v", err)
	}

	ts, err := teamserver.NewTeamServer()
	if err != nil {
		log.Fatalf("Failed Setting up TS: %v", err)
	}
	quit := make(chan os.Signal, 1)
	go func() {
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		fmt.Println("Shutting down...")
		ts.Stop()
	}()

	if err := ts.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed Starting teamserver: %v", err)
	}

}
