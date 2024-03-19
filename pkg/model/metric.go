package model

import "time"

type MetricType int

const (
	Gauge MetricType = iota
	Count
	CumulativeCount
	Summary
	// TODO: add histogram metric type
)

// Metric model variant.
type MetricModel struct {
	Name     string
	Type     MetricType
	Value    Numeric
	Interval time.Duration
	//TODO: summary and histogram metric data model
}

// Make a gauge metric.
func MakeGaugeMetric(name string, value Numeric, timestamp time.Time) MeltModel {
	return MeltModel{
		Type:      Metric,
		Timestamp: timestamp.UnixMilli(),
		Data: MetricModel{
			Name:  name,
			Type:  Gauge,
			Value: value,
		},
	}
}

// Make a delta count metric.
func MakeCountMetric(name string, value Numeric, interval time.Duration, timestamp time.Time) MeltModel {
	return MeltModel{
		Type:      Metric,
		Timestamp: timestamp.UnixMilli(),
		Data: MetricModel{
			Name:     name,
			Type:     Count,
			Value:    value,
			Interval: interval,
		},
	}
}

// Make a cumulative count metric.
func MakeCumulativeCountMetric(name string, value Numeric, timestamp time.Time) MeltModel {
	return MeltModel{
		Type:      Metric,
		Timestamp: timestamp.UnixMilli(),
		Data: MetricModel{
			Name:  name,
			Type:  CumulativeCount,
			Value: value,
		},
	}
}

//TODO: make summary metric
