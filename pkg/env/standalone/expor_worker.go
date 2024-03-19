package standalone

import (
	"newrelic/multienv/pkg/config"
	"newrelic/multienv/pkg/export"
	"newrelic/multienv/pkg/model"
	"time"

	log "github.com/sirupsen/logrus"
)

type ExpWorkerConfig struct {
	InChannel   <-chan model.MeltModel
	BatchSize   int
	HarvestTime int
	Exporter    export.ExportFunc
}

var expWorkerConfig SharedConfig[ExpWorkerConfig]
var pipelineConfig SharedConfig[config.PipelineConfig]

func InitExporter(config ExpWorkerConfig, pipeConf config.PipelineConfig) {
	pipelineConfig.SetConfig(pipeConf)
	expWorkerConfig.SetConfig(config)
	if !expWorkerConfig.SetIsRunning() {
		log.Println("Starting exporter worker...")
		go exporterWorker()
	} else {
		log.Println("Exporter worker already running, config updated.")
	}
}

func exporterWorker() {
	buffer := MakeReservoirBuffer[model.MeltModel](500)
	pre := time.Now().Unix()

	for {
		config := expWorkerConfig.Config()
		harvestTime := time.Duration(config.HarvestTime) * time.Second

		data := <-config.InChannel
		switch data.Type {
		case model.Metric:
			metric, _ := data.Metric()
			log.Println("Exporter received a Metric", metric.Name)
		case model.Event:
			event, _ := data.Event()
			log.Println("Exporter received an Event", event.Type)
		case model.Log:
			dlog, _ := data.Log()
			log.Println("Exporter received a Log", dlog.Message, dlog.Type)
		case model.Trace:
			//TODO
			log.Warn("TODO: Exporter received a Trace")
		}

		buffer.Put(data)

		now := time.Now().Unix()
		bufSize := buffer.Size()

		if now-pre >= int64(harvestTime.Seconds()) || bufSize >= config.BatchSize {
			buf := *buffer.Clear()

			log.Println("Harvest cycle, buffer size = ", bufSize)

			err := config.Exporter(pipelineConfig.Config(), buf[0:bufSize])

			if err != nil {
				log.Error("Exporter failed = ", err)
				//TODO: handle error condition, refill buffer? Discard data? Retry?
			}

			pre = time.Now().Unix()
		}
	}
}
