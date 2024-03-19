package infra

import (
	"errors"
	"newrelic/multienv/pkg/config"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadConfig() (config.PipelineConfig, error) {
	confPath := os.Getenv("CONFIG_PATH")
	if confPath == "" {
		return config.PipelineConfig{}, errors.New("CONFIG_PATH environment variable is empty")
	}

	yamlFile, err := os.ReadFile(confPath)
	if err != nil {
		return config.PipelineConfig{}, err
	}

	var pipeConfigMap map[string]any
	err = yaml.Unmarshal(yamlFile, &pipeConfigMap)
	if err != nil {
		return config.PipelineConfig{}, err
	}

	var pipeConfig = config.PipelineConfig{
		Interval: 0,
		Exporter: config.ExporterType(config.NrInfra),
		Custom:   pipeConfigMap,
	}

	return pipeConfig, nil
}
