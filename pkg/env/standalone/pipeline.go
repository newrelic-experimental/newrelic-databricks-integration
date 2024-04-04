package standalone

import (
	"newrelic/multienv/integration"
	"newrelic/multienv/pkg/config"
	"newrelic/multienv/pkg/export"
	"newrelic/multienv/pkg/model"

	log "github.com/sirupsen/logrus"
)

// Init data pipeline.
func InitPipeline(pipeConf config.PipelineConfig, recvConfig config.RecvConfig, procConfig config.ProcConfig, proc ProcessorFunc) {
	bufferSize, _ := pipeConf.GetInt("buffer")

	if bufferSize < 100 {
		bufferSize = 100
	}

	if bufferSize > 1000 {
		bufferSize = 1000
	}

	batchSize, _ := pipeConf.GetInt("batch_size")

	if batchSize < 10 {
		batchSize = 10
	}

	if batchSize > 100 {
		batchSize = 100
	}

	harvestTime, _ := pipeConf.GetInt("harvest_time")

	if harvestTime < 60 {
		harvestTime = 60
	}

	recvToProcCh := make(chan map[string]any, bufferSize)
	metricsProcToExpCh := make(chan model.MeltModel, bufferSize)
	eventsProcToExpCh := make(chan model.MeltModel, bufferSize)

	metricsExporter := MakeExporterWorker(ExpWorkerConfig{
		InChannel:   metricsProcToExpCh,
		HarvestTime: harvestTime,
		BatchSize:   batchSize,
		Exporter:    export.SelectExporter(config.NrMetrics),
	}, pipeConf)

	metricsExporter.InitExporter()

	eventsExporter := MakeExporterWorker(ExpWorkerConfig{
		InChannel:   eventsProcToExpCh,
		HarvestTime: harvestTime,
		BatchSize:   batchSize,
		Exporter:    export.SelectExporter(config.NrEvents),
	}, pipeConf)

	eventsExporter.InitExporter()

	InitProcessor(ProcWorkerConfig{
		Processor:         proc,
		Model:             procConfig.Model,
		InChannel:         recvToProcCh,
		MetricsOutChannel: metricsProcToExpCh,
		EventsOutChannel:  eventsProcToExpCh,
	})
	InitReceiver(RecvWorkerConfig{
		IntervalSec:  pipeConf.Interval,
		Connectors:   recvConfig.Connectors,
		Deserializer: recvConfig.Deser,
		OutChannel:   recvToProcCh,
	})
}

// Start Integration
func Start(pipeConf config.PipelineConfig) error {
	recvConfig, err := integration.InitRecv(&pipeConf)
	if err != nil {
		log.Error("Error initializing receiver: ", err)
		return err
	}
	procConfig, err := integration.InitProc(&pipeConf)
	if err != nil {
		log.Error("Error initializing processor: ", err)
		return err
	}
	InitPipeline(pipeConf, recvConfig, procConfig, integration.Proc)
	return nil
}
