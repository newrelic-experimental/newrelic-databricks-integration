package databricks

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	databricksSdk "github.com/databricks/databricks-sdk-go"
	databricksSdkConfig "github.com/databricks/databricks-sdk-go/config"
	databricksSdkCompute "github.com/databricks/databricks-sdk-go/service/compute"
	"github.com/newrelic-experimental/newrelic-databricks-integration/internal/spark"

	"github.com/newrelic/newrelic-labs-sdk/pkg/integration"
	"github.com/newrelic/newrelic-labs-sdk/pkg/integration/connectors"
	"github.com/newrelic/newrelic-labs-sdk/pkg/integration/log"
	"github.com/spf13/viper"
)

type DatabricksAuthenticator struct {
	accessToken 		string
}

func NewDatabricksAuthenticator(accessToken string) (
	*DatabricksAuthenticator,
) {
	return &DatabricksAuthenticator{ accessToken }
}

func (b *DatabricksAuthenticator) Authenticate(
	connector *connectors.HttpConnector,
	req *http.Request,
) error {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", b.accessToken))

	return nil
}

func InitPipelines(
	ctx context.Context,
	i *integration.LabsIntegration,
) error {
	// Databricks config
	databricksConfig := &databricksSdk.Config{}

	// Should we collect Spark metrics?
	sparkMetrics := true
	if viper.IsSet("databricks.sparkMetrics") {
		sparkMetrics = viper.GetBool("databricks.sparkMetrics")
	}

	// Databricks credentials
	//
	// If we are not collecting Spark metrics, we support PAT authentication
	// using our config file or environment variables, or any supported SDK
	// authentication method supported by the SDK itself, e.g. via the
	// environment or .databrickscfg or cloud specific methods.
	//
	// If we are collecting Spark metrics, we only support PAT authentication
	// via our config file or environment variables for host and token because
	// I'm not sure how we would use other forms to get through the driver to
	// the Spark UI. HTTP Basic might be a possibility but that is a @TODO.

	databricksAccessToken := viper.GetString("databricks.accessToken")
	if databricksAccessToken != "" {
		databricksConfig.Token = databricksAccessToken
		databricksConfig.Credentials = databricksSdkConfig.PatCredentials{}
	} else if sparkMetrics {
		return fmt.Errorf("missing databricks personal access token")
	}

	// Workspace Host
	databricksWorkspaceHost := viper.GetString("databricks.workspaceHost")
	if databricksWorkspaceHost != "" {
		databricksConfig.Host = databricksWorkspaceHost
	} else if sparkMetrics {
		return fmt.Errorf("missing databricks workspace host")
	}

	w, err := databricksSdk.NewWorkspaceClient(databricksConfig)
	if err != nil {
		return err
	}

	if sparkMetrics {
		// Initialize the Spark pipelines
		authenticator := NewDatabricksAuthenticator(databricksAccessToken)

		all, err := w.Clusters.ListAll(
			ctx,
			databricksSdkCompute.ListClustersRequest{},
		)
		if err != nil {
		  return fmt.Errorf("failed to list clusters: %w", err)
		}

		for _, c := range all {
			if c.State == databricksSdkCompute.StateRunning {
				if c.ClusterSource == databricksSdkCompute.ClusterSourceUi ||
					c.ClusterSource == databricksSdkCompute.ClusterSourceApi {

					/// resolve the spark context UI URL for the cluster
					log.Debugf(
						"resolving Spark context UI URL for cluster %s",
						c.ClusterName,
					)
					sparkContextUiUrl, err := getSparkContextUiUrlForCluster(
						ctx,
						w,
						&c,
						databricksWorkspaceHost,
					)
					if err != nil {
						return err
					}

					// Initialize spark pipelines
					log.Debugf(
						"initializing Spark pipeline for cluster %s with spark context UI URL %s",
						c.ClusterName,
						sparkContextUiUrl,
					)
					err = spark.InitPipelinesForContext(
						i,
						sparkContextUiUrl,
						authenticator,
						map[string] string {
							"clusterProvider": "databricks",
							"databricksClusterId": c.ClusterId,
							"databricksClusterName": c.ClusterName,
						},
					)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	// @TODO: initialize databricks pipelines here

	return nil
}

func getSparkContextUiUrlForCluster(
	ctx context.Context,
	w *databricksSdk.WorkspaceClient,
	c *databricksSdkCompute.ClusterDetails,
	databricksWorkspaceHost string,
) (string, error) {
	// @see https://databrickslabs.github.io/overwatch/assets/_index/realtime_helpers.html

	clusterId := c.ClusterId

	waitContextStatus, err := w.CommandExecution.Create(
		ctx,
		databricksSdkCompute.CreateContext{
			ClusterId: clusterId,
			Language: databricksSdkCompute.LanguagePython,
		},
	)
	if err != nil {
		return "", err
	}

	execContext, err := waitContextStatus.Get()
	if err != nil {
		return "", err
	}

	cmd := databricksSdkCompute.Command{
		ClusterId: clusterId,
		Command: `
		print(f'{spark.conf.get("spark.databricks.clusterUsageTags.clusterOwnerOrgId")}')
		print(f'{spark.conf.get("spark.ui.port")}')
		`,
		ContextId: execContext.Id,
		Language: databricksSdkCompute.LanguagePython,
	}

	waitCommandStatus, err := w.CommandExecution.Execute(
		ctx,
		cmd,
	)
	if err != nil {
		return "", err
	}

	resp, err := waitCommandStatus.Get()
	if err != nil {
		return "", err
	}

	data, ok := resp.Results.Data.(string);
	if !ok {
		return "", fmt.Errorf("command result is not a string value")
	}

	vals := strings.Split(data, "\n")
	if len(vals) != 2 {
		return "", fmt.Errorf("invalid command result")
	}

	if vals[0] == "" || vals[1] == "" {
		return "", fmt.Errorf("empty command results")
	}

	// @TODO: I think this URL pattern only works for multi-tenant accounts.
	// We may need a flag for single tenant accounts and use the o/0 form
	// shown on the overwatch site.

	url := fmt.Sprintf(
		"https://%s/driver-proxy-api/o/%s/%s/%s",
		databricksWorkspaceHost,
		vals[0],
		clusterId,
		vals[1],
	)

	return url, nil
}
