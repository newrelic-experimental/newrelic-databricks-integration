package main

import (
	"context"
	awslambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"newrelic/multienv/integration"
	"newrelic/multienv/pkg/config"
	nrilambda "newrelic/multienv/pkg/env/lambda"
	model "newrelic/multienv/pkg/model"
)

var pipeConf config.PipelineConfig
var initConf config.ProcConfig
var initErr error

func init() {
	pipeConf = nrilambda.LoadConfig()
	initConf, initErr = integration.InitProc(&pipeConf)
}

func HandleRequest(ctx context.Context, event map[string]any) (any, error) {
	if initErr != nil {
		log.Error("Error initializing = ", initErr)
		return nil, initErr
	}

	processorModel := initConf.Model
	responsePayload, ok := event["responsePayload"].([]any)

	var dataArray []model.MeltModel

	if ok {
		for _, payload := range responsePayload {
			err := mapstructure.Decode(payload, &processorModel)
			if err == nil {
				data := integration.Proc(processorModel)
				dataArray = append(dataArray, data...)
			} else {
				// We don't return an error because we don't want AWS to retry, just ignore this event.
				log.Error("Error mapping data = ", err)
			}
		}

		return dataArray, nil
	} else {
		log.Error("responsePayload not present")
		return nil, nil
	}
}

func main() {
	awslambda.Start(HandleRequest)
}
