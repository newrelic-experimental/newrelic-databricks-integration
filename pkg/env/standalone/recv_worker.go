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

				log.Println("Data received: ", string(data))

				var deserBuffer map[string]any
				desErr := config.Deserializer(data, &deserBuffer)
				if desErr == nil {
					config.OutChannel <- deserBuffer
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
