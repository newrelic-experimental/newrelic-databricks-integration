package integration

import (
	"newrelic/multienv/examples/integrations/random"
	"newrelic/multienv/pkg/config"
	"newrelic/multienv/pkg/model"
)

// Integration Receiver Initializer
func InitRecv(pipeConfig *config.PipelineConfig) (config.RecvConfig, error) {
	// CALL YOUR RECEIVER INITIALIZER HERE
	return random.InitRecv(pipeConfig)
}

// Integration Processor Initializer
func InitProc(pipeConfig *config.PipelineConfig) (config.ProcConfig, error) {
	// CALL YOUR PROCESSOR INITIALIZER HERE
	return random.InitProc(pipeConfig)
}

// Integration Processor
func Proc(data any) []model.MeltModel {
	// CALL YOUR PROCESSOR HERE
	return random.Proc(data)
}
