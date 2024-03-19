package status

import (
	"newrelic/multienv/pkg/config"
	"newrelic/multienv/pkg/connect"
	"newrelic/multienv/pkg/deser"
	"newrelic/multienv/pkg/model"
	"time"

	log "github.com/sirupsen/logrus"
)

///// New Relic Status API example \\\\\

type nrStatus struct {
	Page struct {
		Id         string
		Name       string
		Url        string
		Time_zone  string
		Updated_at string
	}
	Status struct {
		Indicator   string
		Description string
	}
}

func InitRecv(pipeConfig *config.PipelineConfig) (config.RecvConfig, error) {
	connector := connect.MakeHttpGetConnector("https://status.newrelic.com/api/v2/status.json", nil)

	return config.RecvConfig{
		Connectors: []connect.Connector{&connector},
		Deser:      deser.DeserJson,
	}, nil
}

func InitProc(pipeConfig *config.PipelineConfig) (config.ProcConfig, error) {
	return config.ProcConfig{
		Model: nrStatus{},
	}, nil
}

// Processor function
func Proc(data any) []model.MeltModel {
	out := make([]model.MeltModel, 0)
	if nrStatus, ok := data.(nrStatus); ok {
		mlog := model.MakeLog(nrStatus.Status.Description, "NRStatus", time.Now())
		mlog.Attributes = map[string]any{
			"updatedAt": nrStatus.Page.Updated_at,
			"indicator": nrStatus.Status.Indicator,
			"id":        nrStatus.Page.Id,
		}
		out = append(out, mlog)
	} else {
		log.Warn("Unknown type for data = ", data)
	}
	return out
}
