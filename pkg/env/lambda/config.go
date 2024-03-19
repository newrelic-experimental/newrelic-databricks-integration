package lambda

import (
	"os"
	"strconv"
	"strings"

	"newrelic/multienv/pkg/config"
)

func LoadConfig() config.PipelineConfig {
	pipeConf := config.PipelineConfig{
		Custom: map[string]any{},
	}

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		key := pair[0]
		val := pair[1]
		switch key {
		case "interval":
			interval, err := strconv.ParseUint(val, 10, 0)
			if err != nil {
				interval = 60
			}
			pipeConf.Interval = uint(interval)
		case "exporter":
			exporter := config.ExporterType(val)
			if exporter.Check() {
				pipeConf.Exporter = exporter
			}
		default:
			pipeConf.Custom[key] = val
		}
	}

	return pipeConf
}
