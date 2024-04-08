package standalone

import (
	"sync"
	"time"

	"newrelic/multienv/pkg/connect"
	"newrelic/multienv/pkg/deser"

	log "github.com/sirupsen/logrus"
)

type RecvWorkerConfig struct {
	IntervalSec  uint
	Connectors   []connect.Connector
	Deserializer deser.DeserFunc
	OutChannel   chan<- map[string]any
}

var recvWorkerConfigHoldr SharedConfig[RecvWorkerConfig]

func InitReceiver(config RecvWorkerConfig) {
	recvWorkerConfigHoldr.SetConfig(config)
	if !recvWorkerConfigHoldr.SetIsRunning() {
		log.Println("Starting receiver worker...")
		go receiverWorker()
	} else {
		log.Println("Receiver worker already running, config updated.")
	}
}

func receiverWorker() {
	for {
		config := recvWorkerConfigHoldr.Config()
		pre := time.Now().Unix()

		wg := &sync.WaitGroup{}
		wg.Add(len(config.Connectors))

		for _, connector := range config.Connectors {
			go func(connector connect.Connector) {
				defer wg.Done()

				data, err := connector.Request()
				if err.Err != nil {
					log.Error("Http Get error = ", err.Err.Error())
					delayBeforeNextReq(pre, &config)
					return
				}

				//log.Println("Data received: ", string(data))

				response, jsonType, desErr := config.Deserializer(data)
				if desErr == nil {
					switch jsonType {
					case "array":
						var rdata = map[string]any{}
						for _, responseObject := range response.([]map[string]any) {
							rdata["model"] = connector.ConnectorModel()
							rdata["customData"] = connector.ConnectorCustomData()
							rdata["response"] = responseObject
							config.OutChannel <- rdata
						}
					case "object":
						var rdata = map[string]any{}
						rdata["model"] = connector.ConnectorModel()
						rdata["customData"] = connector.ConnectorCustomData()
						rdata["response"] = response.(map[string]any)
						config.OutChannel <- rdata
					default:
						log.Warn("Response data type couldn't be identified")
					}
				}
			}(connector)
		}

		wg.Wait()
		log.Println("All Requests Completed")

		// Delay before the next request
		delayBeforeNextReq(pre, &config)
	}
}

func delayBeforeNextReq(pre int64, config *RecvWorkerConfig) {
	timeDiff := time.Now().Unix() - pre
	if timeDiff < int64(config.IntervalSec) {
		remainingDelay := int64(config.IntervalSec) - timeDiff
		time.Sleep(time.Duration(remainingDelay) * time.Second)
	}
}
