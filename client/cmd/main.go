package main

import (
	"log"
	"os"

	"github.com/z3vxo/kronos/internal/cli"
	"github.com/z3vxo/kronos/internal/config"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "setup" {
		if err := config.Setup(); err != nil {
			log.Fatalf("Setup failed: %v", err)
		}
		log.Println("Config written to client.yaml")
		return
	}

	if err := config.LoadCfg(); err != nil {
		log.Fatalf("%v", err)
	}

	cli, err := cli.NewCli()
	if err != nil {
		log.Fatalf("Failed Loading config: %v", err)
	}

	cli.Run()

}
