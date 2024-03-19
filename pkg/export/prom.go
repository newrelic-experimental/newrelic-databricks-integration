package export

import (
	"context"
	"errors"
	"fmt"
	"newrelic/multienv/pkg/config"
	"newrelic/multienv/pkg/model"
	"regexp"
	"strconv"
	"time"

	"github.com/castai/promwrite"

	log "github.com/sirupsen/logrus"
)

func exportProm(pipeConf config.PipelineConfig, data []model.MeltModel) error {
	log.Print("------> Prometheus Exporter = ", data)

	endpoint, ok := pipeConf.GetString("prom_endpoint")
	if !ok {
		return errors.New("Config key 'prom_endpoint' doesn't exist")
	}

	metrics := make([]promwrite.TimeSeries, 0)

	for _, m := range data {
		switch m.Type {
		case model.Metric:
			timeSeries, ok := metricIntoProm(m)
			if ok {
				metrics = append(metrics, timeSeries)
			} else {
				log.Warn("Metric can't be converted to Prometheus format")
			}
		default:
			log.Warn("Data sample is not a metric, can't be sent to Prometheus")
		}
	}

	if len(metrics) == 0 {
		return nil
	}

	promHeaders := getPromHeaders(pipeConf)
	credentials, ok_credentials := pipeConf.GetString("prom_credentials")
	if ok_credentials {
		promHeaders["Authorization"] = "Bearer " + credentials
	}

	client := promwrite.NewClient(endpoint)
	resp, err := client.Write(context.Background(), &promwrite.WriteRequest{TimeSeries: metrics}, promwrite.WriteHeaders(promHeaders))

	if err != nil {
		return err
	}

	log.Print("Prom response = ", resp)

	return nil
}

func metricIntoProm(melt model.MeltModel) (promwrite.TimeSeries, bool) {
	metric, ok := melt.Metric()
	if !ok {
		return promwrite.TimeSeries{}, false
	}

	labels := attributesToPromLabels(melt.Attributes)
	labels = append(labels, promwrite.Label{
		Name:  "__name__",
		Value: nameToProm(metric.Name),
	})

	switch metric.Type {
	// We are just ignoring the interval, because prometheus doesn't support delta counters
	case model.Gauge, model.Count, model.CumulativeCount:
		return promwrite.TimeSeries{
			Labels: labels,
			Sample: promwrite.Sample{
				Time:  time.UnixMilli(melt.Timestamp),
				Value: metric.Value.Float(),
			},
		}, true
	default:
		return promwrite.TimeSeries{}, false
	}
}

func attributesToPromLabels(attr map[string]any) []promwrite.Label {
	labels := []promwrite.Label{}
	for k, v := range attr {
		switch val := v.(type) {
		case string:
			labels = append(labels, promwrite.Label{
				Name:  nameToProm(k),
				Value: val,
			})
		case int:
			labels = append(labels, promwrite.Label{
				Name:  nameToProm(k),
				Value: strconv.Itoa(val),
			})
		case float32:
			labels = append(labels, promwrite.Label{
				Name:  nameToProm(k),
				Value: strconv.FormatFloat(float64(val), 'f', 2, 32),
			})
		case float64:
			labels = append(labels, promwrite.Label{
				Name:  nameToProm(k),
				Value: strconv.FormatFloat(val, 'f', 2, 32),
			})
		case fmt.Stringer:
			labels = append(labels, promwrite.Label{
				Name:  nameToProm(k),
				Value: val.String(),
			})
		default:
			log.Warn("Attribute of unsupported type: ", k, v)
		}
	}
	return labels
}

// Convert name into prometheus naming conventions: only allow [a-zA-Z0-9_:]
func nameToProm(name string) string {
	namePattern := regexp.MustCompile(`[^a-zA-Z0-9_:]`)
	return string(namePattern.ReplaceAll([]byte(name), []byte("_")))
}

func getPromHeaders(pipeConf config.PipelineConfig) map[string]string {
	any_headers, ok := pipeConf.GetMap("prom_headers")
	if ok {
		headers := map[string]string{}
		for k, v := range any_headers {
			switch v := v.(type) {
			case string:
				headers[k] = v
			default:
				log.Warn("Found a non string value in 'prom_headers', ignoring")
			}
		}
		return headers
	} else {
		return map[string]string{}
	}
}
