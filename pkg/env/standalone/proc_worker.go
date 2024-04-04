package standalone

import (
	"newrelic/multienv/pkg/model"

	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

type ProcessorFunc = func(any) []model.MeltModel

type ProcWorkerConfig struct {
	Processor         ProcessorFunc
	Model             any
	InChannel         <-chan map[string]any
	MetricsOutChannel chan<- model.MeltModel
	EventsOutChannel  chan<- model.MeltModel
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
		configModel := config.Model
		data := <-config.InChannel
		err := mapstructure.Decode(data, &configModel)
		if err == nil {
			for _, val := range config.Processor(configModel) {

				switch val.Type {
				case model.Metric:
					config.MetricsOutChannel <- val
				case model.Event:
					config.EventsOutChannel <- val

				default:
					log.Warn("Model type unknown")
				}
			}
		} else {
			log.Error("Error decoding data = ", err)
		}
	}
}
