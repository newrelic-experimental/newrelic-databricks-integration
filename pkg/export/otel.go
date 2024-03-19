package export

import (
	"encoding/json"
	"newrelic/multienv/pkg/config"
	"newrelic/multienv/pkg/connect"
	"newrelic/multienv/pkg/model"

	log "github.com/sirupsen/logrus"
)

//// OTel data model for Logs

type otelLogsData struct {
	ResourceLogs []otelResourceLog `json:"resourceLogs"`
}

type otelResourceLog struct {
	//TODO: add "resource" model
	ScopeLogs []otelScopeLog `json:"scopeLogs"`
}

type otelScopeLog struct {
	//TODO: add "scope" model
	LogRecords []otelLogRecord `json:"logRecords"`
}

type otelLogRecord struct {
	TimeUnixNano         int64           `json:"timeUnixNano"`
	ObservedTimeUnixNano int64           `json:"observedTimeUnixNano"`
	SeverityText         string          `json:"severityText"`
	Body                 stringAttribute `json:"body"`
	Attributes           []attribute     `json:"attributes"`
}

//// OTel data model for Metrics

type otelMetricsData struct {
	ResourceMetrics []otelResourceMetrics `json:"resourceMetrics"`
}

type otelResourceMetrics struct {
	//TODO: add "resource" model
	ScopeMetrics []otelScopeMetrics `json:"scopeMetrics"`
}

type otelScopeMetrics struct {
	//TODO: add "scope" model
	Metrics []any `json:"metrics"`
}

type otelMetricSum struct {
	Name string            `json:"name"`
	Sum  otelMetricSumData `json:"sum"`
}

type otelAggrTemp int

const (
	AggrTempUnspecified otelAggrTemp = 0
	AggrTempDelta       otelAggrTemp = 1
	AggrTempCumulative  otelAggrTemp = 2
)

type otelMetricSumData struct {
	AggregationTemporality otelAggrTemp          `json:"aggregationTemporality"`
	IsMonotonic            bool                  `json:"isMonotonic"`
	DataPoints             []otelNumberDataPoint `json:"dataPoints"`
}

type otelMetricGauge struct {
	Name  string              `json:"name"`
	Gauge otelMetricGaugeData `json:"gauge"`
}

type otelMetricGaugeData struct {
	DataPoints []otelNumberDataPoint `json:"dataPoints"`
}

type otelNumberDataPoint struct {
	//TODO: also include AsInt
	AsDouble          float64     `json:"asDouble"`
	TimeUnixNano      int64       `json:"timeUnixNano"`
	StartTimeUnixNano int64       `json:"startTimeUnixNano,omitempty"`
	Attributes        []attribute `json:"attributes"`
}

//TODO: summary metric model

//// Common OTel models

//TODO: resource model
//TODO: scope model

type attribute struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

type stringAttribute struct {
	StringValue string `json:"stringValue"`
}

type intAttribute struct {
	IntValue int64 `json:"intValue"`
}

type doubleAttribute struct {
	DoubleValue float64 `json:"doubleValue"`
}

type boolAttribute struct {
	BoolValue bool `json:"boolValue"`
}

func exportOtel(pipeConf config.PipelineConfig, data []model.MeltModel) error {
	metrics := make([]model.MeltModel, 0)
	logs := make([]model.MeltModel, 0)

	for _, m := range data {
		switch m.Type {
		case model.Metric:
			metrics = append(metrics, m)
		case model.Event, model.Log:
			logs = append(logs, m)
		case model.Trace:
			//TODO: append traces
		}
	}

	if len(metrics) > 0 {
		err := exportOtelMetrics(pipeConf, metrics)
		if err != nil {
			return err
		}
	}

	if len(logs) > 0 {
		err := exportOtelLogs(pipeConf, logs)
		if err != nil {
			return err
		}
	}

	return nil
}

func exportOtelMetrics(pipeConf config.PipelineConfig, data []model.MeltModel) error {
	log.Print("------> OpenTelemetry Metric Exporter = ", data)

	metrics := make([]any, 0)

	// Generate metrics
	for _, d := range data {
		m, _ := d.Metric()
		switch m.Type {
		case model.Gauge:
			gauge := buildOtelGaugeMetric(d)
			metrics = append(metrics, gauge)
		case model.Count:
			count := buildOtelCountMetric(d)
			metrics = append(metrics, count)
		case model.CumulativeCount:
			count := buildOtelCumulCountMetric(d)
			metrics = append(metrics, count)
		case model.Summary:
			//TODO
		}
	}

	metricData := otelMetricsData{
		ResourceMetrics: []otelResourceMetrics{
			{
				ScopeMetrics: []otelScopeMetrics{
					{
						Metrics: metrics,
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(metricData)
	if err != nil {
		return err
	}

	log.Print("Metric data in JSON = ", string(jsonData))

	return otlpMetricRequest(pipeConf, jsonData)
}

func buildOtelGaugeMetric(melt model.MeltModel) otelMetricGauge {
	metric, _ := melt.Metric()
	return otelMetricGauge{
		Name: metric.Name,
		Gauge: otelMetricGaugeData{
			DataPoints: []otelNumberDataPoint{
				{
					AsDouble:          metric.Value.Float(),
					TimeUnixNano:      melt.Timestamp * 1000000,
					StartTimeUnixNano: 0,
					Attributes:        convertAttributes(melt.Attributes),
				},
			},
		},
	}
}

func buildOtelCountMetric(melt model.MeltModel) otelMetricSum {
	metric, _ := melt.Metric()
	return otelMetricSum{
		Name: metric.Name,
		Sum: otelMetricSumData{
			IsMonotonic:            true,
			AggregationTemporality: AggrTempDelta,
			DataPoints: []otelNumberDataPoint{
				{
					AsDouble:          metric.Value.Float(),
					TimeUnixNano:      melt.Timestamp * 1000000,
					StartTimeUnixNano: (melt.Timestamp - metric.Interval.Milliseconds()) * 1000000,
					Attributes:        convertAttributes(melt.Attributes),
				},
			},
		},
	}
}

func buildOtelCumulCountMetric(melt model.MeltModel) otelMetricSum {
	metric, _ := melt.Metric()
	return otelMetricSum{
		Name: metric.Name,
		Sum: otelMetricSumData{
			IsMonotonic:            true,
			AggregationTemporality: AggrTempCumulative,
			DataPoints: []otelNumberDataPoint{
				{
					AsDouble:     metric.Value.Float(),
					TimeUnixNano: melt.Timestamp * 1000000,
					Attributes:   convertAttributes(melt.Attributes),
				},
			},
		},
	}
}

func exportOtelLogs(pipeConf config.PipelineConfig, data []model.MeltModel) error {
	log.Print("------> OpenTelemetry Log Exporter = ", data)

	logRecords := []otelLogRecord{}
	for _, d := range data {
		l, _ := d.Log()
		logRecord := otelLogRecord{
			TimeUnixNano:         d.Timestamp * 1000000,
			ObservedTimeUnixNano: d.Timestamp * 1000000,
			SeverityText:         l.Type,
			Body: stringAttribute{
				StringValue: l.Message,
			},
			Attributes: convertAttributes(d.Attributes),
		}
		logRecords = append(logRecords, logRecord)
	}

	logsData := otelLogsData{
		ResourceLogs: []otelResourceLog{
			{
				ScopeLogs: []otelScopeLog{
					{
						LogRecords: logRecords,
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(logsData)
	if err != nil {
		return err
	}

	log.Print("Logs data in JSON = ", string(jsonData))

	return otlpLogRequest(pipeConf, jsonData)
}

func convertAttributes(attr map[string]any) []attribute {
	attrArr := []attribute{}
	for k, v := range attr {
		switch v := v.(type) {
		case string:
			a := attribute{
				Key: k,
				Value: stringAttribute{
					StringValue: v,
				},
			}
			attrArr = append(attrArr, a)
		case int:
			a := attribute{
				Key: k,
				Value: intAttribute{
					IntValue: int64(v),
				},
			}
			attrArr = append(attrArr, a)
		case float64:
			a := attribute{
				Key: k,
				Value: doubleAttribute{
					DoubleValue: v,
				},
			}
			attrArr = append(attrArr, a)
		case bool:
			a := attribute{
				Key: k,
				Value: boolAttribute{
					BoolValue: v,
				},
			}
			attrArr = append(attrArr, a)
		}
	}
	return attrArr
}

func otlpMetricRequest(pipeConf config.PipelineConfig, jsonModel []byte) error {
	endpoint, ok := pipeConf.GetString("otel_metric_endpoint")
	if !ok {
		endpoint = getOtelScheme(pipeConf) + "://" + getOtelEndpoint(pipeConf) + "/v1/metrics"
	}
	return otlpRequest(pipeConf, jsonModel, endpoint)
}

func otlpLogRequest(pipeConf config.PipelineConfig, jsonModel []byte) error {
	endpoint, ok := pipeConf.GetString("otel_log_endpoint")
	if !ok {
		endpoint = getOtelScheme(pipeConf) + "://" + getOtelEndpoint(pipeConf) + "/v1/logs"
	}
	return otlpRequest(pipeConf, jsonModel, endpoint)
}

func otlpRequest(pipeConf config.PipelineConfig, jsonModel []byte, endpoint string) error {
	headers := getOtelHeaders(pipeConf)
	headers["Content-Type"] = "application/json"
	headers["Content-Encoding"] = "gzip"

	gzipBody, errGzip := gzipString(string(jsonModel))
	if errGzip != nil {
		log.Error("Error compressing body = ", errGzip)
		return errGzip
	}

	connector := connect.MakeHttpPostConnector(endpoint, gzipBody, headers)
	response, errReq := connector.Request()
	if errReq.Err != nil {
		log.Error("Error sending request to OTel collector = ", errReq)
		return errReq.Err
	}

	log.Print("Response from OTel collector = ", string(response))

	return nil
}

func getOtelEndpoint(pipeConf config.PipelineConfig) string {
	endpoint, ok := pipeConf.GetString("otel_endpoint")
	if ok {
		return endpoint
	} else {
		log.Warn("'otel_endpoint' not specified, fallback to 'otlp.nr-data.net:4318'")
		return "otlp.nr-data.net:4318"
	}
}

func getOtelScheme(pipeConf config.PipelineConfig) string {
	scheme, ok := pipeConf.GetString("otel_scheme")
	if ok {
		return scheme
	} else {
		return "https"
	}
}

func getOtelHeaders(pipeConf config.PipelineConfig) map[string]string {
	any_headers, ok := pipeConf.GetMap("otel_headers")
	if ok {
		headers := map[string]string{}
		for k, v := range any_headers {
			switch v := v.(type) {
			case string:
				headers[k] = v
			default:
				log.Warn("Found a non string value in 'otel_headers', ignoring")
			}
		}
		return headers
	} else {
		return map[string]string{}
	}
}
