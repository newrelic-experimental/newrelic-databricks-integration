package databricks

import (
	"context"

	databricksSdk "github.com/databricks/databricks-sdk-go"
	"github.com/newrelic/newrelic-labs-sdk/pkg/integration/model"
)

type DatabricksSdkReceiver struct {
	w			*databricksSdk.WorkspaceClient
}

func NewDatabricksSdkReceiver() (*DatabricksSdkReceiver, error) {
	w, err := databricksSdk.NewWorkspaceClient()
	if err != nil {
		return nil, err
	}

	return &DatabricksSdkReceiver{ w }, nil
}

func (d *DatabricksSdkReceiver) GetId() string {
	return "databricks-sdk-receiver"
}

func (d *DatabricksSdkReceiver) PollMetrics(
	context context.Context,
	writer chan <- model.Metric,
) error {
	return nil
}
