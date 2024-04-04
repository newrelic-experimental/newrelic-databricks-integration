package spark

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"newrelic/multienv/integration/Utils"
	"newrelic/multienv/pkg/config"
	"newrelic/multienv/pkg/connect"
	"newrelic/multienv/pkg/model"
	"reflect"
)

func GetSparkConnectors(pipeConfig *config.PipelineConfig) ([]connect.Connector, error) {
	recv_interval = int(pipeConfig.Interval)
	if recv_interval == 0 {
		log.Warn("Interval not set, using 60 seconds")
		recv_interval = 60
	}

	licenseKey := ""
	sparkEndpoint := ""

	if databricksLicenseKey, ok := pipeConfig.GetString("db_access_token"); ok {
		licenseKey = databricksLicenseKey
	} else {
		return nil, errors.New("config key 'db_access_token' doesn't exist")
	}

	if sparkEndpointUrl, ok := pipeConfig.GetString("spark_endpoint"); ok {
		sparkEndpoint = sparkEndpointUrl
	} else {
		return nil, errors.New("config key 'spark_endpoint' doesn't exist")
	}

	var headers = make(map[string]string)
	headers["Authorization"] = "Bearer " + licenseKey
	headers["Content-Type"] = "application/json"
	headers["Accept"] = "application/json"

	apps, _ := GetSparkApplications(sparkEndpoint, headers)
	fetchLogs()

	var connectors []connect.Connector

	for i := 0; i < len(apps); i++ {

		jobsConnector := connect.MakeHttpGetConnector(sparkEndpoint+"/api/v1/applications/"+apps[i].ID+"/jobs", headers)
		jobsConnector.SetConnectorModelName("SparkJob")

		executorsConnector := connect.MakeHttpGetConnector(sparkEndpoint+"/api/v1/applications/"+apps[i].ID+"/executors", headers)
		executorsConnector.SetConnectorModelName("SparkExecutor")

		stagesConnector := connect.MakeHttpGetConnector(sparkEndpoint+"/api/v1/applications/"+apps[i].ID+"/stages", headers)
		stagesConnector.SetConnectorModelName("SparkStage")

		connectors = append(connectors, &jobsConnector, &executorsConnector, &stagesConnector)
	}

	return connectors, nil
}

func InitSparkProc(pipeConfig *config.PipelineConfig) (config.ProcConfig, error) {
	recv_interval = int(pipeConfig.Interval)
	if recv_interval == 0 {
		log.Warn("Interval not set, using 60 seconds")
		recv_interval = 60
	}

	return config.ProcConfig{}, nil
}

// SparkProc Proc Generate all kinds of data.
func SparkProc(data any) []model.MeltModel {

	responseModel := data.(map[string]any)["model"]
	responseData := data.(map[string]any)["response"]

	if responseModel == "SparkStage" {

		var sparkStagesModel = SparkStage{}
		modelJson, _ := json.Marshal(responseData)
		err := json.Unmarshal(modelJson, &sparkStagesModel)
		if err != nil {
			return nil
		}

		e := reflect.ValueOf(&sparkStagesModel).Elem()
		tags := make(map[string]interface{})
		stagesTags := make(map[string]interface{})
		Utils.SetTags("spark.stage.", e, tags, stagesTags)

		return Utils.CreateMetricModels("spark.stage.", e, stagesTags)

	} else if responseModel == "SparkJob" {

		var sparkStagesModel = SparkJob{}
		modelJson, _ := json.Marshal(responseData)
		err := json.Unmarshal(modelJson, &sparkStagesModel)
		if err != nil {
			return nil
		}

		e := reflect.ValueOf(&sparkStagesModel).Elem()
		tags := make(map[string]interface{})
		stagesTags := make(map[string]interface{})
		Utils.SetTags("spark.job.", e, tags, stagesTags)

		return Utils.CreateMetricModels("spark.job.", e, stagesTags)

	} else if responseModel == "SparkExecutor" {

		var sparkStagesModel = SparkExecutor{}
		modelJson, _ := json.Marshal(responseData)
		err := json.Unmarshal(modelJson, &sparkStagesModel)
		if err != nil {
			return nil
		}

		e := reflect.ValueOf(&sparkStagesModel).Elem()
		tags := make(map[string]interface{})
		stagesTags := make(map[string]interface{})
		Utils.SetTags("spark.executor.", e, tags, stagesTags)

		return Utils.CreateMetricModels("spark.executor.", e, stagesTags)

	} else {
		log.Println("Unknown response model in Spark integration")
	}
	return nil
}
