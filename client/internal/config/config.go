package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type HttpConf struct {
	User   string `yaml:"user"`
	Passwd string `yaml:"pass"`
	Host   string `yaml:"host"`
}

type Config struct {
	Http HttpConf `yaml:"http"`
}

var Cfg *Config

func Setup() error {
	r := bufio.NewReader(os.Stdin)
	prompt := func(label string) string {
		fmt.Print(label)
		line, _ := r.ReadString('\n')
		return strings.TrimSpace(line)
	}

	host := prompt("Teamserver host (e.g. http://127.0.0.1:50050): ")
	user := prompt("Username: ")
	passwd := prompt("Password: ")

	cfg := Config{Http: HttpConf{Host: host, User: user, Passwd: passwd}}
	data, err := yaml.Marshal(&cfg)
	if err != nil {
		return err
	}
	return os.WriteFile("./client.yaml", data, 0600)
}

func LoadCfg() error {
	data, err := os.ReadFile("./client.yaml")
	if os.IsNotExist(err) {
		return fmt.Errorf("missing config: run ./client setup")
	}
	if err != nil {
		return err
	}
	Cfg = &Config{}
	return yaml.Unmarshal(data, Cfg)
}
