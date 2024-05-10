package databricks

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"newrelic/multienv/integration/utils"
	"newrelic/multienv/pkg/config"
	"newrelic/multienv/pkg/connect"
	"newrelic/multienv/pkg/model"
	"reflect"
	"time"
)

func GetDatabricksConnectors(pipeConfig *config.PipelineConfig) ([]connect.Connector, error) {

	licenseKey := ""
	databricksEndpoint := ""
	databricksCloudProvider := ""

	if databricksLicenseKey, ok := pipeConfig.GetString("db_access_token"); ok {
		licenseKey = databricksLicenseKey
	} else {
		return nil, errors.New("config key 'db_access_token' doesn't exist")
	}

	if databricksEndpointUrl, ok := pipeConfig.GetString("databricks_endpoint"); ok {
		databricksEndpoint = databricksEndpointUrl
	} else {
		return nil, errors.New("config key 'databricks_endpoint' doesn't exist")
	}

	if databricksCloud, ok := pipeConfig.GetString("db_cloud"); ok {
		databricksCloudProvider = databricksCloud
	} else {
		return nil, errors.New("config key 'db_cloud' doesn't exist")
	}

	var headers = make(map[string]string)
	headers["Authorization"] = "Bearer " + licenseKey
	headers["Content-Type"] = "application/json"
	headers["Accept"] = "application/json"

	var connectors []connect.Connector

	queryHistoryConnector := connect.MakeHttpGetConnector(databricksEndpoint+"/api/2.0/sql/history/queries?include_metrics=true", headers)

	if databricksCloudProvider == "GCP" {
		queryHistoryConnector.SetConnectorModelName("GCPDatabricksQueryList")
	} else if databricksCloudProvider == "AWS" {
		queryHistoryConnector.SetConnectorModelName("AWSDatabricksQueryList")
	} else if databricksCloudProvider == "AZURE" {
		queryHistoryConnector.SetConnectorModelName("AzureDatabricksQueryList")
	}

	jobRunsListConnector := connect.MakeHttpGetConnector(databricksEndpoint+"/api/2.1/jobs/runs/list", headers)
	if databricksCloudProvider == "GCP" {
		jobRunsListConnector.SetConnectorModelName("GCPDatabricksJobsRunsList")
	} else if databricksCloudProvider == "AWS" {
		jobRunsListConnector.SetConnectorModelName("AWSDatabricksJobsRunsList")
	} else if databricksCloudProvider == "AZURE" {
		jobRunsListConnector.SetConnectorModelName("AzureDatabricksJobsRunsList")
	}

	pipelinesListConnector := connect.MakeHttpGetConnector(databricksEndpoint+"/api/2.0/pipelines?max_results=100", headers)
	if databricksCloudProvider == "GCP" {
		pipelinesListConnector.SetConnectorModelName("GCPDatabricksPipelinesList")
	} else if databricksCloudProvider == "AWS" {
		pipelinesListConnector.SetConnectorModelName("AWSDatabricksPipelinesList")
	} else if databricksCloudProvider == "AZURE" {
		pipelinesListConnector.SetConnectorModelName("AzureDatabricksPipelinesList")
	}

	connectors = append(connectors, &jobRunsListConnector, &queryHistoryConnector)

	return connectors, nil
}

func InitDatabricksProc(pipeConfig *config.PipelineConfig) (config.ProcConfig, error) {

	recv_interval = int(pipeConfig.Interval)
	if recv_interval == 0 {
		log.Warn("Interval not set, using 5 seconds")
		recv_interval = 5
	}

	return config.ProcConfig{}, nil
}

// DatabricksProc Proc Generate all kinds of data.
func DatabricksProc(data any) []model.MeltModel {

	responseModel := data.(map[string]any)["model"]
	responseData := data.(map[string]any)["response"]

	out := make([]model.MeltModel, 0)

	switch responseModel {

	case "AWSDatabricksQueryList":

		var databricksJobsListModel = AWSQueriesList{}
		modelJson, _ := json.Marshal(responseData)
		err := json.Unmarshal(modelJson, &databricksJobsListModel)
		if err != nil {
			return nil
		}

		for i := 0; i < len(databricksJobsListModel.Res); i++ {
			e := reflect.ValueOf(&databricksJobsListModel.Res[i]).Elem()
			tags := make(map[string]interface{})
			tags["QueryId"] = databricksJobsListModel.Res[i].QueryId
			metricModels := utils.CreateMetricModels("AWSDatabricksQuery", e, tags)
			out = append(out, metricModels...)
		}

	case "AWSDatabricksJobsRunsList":
		var databricksJobsListModel = AWSJobRuns{}
		modelJson, _ := json.Marshal(responseData)
		err := json.Unmarshal(modelJson, &databricksJobsListModel)
		if err != nil {
			return nil
		}

		for i := 0; i < len(databricksJobsListModel.Runs); i++ {

			metricModel := createEventModels("AWS", databricksJobsListModel.Runs[i])
			out = append(out, metricModel)
		}

	case "AWSDatabricksPipelinesList":
		var databricksJobsListModel = AWSPipelinesList{}
		modelJson, _ := json.Marshal(responseData)
		err := json.Unmarshal(modelJson, &databricksJobsListModel)
		if err != nil {
			return nil
		}

		for i := 0; i < len(databricksJobsListModel.Statuses); i++ {

			metricModel := createEventModels("AWS", databricksJobsListModel.Statuses[i])
			out = append(out, metricModel)
		}

	case "GCPDatabricksQueryList":

		var databricksJobsListModel = GCPQueriesList{}
		modelJson, _ := json.Marshal(responseData)
		err := json.Unmarshal(modelJson, &databricksJobsListModel)
		if err != nil {
			return nil
		}

		for i := 0; i < len(databricksJobsListModel.Res); i++ {
			e := reflect.ValueOf(&databricksJobsListModel.Res[i]).Elem()
			tags := make(map[string]interface{})
			tags["QueryId"] = databricksJobsListModel.Res[i].QueryId
			metricModels := utils.CreateMetricModels("GCPDatabricksQuery", e, tags)
			out = append(out, metricModels...)
		}

	case "GCPDatabricksJobsRunsList":
		var databricksJobsListModel = GCPJobRuns{}
		modelJson, _ := json.Marshal(responseData)
		err := json.Unmarshal(modelJson, &databricksJobsListModel)
		if err != nil {
			return nil
		}

		for i := 0; i < len(databricksJobsListModel.Runs); i++ {

			metricModel := createEventModels("GCP", databricksJobsListModel.Runs[i])
			out = append(out, metricModel)
		}

	case "GCPDatabricksPipelinesList":
		var databricksJobsListModel = GCPPipelinesList{}
		modelJson, _ := json.Marshal(responseData)
		err := json.Unmarshal(modelJson, &databricksJobsListModel)
		if err != nil {
			return nil
		}

		for i := 0; i < len(databricksJobsListModel.Statuses); i++ {

			metricModel := createEventModels("GCP", databricksJobsListModel.Statuses[i])
			out = append(out, metricModel)
		}
	case "AzureDatabricksQueryList":

		var databricksJobsListModel = AzureQueriesList{}
		modelJson, _ := json.Marshal(responseData)
		err := json.Unmarshal(modelJson, &databricksJobsListModel)
		if err != nil {
			return nil
		}

		for i := 0; i < len(databricksJobsListModel.Res); i++ {
			e := reflect.ValueOf(&databricksJobsListModel.Res[i]).Elem()
			tags := make(map[string]interface{})
			tags["QueryId"] = databricksJobsListModel.Res[i].QueryId
			metricModels := utils.CreateMetricModels("AzureDatabricksQuery", e, tags)
			out = append(out, metricModels...)
		}

	case "AzureDatabricksJobsRunsList":
		var databricksJobsListModel = AzureJobRuns{}
		modelJson, _ := json.Marshal(responseData)
		err := json.Unmarshal(modelJson, &databricksJobsListModel)
		if err != nil {
			return nil
		}

		for i := 0; i < len(databricksJobsListModel.Runs); i++ {

			metricModel := createEventModels("Azure", databricksJobsListModel.Runs[i])
			out = append(out, metricModel)
		}

	case "AzureDatabricksPipelinesList":
		var databricksJobsListModel = AzurePipelinesList{}
		modelJson, _ := json.Marshal(responseData)
		err := json.Unmarshal(modelJson, &databricksJobsListModel)
		if err != nil {
			return nil
		}

		for i := 0; i < len(databricksJobsListModel.Statuses); i++ {

			metricModel := createEventModels("Azure", databricksJobsListModel.Statuses[i])
			out = append(out, metricModel)
		}
	default:
		log.Println("Unknown response model in Databricks integration")
	}

	return out
}

func createEventModels(cloud string, data any) model.MeltModel {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		log.Println(err.Error())
	}
	var jsonString map[string]interface{}

	err = json.Unmarshal(jsonBytes, &jsonString)
	if err != nil {
		log.Println(err.Error())
	}

	return model.MakeEvent(cloud+"DatabricksJobsRuns", jsonString, time.Now())
}
