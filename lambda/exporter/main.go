package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	awslambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"newrelic/multienv/pkg/config"
	nrilambda "newrelic/multienv/pkg/env/lambda"
	"newrelic/multienv/pkg/export"
	"newrelic/multienv/pkg/model"
)

var pipeConf config.PipelineConfig

func init() {
	pipeConf = nrilambda.LoadConfig()
}

func HandleRequest(ctx context.Context, event events.SQSEvent) error {

	samples := make([]model.MeltModel, 0)
	var meltModel model.MeltModel

	for _, record := range event.Records {
		var obj map[string]any

		err := json.Unmarshal([]byte(record.Body), &obj)
		if err != nil {
			log.Error("Error unmarshalling = ", err)
			return err
		}

		responsePayload, ok := obj["responsePayload"].([]any)
		if ok {
			for _, payload := range responsePayload {

				mapErr := mapstructure.Decode(payload, &meltModel)
				if mapErr != nil {
					log.Warn("error while decoding", mapErr)
					return nil
				}

				switch meltModel.Type {
				case model.Metric:
					var metricModel model.MetricModel
					metricMapError := mapstructure.Decode(meltModel.Data, &metricModel)
					if metricMapError != nil {
						log.Warn("error while decoding metric", metricMapError)
						return nil
					}
					meltModel.Data = metricModel
				case model.Event:
					var eventModel model.EventModel
					eventMapError := mapstructure.Decode(meltModel.Data, &eventModel)
					if eventMapError != nil {
						log.Warn("error while decoding event", eventMapError)
						return nil
					}
					meltModel.Data = eventModel
				case model.Log:
					var logModel model.LogModel
					logMapError := mapstructure.Decode(meltModel.Data, &logModel)
					if logMapError != nil {
						log.Warn("error while decoding Log", logMapError)
						return nil
					}
					meltModel.Data = logModel
				case model.Trace:
					var traceModel model.TraceModel
					traceMapError := mapstructure.Decode(meltModel.Data, &traceModel)
					if traceMapError != nil {
						log.Warn("error while decoding Trace", traceMapError)
						return nil
					}
					meltModel.Data = traceModel
				case model.Custom:
					// Custom types should be handled explicitly based on requirement
					log.Println("Model is custom type")
				}
				samples = append(samples, meltModel)
			}

		} else {
			log.Error("Couldn't get 'responsePayload'")
			return nil
		}
	}

	//Export data
	exporter := export.SelectExporter(pipeConf.Exporter)
	err := exporter(pipeConf, samples)
	if err != nil {
		log.Error("Exporter failed = ", err)
		//TODO: handle error condition, refill buffer? Discard data? Retry?
	}

	return nil
}

func main() {
	awslambda.Start(HandleRequest)
}
