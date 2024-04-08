package integration

import (
	log "github.com/sirupsen/logrus"
	"newrelic/multienv/integration/databricks"
	"newrelic/multienv/integration/spark"
	"newrelic/multienv/pkg/config"
	"newrelic/multienv/pkg/connect"
	"newrelic/multienv/pkg/deser"
	"newrelic/multienv/pkg/model"
)

var recv_interval = 60

// InitRecv Integration Receiver Initializer
func InitRecv(pipeConfig *config.PipelineConfig) (config.RecvConfig, error) {

	recv_interval = int(pipeConfig.Interval)
	if recv_interval == 0 {
		log.Warn("Interval not set, using 60 seconds")
		recv_interval = 60
	}

	var connectors []connect.Connector

	sparkConnectors, _ := spark.GetSparkConnectors(pipeConfig)
	databricksConnectors, _ := databricks.GetDatabricksConnectors(pipeConfig)

	connectors = append(connectors, sparkConnectors...)
	connectors = append(connectors, databricksConnectors...)

	return config.RecvConfig{
		Connectors: connectors,
		Deser:      deser.DeserJson,
	}, nil
}

// InitProc Integration Processor Initializer
func InitProc(pipeConfig *config.PipelineConfig) (config.ProcConfig, error) {

	recv_interval = int(pipeConfig.Interval)
	if recv_interval == 0 {
		log.Warn("Interval not set, using 60 seconds")
		recv_interval = 60
	}

	_, sparkError := spark.InitSparkProc(pipeConfig)
	if sparkError != nil {
		return config.ProcConfig{}, sparkError
	}

	_, databricksError := databricks.InitDatabricksProc(pipeConfig)
	if databricksError != nil {
		return config.ProcConfig{}, databricksError
	}
	return config.ProcConfig{
		Model: map[string]any{},
	}, nil
}

// Integration Processor
func Proc(data any) []model.MeltModel {

	modelName := data.(map[string]any)["model"].(string)

	switch modelName {
	case "SparkJob", "SparkExecutor", "SparkStage":
		return spark.SparkProc(data)

	case "AWSDatabricksQueryList",
		"AWSDatabricksJobsRunsList",
		"GCPDatabricksQueryList",
		"GCPDatabricksJobsRunsList",
		"AzureDatabricksQueryList",
		"AzureDatabricksJobsRunsList":
		return databricks.DatabricksProc(data)

	default:
		log.Println("Unknown response model " + modelName)
	}

	return nil
}
