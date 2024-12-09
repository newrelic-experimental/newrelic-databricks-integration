package databricks

import (
	"time"

	"github.com/newrelic/newrelic-labs-sdk/v2/pkg/integration/model"
)

func writeGauge(
	prefix string,
	metricName string,
	metricValue any,
	attrs map[string]interface{},
	writer chan <- model.Metric,
) {
	metric := model.NewGaugeMetric(
		prefix + metricName,
		model.MakeNumeric(metricValue),
		time.Now(),
	)

	for k, v := range attrs {
		metric.Attributes[k] = v
	}

	writer <- metric
}

func makeAttributesMap(
	tags map[string]string,
) map[string]interface{} {
	attrs := make(map[string]interface{})

	for k, v := range tags {
		attrs[k] = v
	}

	return attrs
}
