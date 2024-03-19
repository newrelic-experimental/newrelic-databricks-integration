package standalone

import (
	"errors"
	"os"
	"sync"
	"sync/atomic"

	"newrelic/multienv/pkg/config"

	"gopkg.in/yaml.v3"
)

type SharedConfig[T any] struct {
	workerConfig    T
	workerConfigMu  sync.Mutex
	workerIsRunning atomic.Bool
}

func (w *SharedConfig[T]) SetIsRunning() bool {
	return w.workerIsRunning.Swap(true)
}

func (w *SharedConfig[T]) SetConfig(config T) {
	w.workerConfigMu.Lock()
	w.workerConfig = config
	w.workerConfigMu.Unlock()
}

func (w *SharedConfig[T]) Config() T {
	w.workerConfigMu.Lock()
	config := w.workerConfig
	w.workerConfigMu.Unlock()
	return config
}

func LoadConfig(filePath string) (config.PipelineConfig, error) {
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return config.PipelineConfig{}, err
	}

	var pipeConfigMap map[string]any
	err = yaml.Unmarshal(yamlFile, &pipeConfigMap)
	if err != nil {
		return config.PipelineConfig{}, err
	}

	interval, ok := pipeConfigMap["interval"]
	if !ok {
		interval = 60
	}

	if _, ok := interval.(int); !ok {
		return config.PipelineConfig{}, errors.New("Interval must be an integer")
	}

	if interval.(int) <= 0 {
		interval = 60
	}

	exporter, ok := pipeConfigMap["exporter"]
	if !ok {
		return config.PipelineConfig{}, errors.New("Exporter must be specified")
	}

	if _, ok := exporter.(string); !ok {
		return config.PipelineConfig{}, errors.New("Exporter must be a string")
	}

	if !config.ExporterType(exporter.(string)).Check() {
		return config.PipelineConfig{}, errors.New("Invalid 'exporter' value in config.")
	}

	delete(pipeConfigMap, "interval")
	delete(pipeConfigMap, "exporter")

	var pipeConfig = config.PipelineConfig{
		Interval: uint(interval.(int)),
		Exporter: config.ExporterType(exporter.(string)),
		Custom:   pipeConfigMap,
	}

	return pipeConfig, nil
}
