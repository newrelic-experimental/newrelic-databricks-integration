package spark

import (
	"context"
	"encoding/json"

	"github.com/newrelic/newrelic-labs-sdk/pkg/integration/connectors"
	"github.com/newrelic/newrelic-labs-sdk/pkg/integration/log"
)

type SparkApiClient interface {
	GetApplications(ctx context.Context) ([]SparkApplication, error)
	GetApplicationExecutors(
		ctx context.Context,
		app *SparkApplication,
	) ([]SparkExecutor, error)
	GetApplicationJobs(
		ctx context.Context,
		app *SparkApplication,
	) ([]SparkJob, error)
	GetApplicationStages(
		ctx context.Context,
		app *SparkApplication,
	) ([]SparkStage, error)
	GetApplicationRDDs(
		ctx context.Context,
		app *SparkApplication,
	) ([]SparkRDD, error)
}

type NativeSparkApiClient struct {
	sparkContextUiUrl		string
	authenticator			connectors.HttpAuthenticator
}

func NewNativeSparkApiClient(
	sparkContextUiUrl string,
	authenticator connectors.HttpAuthenticator,
) *NativeSparkApiClient {
	return &NativeSparkApiClient{
		sparkContextUiUrl,
		authenticator,
	}
}

func (s *NativeSparkApiClient) GetApplications(
	ctx context.Context,
) ([]SparkApplication, error) {
	sparkApps := []SparkApplication{}

	err := makeRequest(
		s.sparkContextUiUrl + "/api/v1/applications",
		s.authenticator,
		&sparkApps,
	)
	if err !=  nil {
		return nil, err
	}

	return sparkApps, nil
}

func (s *NativeSparkApiClient) GetApplicationExecutors(
	ctx context.Context,
	app *SparkApplication,
) ([]SparkExecutor, error) {
	executors := []SparkExecutor{}

	err := makeRequest(
		s.sparkContextUiUrl + "/api/v1/applications/" + app.Id + "/executors",
		s.authenticator,
		&executors,
	)
	if err !=  nil {
		return nil, err
	}

	return executors, nil
}

func (s *NativeSparkApiClient) GetApplicationJobs(
	ctx context.Context,
	app *SparkApplication,
) ([]SparkJob, error) {
	jobs := []SparkJob{}

	err := makeRequest(
		s.sparkContextUiUrl + "/api/v1/applications/" + app.Id + "/jobs",
		s.authenticator,
		&jobs,
	)
	if err !=  nil {
		return nil, err
	}

	return jobs, nil
}

func (s *NativeSparkApiClient) GetApplicationStages(
	ctx context.Context,
	app *SparkApplication,
) ([]SparkStage, error) {
	stages := []SparkStage{}

	err := makeRequest(
		s.sparkContextUiUrl + "/api/v1/applications/" + app.Id + "/stages?details=true",
		s.authenticator,
		&stages,
	)
	if err != nil {
		return nil, err
	}

	return stages, nil
}

func (s *NativeSparkApiClient) GetApplicationRDDs(
	ctx context.Context,
	app *SparkApplication,
) ([]SparkRDD, error) {
	rdds := []SparkRDD{}

	err := makeRequest(
		s.sparkContextUiUrl + "/api/v1/applications/" + app.Id + "/storage/rdd",
		s.authenticator,
		&rdds,
	)
	if err !=  nil {
		return nil, err
	}

	return rdds, nil
}

func makeRequest(
	url string,
	authenticator connectors.HttpAuthenticator,
	response interface{},
) error {
	connector := connectors.NewHttpGetConnector(url)

	if authenticator != nil {
		connector.SetAuthenticator(authenticator)
	}

	connector.SetHeaders(map[string]string {
		"Content-Type": "application/json",
		"Accept": "application/json",
	})

	in, err := connector.Request()
	if err != nil {
		return err
	}

	log.Debugf("decoding spark JSON response for URL %s", url)

	dec := json.NewDecoder(in)

	err = dec.Decode(response)
	if err != nil {
		return err
	}

	if log.IsDebugEnabled() {
		log.PrettyPrintJson(response)
	}

	return nil
}
