package spark

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"newrelic/multienv/integration/utils"
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

	var connectors []connect.Connector

	for i := 0; i < len(apps); i++ {

		var appData map[string]interface{}
		appData = make(map[string]interface{})
		appData["appId"] = apps[i].ID
		appData["appName"] = apps[i].Name

		jobsConnector := connect.MakeHttpGetConnector(sparkEndpoint+"/api/v1/applications/"+apps[i].ID+"/jobs", headers)
		jobsConnector.SetConnectorModelName("SparkJob")
		jobsConnector.SetCustomData(appData)

		executorsConnector := connect.MakeHttpGetConnector(sparkEndpoint+"/api/v1/applications/"+apps[i].ID+"/executors", headers)
		executorsConnector.SetConnectorModelName("SparkExecutor")
		executorsConnector.SetCustomData(appData)

		stagesConnector := connect.MakeHttpGetConnector(sparkEndpoint+"/api/v1/applications/"+apps[i].ID+"/stages", headers)
		stagesConnector.SetConnectorModelName("SparkStage")
		stagesConnector.SetCustomData(appData)

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
	customData := data.(map[string]any)["customData"].(map[string]interface{})

	switch responseModel {
	case "SparkStage":
		var sparkStagesModel = SparkStage{}
		tagPrefix := "spark.stage."
		modelJson, _ := json.Marshal(responseData)
		err := json.Unmarshal(modelJson, &sparkStagesModel)
		if err != nil {
			return nil
		}

		e := reflect.ValueOf(&sparkStagesModel).Elem()
		tags := make(map[string]interface{})
		tags["appId"] = customData["appId"]
		tags["appName"] = customData["appName"]
		stagesTags := make(map[string]interface{})
		utils.SetTags(tagPrefix, e, tags, stagesTags)
		return utils.CreateMetricModels(tagPrefix, e, stagesTags)

	case "SparkJob":

		var sparkJobsModel = SparkJob{}
		tagPrefix := "spark.job."
		modelJson, _ := json.Marshal(responseData)
		err := json.Unmarshal(modelJson, &sparkJobsModel)
		if err != nil {
			return nil
		}

		e := reflect.ValueOf(&sparkJobsModel).Elem()
		tags := make(map[string]interface{})
		tags["appId"] = customData["appId"]
		tags["appName"] = customData["appName"]
		sparkJobTags := make(map[string]interface{})
		utils.SetTags(tagPrefix, e, tags, sparkJobTags)
		return utils.CreateMetricModels(tagPrefix, e, sparkJobTags)

	case "SparkExecutor":

		var sparkExecutorsModel = SparkExecutor{}
		tagPrefix := "spark.executor."
		modelJson, _ := json.Marshal(responseData)
		err := json.Unmarshal(modelJson, &sparkExecutorsModel)
		if err != nil {
			return nil
		}

		e := reflect.ValueOf(&sparkExecutorsModel).Elem()
		tags := make(map[string]interface{})
		tags["appId"] = customData["appId"]
		tags["appName"] = customData["appName"]
		sparkExecutorTags := make(map[string]interface{})
		utils.SetTags(tagPrefix, e, tags, sparkExecutorTags)
		return utils.CreateMetricModels(tagPrefix, e, sparkExecutorTags)

	default:
		log.Println("Unknown response model in Spark integration")
	}
	return nil
}
