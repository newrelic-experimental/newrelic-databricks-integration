package main

import (
	"context"
	awslambda "github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
	"newrelic/multienv/integration"
	"newrelic/multienv/pkg/config"
	"newrelic/multienv/pkg/connect"
	nrilambda "newrelic/multienv/pkg/env/lambda"
	"sync"
)

var pipeConf config.PipelineConfig
var initConf config.RecvConfig
var initErr error

func init() {
	pipeConf = nrilambda.LoadConfig()
	initConf, initErr = integration.InitRecv(&pipeConf)
}

// TODO: maybe we can get the scheduler rule from the context
// https://docs.aws.amazon.com/lambda/latest/dg/golang-context.html

func HandleRequest(ctx context.Context, event any) ([]map[string]any, error) {

	if initErr != nil {
		log.Error("Error initializing = ", initErr)
		return nil, initErr
	}

	wg := &sync.WaitGroup{}
	var deserBuffers []map[string]any

	for _, connector := range initConf.Connectors {
		wg.Add(1)

		go func(connector connect.Connector) {
			defer wg.Done()

			var deserBuffer map[string]any

			data, reqErr := connector.Request()
			if reqErr.Err != nil {
				log.Error("Http Get error = ", reqErr.Err.Error())
			}

			desErr := initConf.Deser(data, &deserBuffer)
			if desErr != nil {
				log.Error("Error deserializing data = ", desErr)
			}

			log.Print("Data received: ", deserBuffer)
			deserBuffers = append(deserBuffers, deserBuffer)
		}(connector)
	}
	wg.Wait()

	return deserBuffers, nil
}

func main() {
	awslambda.Start(HandleRequest)
}
