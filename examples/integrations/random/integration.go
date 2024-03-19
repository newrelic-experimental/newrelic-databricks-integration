package random

import (
	"fmt"
	"math/rand"
	"newrelic/multienv/pkg/config"
	"newrelic/multienv/pkg/connect"
	"newrelic/multienv/pkg/deser"
	"newrelic/multienv/pkg/model"
	"time"

	log "github.com/sirupsen/logrus"
)

var recv_interval = 0

type dummyConnector struct{}

func (c *dummyConnector) SetConfig(config any) {
}

func (c *dummyConnector) Request() ([]byte, connect.ConnecError) {
	return []byte("{}"), connect.ConnecError{}
}

func (c *dummyConnector) ConnectorID() string {
	return "dummy"
}

func InitRecv(pipeConfig *config.PipelineConfig) (config.RecvConfig, error) {
	recv_interval = int(pipeConfig.Interval)
	if recv_interval == 0 {
		log.Warn("Random: Interval not set, using 5 seconds")
		recv_interval = 5
	}
	return config.RecvConfig{
		Connectors: []connect.Connector{&dummyConnector{}},
		Deser:      deser.DeserJson,
	}, nil
}

func InitProc(pipeConfig *config.PipelineConfig) (config.ProcConfig, error) {
	return config.ProcConfig{}, nil
}

// Generate all kinds of data.
func Proc(data any) []model.MeltModel {
	log.Print("Radom integration")

	randVal := model.MakeIntNumeric(int64(rand.Intn(100)))

	gauge := model.MakeGaugeMetric("rnd.gauge", randVal, time.Now())
	gauge.Attributes = map[string]any{"test.attr": 1001}

	return []model.MeltModel{
		model.MakeLog("Random number is "+fmt.Sprint(randVal.Int()), "Random", time.Now()),
		model.MakeCountMetric("rnd.count", randVal, time.Duration(recv_interval)*time.Second, time.Now()),
		model.MakeCumulativeCountMetric("time.count", model.MakeIntNumeric(time.Now().Unix()), time.Now()),
		gauge,
		model.MakeEvent("randomEvent", map[string]any{"rnd": randVal.Int()}, time.Now()),
	}
}
