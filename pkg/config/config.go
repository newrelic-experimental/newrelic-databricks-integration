package config

import (
	"newrelic/multienv/pkg/connect"
	"newrelic/multienv/pkg/deser"
)

// Recever configuration
type RecvConfig struct {
	Connectors []connect.Connector
	Deser      deser.DeserFunc
}

// Processor configuration
type ProcConfig struct {
	Model any
}

type ExporterType string

func (expor ExporterType) Check() bool {
	switch expor {
	case NrApi:
	case NrMetrics:
	case NrEvents:
	case NrLogs:
	case NrTraces:
	case NrInfra:
	case Prometheus:
	case Otel:
	case Dummy:
	default:
		return false
	}
	return true
}

const (
	NrInfra    ExporterType = "nrinfra"
	NrApi      ExporterType = "nrapi"
	NrMetrics  ExporterType = "nrmetrics"
	NrEvents   ExporterType = "nrevents"
	NrLogs     ExporterType = "nrlogs"
	NrTraces   ExporterType = "nrtraces"
	Otel       ExporterType = "otel"
	Prometheus ExporterType = "prom"
	Dummy      ExporterType = "dummy"
)

// Data pipeline configuration.
type PipelineConfig struct {
	Interval uint
	Exporter ExporterType
	Custom   map[string]any
}

func (conf PipelineConfig) GetString(key string) (string, bool) {
	val, ok := conf.Custom[key]
	if ok {
		val_str, ok := val.(string)
		if ok {
			return val_str, true
		} else {
			return "", false
		}
	} else {
		return "", false
	}
}

func (conf PipelineConfig) GetInt(key string) (int, bool) {
	val, ok := conf.Custom[key]
	if ok {
		val_int, ok := val.(int)
		if ok {
			return val_int, true
		} else {
			return 0, false
		}
	} else {
		return 0, false
	}
}

func (conf PipelineConfig) GetMap(key string) (map[string]any, bool) {
	val, ok := conf.Custom[key]
	if ok {
		val_map, ok := val.(map[string]any)
		if ok {
			return val_map, true
		} else {
			return nil, false
		}
	} else {
		return nil, false
	}
}
