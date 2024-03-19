package standalone

import (
	"newrelic/multienv/pkg/model"

	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

type ProcessorFunc = func(any) []model.MeltModel

type ProcWorkerConfig struct {
	Processor  ProcessorFunc
	Model      any
	InChannel  <-chan map[string]any
	OutChannel chan<- model.MeltModel
}

var procWorkerConfigHoldr SharedConfig[ProcWorkerConfig]

func InitProcessor(config ProcWorkerConfig) {
	procWorkerConfigHoldr.SetConfig(config)
	if !procWorkerConfigHoldr.SetIsRunning() {
		log.Println("Starting processor worker...")
		go processorWorker()
	} else {
		log.Println("Processor worker already running, config updated.")
	}
}

func processorWorker() {
	for {
		config := procWorkerConfigHoldr.Config()
		model := config.Model
		data := <-config.InChannel
		err := mapstructure.Decode(data, &model)
		if err == nil {
			for _, val := range config.Processor(model) {
				config.OutChannel <- val
			}
		} else {
			log.Error("Error decoding data = ", err)
			// Sending to processor anyway, just in case it knows a better way to decode
			for _, val := range config.Processor(data) {
				config.OutChannel <- val
			}
		}
	}
}
