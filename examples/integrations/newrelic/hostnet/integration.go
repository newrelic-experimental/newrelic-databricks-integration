package hostnet

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"newrelic/multienv/pkg/config"
	"newrelic/multienv/pkg/connect"
	"newrelic/multienv/pkg/deser"
	"newrelic/multienv/pkg/model"
)

///// New Relic GraphQL example, avarage host.net.receiveBytesPerSecond metric \\\\\

type nerdGraphResp struct {
	Data struct {
		Actor struct {
			Account struct {
				Nrql struct {
					Results []struct {
						Value float64 `mapstructure:"average.host.net.receiveBytesPerSecond"`
					}
				}
			}
		}
	}
}

type nrCred struct {
	AccountID string
	UserKey   string
}

// No danger of data races because it's only set on Init and after that only read
var recv_interval = 0

func InitRecv(pipeConfig *config.PipelineConfig) (config.RecvConfig, error) {
	recv_interval = int(pipeConfig.Interval)
	if recv_interval == 0 {
		log.Warn("NR Graph QL: Interval not set, using 5 seconds")
		recv_interval = 5
	}

	var nrCred nrCred

	if nr_account_id, ok := pipeConfig.GetString("nr_account_id"); ok {
		nrCred.AccountID = nr_account_id
	} else {
		return config.RecvConfig{}, errors.New("Config key 'nr_account_id' doesn't exist")
	}

	if nr_user_key, ok := pipeConfig.GetString("nr_user_key"); ok {
		nrCred.UserKey = nr_user_key
	} else {
		return config.RecvConfig{}, errors.New("Config key 'nr_user_key' doesn't exist")
	}

	url := "https://api.newrelic.com/graphql"
	query := "SELECT average(host.net.receiveBytesPerSecond) FROM Metric SINCE " + strconv.Itoa(recv_interval) + " seconds AGO"
	body := fmt.Sprintf(`{
		actor { account(id: %s) 
		{ nrql
		(query: "%s")
		{ results } } } 
	}`, nrCred.AccountID, query)

	headers := map[string]string{"API-Key": nrCred.UserKey}

	connector := connect.MakeHttpPostConnector(url, body, headers)
	connector.SetTimeout(10 * time.Second)

	return config.RecvConfig{
		Connectors: []connect.Connector{&connector},
		Deser:      deser.DeserJson,
	}, nil
}

func InitProc(pipeConfig *config.PipelineConfig) (config.ProcConfig, error) {
	return config.ProcConfig{
		Model: nerdGraphResp{},
	}, nil
}

// Processor function
func Proc(data any) []model.MeltModel {
	out := make([]model.MeltModel, 0)
	if nrdata, ok := data.(nerdGraphResp); ok {
		if len(nrdata.Data.Actor.Account.Nrql.Results) == 0 {
			return out
		}

		log.Printf("NR value received = %v\n", nrdata.Data.Actor.Account.Nrql.Results[0].Value)

		avrgVal := model.MakeFloatNumeric(nrdata.Data.Actor.Account.Nrql.Results[0].Value)
		avrgMetric := model.MakeGaugeMetric("nr.test.AvrgBytesPerSec", avrgVal, time.Now())

		// Simulate a counter metric
		interval := time.Duration(recv_interval) * time.Second
		randVal := model.MakeIntNumeric(int64(rand.Intn(10)))
		randMetric := model.MakeCountMetric("nr.test.Random", randVal, interval, time.Now())
		randMetric.Attributes = map[string]any{"range": "[0,10)"}

		out = append(out, avrgMetric, randMetric)
	} else {
		log.Println("Unknown type for NR data = ", data)
	}
	return out
}
