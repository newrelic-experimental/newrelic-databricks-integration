package databricks

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"newrelic/multienv/integration/Utils"
	"newrelic/multienv/pkg/config"
	"newrelic/multienv/pkg/connect"
	"newrelic/multienv/pkg/model"
	"reflect"
	"time"
)

func GetDatabricksConnectors(pipeConfig *config.PipelineConfig) ([]connect.Connector, error) {

	licenseKey := ""
	databricksEndpoint := ""

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

	var headers = make(map[string]string)
	headers["Authorization"] = "Bearer " + licenseKey
	headers["Content-Type"] = "application/json"
	headers["Accept"] = "application/json"

	var connectors []connect.Connector

	queryHistoryConnector := connect.MakeHttpGetConnector(databricksEndpoint+"/api/2.0/sql/history/queries?include_metrics=true", headers)
	queryHistoryConnector.SetConnectorModelName("DatabricksQueryList")

	jobRunsListConnector := connect.MakeHttpGetConnector(databricksEndpoint+"/api/2.1/jobs/runs/list", headers)
	jobRunsListConnector.SetConnectorModelName("DatabricksJobsRunsList")

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

// Proc Generate all kinds of data.
func DatabricksProc(data any) []model.MeltModel {

	responseModel := data.(map[string]any)["model"]
	responseData := data.(map[string]any)["response"]

	out := make([]model.MeltModel, 0)
	if responseModel == "DatabricksQueryList" {

		var databricksJobsListModel = AWSQueriesList{}
		modelJson, _ := json.Marshal(responseData)
		err := json.Unmarshal(modelJson, &databricksJobsListModel)
		if err != nil {
			return nil
		}

		for i := 0; i < len(databricksJobsListModel.Res); i++ {

			jsonBytes, err := json.Marshal(databricksJobsListModel.Res[i])
			if err != nil {
				log.Println(err.Error())
			}
			var jsonString map[string]interface{}

			err = json.Unmarshal(jsonBytes, &jsonString)
			if err != nil {
				log.Println(err.Error())
			}
			result := make(map[string]interface{})
			result = Utils.Flatten(jsonString, "", result)

			e := reflect.ValueOf(&databricksJobsListModel.Res[i]).Elem()
			tags := make(map[string]interface{})
			stagesTags := make(map[string]interface{})
			Utils.SetTags("databricks.query.", e, tags, stagesTags)

			metricModels := Utils.CreateMetricModels("databricks.query.", e, stagesTags)

			out = append(out, metricModels...)
		}

	} else if responseModel == "DatabricksJobsRunsList" {

		var databricksJobsListModel = AWSJobRuns{}
		modelJson, _ := json.Marshal(responseData)
		err := json.Unmarshal(modelJson, &databricksJobsListModel)
		if err != nil {
			return nil
		}

		for i := 0; i < len(databricksJobsListModel.Runs); i++ {

			jsonBytes, err := json.Marshal(databricksJobsListModel.Runs[i])
			if err != nil {
				log.Println(err.Error())
			}
			var jsonString map[string]interface{}

			err = json.Unmarshal(jsonBytes, &jsonString)
			if err != nil {
				log.Println(err.Error())
			}

			metricModel := model.MakeEvent("DatabricksJobsRuns", jsonString, time.Now())

			out = append(out, metricModel)
		}

	} else {
		log.Println("Unknown response model in Databricks integration")
	}

	return out
}
