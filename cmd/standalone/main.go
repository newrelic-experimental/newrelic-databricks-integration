package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"newrelic/multienv/pkg/env/standalone"
)

func main() {
	configFile := "./config.yaml"
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}

	log.Print("Loading config file: " + configFile)

	pipeConf, err := standalone.LoadConfig(configFile)
	if err != nil {
		log.Error("Error loading config: ", err)
		os.Exit(1)
	}

	err = standalone.Start(pipeConf)
	if err != nil {
		os.Exit(2)
	}

	select {}
}
