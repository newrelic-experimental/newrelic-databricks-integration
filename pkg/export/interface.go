package export

import (
	"newrelic/multienv/pkg/config"
	"newrelic/multienv/pkg/model"

	log "github.com/sirupsen/logrus"
)

type ExportFunc = func(config.PipelineConfig, []model.MeltModel) error

func dummyExporter(pipeConf config.PipelineConfig, data []model.MeltModel) error {
	log.Warn("Dummy Exporter, do nothing")
	log.Warn("    Data = ", data)
	log.Warn("    Config = ", pipeConf)
	return nil
}

func SelectExporter(exporterType config.ExporterType) ExportFunc {
	switch exporterType {
	case config.NrApi:
		return exportNrApi
	case config.NrEvents:
		return exportNrEvent
	case config.NrMetrics:
		return exportNrMetric
	case config.NrLogs:
		return exportNrLog
	case config.NrTraces:
		return exportNrTrace
	case config.NrInfra:
		return exportNrInfra
	case config.Otel:
		return exportOtel
	case config.Prometheus:
		return exportProm
	default:
		return dummyExporter
	}
}
