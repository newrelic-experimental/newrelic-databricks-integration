package ipify

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"newrelic/multienv/pkg/config"
	"newrelic/multienv/pkg/connect"
	"newrelic/multienv/pkg/deser"
	"newrelic/multienv/pkg/model"
)

///// Ipify API example \\\\\

type ipify struct {
	IpAddress string `mapstructure:"ip"`
}

func InitRecvSimple(pipeConfig *config.PipelineConfig) (config.RecvConfig, error) {
	example1Connector := connect.MakeHttpGetConnector("https://api.ipify.org/?format=json", nil)
	example2Connector := connect.MakeHttpGetConnector("https://api.ipify.org/?format=json", nil)

	return config.RecvConfig{
		Connectors: []connect.Connector{&example1Connector, &example2Connector},
		Deser:      deser.DeserJson,
	}, nil
}

func InitRecvWithReqBuilder(pipeConfig *config.PipelineConfig) (config.RecvConfig, error) {
	example1Connector := connect.MakeHttpConnectorWithBuilder(example1RequestBuilder)
	example2Connector := connect.MakeHttpConnectorWithBuilder(example2RequestBuilder)

	return config.RecvConfig{
		Connectors: []connect.Connector{&example1Connector, &example2Connector},
		Deser:      deser.DeserJson,
	}, nil
}

func example2RequestBuilder(conf *connect.HttpConfig) (*http.Request, error) {
	req, err := http.NewRequest("GET", "https://api.ipify.org/?format=json", nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// Custom request builder. Not necessary, only to show how it works.
func example1RequestBuilder(conf *connect.HttpConfig) (*http.Request, error) {
	req, err := http.NewRequest("GET", "https://api.ipify.org/?format=json", nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func InitProc(pipeConfig *config.PipelineConfig) (config.ProcConfig, error) {
	return config.ProcConfig{
		Model: ipify{},
	}, nil
}

// Processor function
func Proc(data any) []model.MeltModel {
	out := make([]model.MeltModel, 0)
	if ipify, ok := data.(ipify); ok {
		log.Println("My IP is = " + ipify.IpAddress)
		mlog := model.MakeLog(ipify.IpAddress, "IPAddress", time.Now())
		mlog.Attributes = map[string]any{"type": "ip"}
		out = append(out, mlog)
	} else {
		log.Warn("Unknown type for data = ", data)
	}
	return out
}
