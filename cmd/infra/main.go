package main

import (
	"newrelic/multienv/integration"
	"newrelic/multienv/pkg/connect"
	"newrelic/multienv/pkg/env/infra"
	"newrelic/multienv/pkg/export"
	"os"
	"sync"

	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

func main() {
	pipeConf, err := infra.LoadConfig()
	if err != nil {
		log.Error("Error loading config: ", err)
		os.Exit(1)
	}

	recvConfig, err := integration.InitRecv(&pipeConf)
	if err != nil {
		log.Error("Error initializing receiver: ", err)
		os.Exit(2)
	}

	procConfig, err := integration.InitProc(&pipeConf)
	if err != nil {
		log.Error("Error initializing processor: ", err)
		os.Exit(3)
	}

	wg := &sync.WaitGroup{}

	for _, connector := range recvConfig.Connectors {
		wg.Add(1)

		go func(connector connect.Connector) {
			defer wg.Done()

			res, reqErr := connector.Request()
			if reqErr.Err != nil {
				log.Error("Error connecting: ", reqErr.Err)
				return
			}

			var deserBuffer map[string]any
			errDes := recvConfig.Deser(res, &deserBuffer)
			if errDes != nil {
				log.Error("Error deserializing: ", errDes)
				return
			}

			model := procConfig.Model
			errMap := mapstructure.Decode(deserBuffer, &model)
			if errMap != nil {
				log.Error("Error mapping: ", errMap)
				return
			}

			meltData := integration.Proc(model)

			exporter := export.SelectExporter(pipeConf.Exporter)

			errExp := exporter(pipeConf, meltData)
			if errExp != nil {
				log.Error("Error exporting: ", errExp)
				return
			}
		}(connector)
	}
	wg.Wait()
}
