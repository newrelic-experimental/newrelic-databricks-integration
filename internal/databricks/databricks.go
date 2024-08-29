package databricks

import (
	"context"
	"fmt"
	"strings"

	databricksSdk "github.com/databricks/databricks-sdk-go"
	databricksSdkConfig "github.com/databricks/databricks-sdk-go/config"
	databricksSdkCompute "github.com/databricks/databricks-sdk-go/service/compute"
	"github.com/newrelic-experimental/newrelic-databricks-integration/internal/spark"

	"github.com/newrelic/newrelic-labs-sdk/v2/pkg/integration"
	"github.com/newrelic/newrelic-labs-sdk/v2/pkg/integration/log"
	"github.com/spf13/viper"
)

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

	err := configureAuth(databricksConfig)
	if err != nil {
		return err
	}

	/*
	 * If the user explicitly specifies a host in the config, use that.
	 * Otherwise the user can specify using an SDK-supported mechanism.
	 */
	databricksWorkspaceHost := viper.GetString("databricks.workspaceHost")
	if databricksWorkspaceHost != "" {
		databricksConfig.Host = databricksWorkspaceHost
	}

	w, err := databricksSdk.NewWorkspaceClient(databricksConfig)
	if err != nil {
		return err
	}

	if sparkMetrics {
		all, err := w.Clusters.ListAll(
			ctx,
			databricksSdkCompute.ListClustersRequest{},
		)
		if err != nil {
		  return fmt.Errorf("failed to list clusters: %w", err)
		}

		for _, c := range all {
			if c.State != databricksSdkCompute.StateRunning {
				log.Debugf(
					"skipping cluster %s because it is not running",
					c.ClusterName,
				)
				continue
			}

			if c.ClusterSource == databricksSdkCompute.ClusterSourceUi ||
				c.ClusterSource == databricksSdkCompute.ClusterSourceApi {

				/// resolve the spark context UI URL for the cluster
				log.Debugf(
					"resolving Spark context UI URL for cluster %s",
					c.ClusterName,
				)
				sparkContextUiPath, err := getSparkContextUiPathForCluster(
					ctx,
					w,
					&c,
				)
				if err != nil {
					return err
				}

				databricksSparkApiClient, err := NewDatabricksSparkApiClient(
					sparkContextUiPath,
					w,
				)
				if err != nil {
					return err
				}

				// Initialize spark pipelines
				log.Debugf(
					"initializing Spark pipeline for cluster %s with spark context UI URL %s",
					c.ClusterName,
					sparkContextUiPath,
				)
				err = spark.InitPipelines(
					i,
					databricksSparkApiClient,
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

	// @TODO: initialize databricks pipelines here

	return nil
}

func configureAuth(config *databricksSdk.Config) error {
	/*
	 * Any of the variables below can be specified in any of the ways that
	 * are supported by the Databricks SDK so if we don't explicitly find one
	 * in the config file, it's not an error.  We assume the user has used one
	 * of the SDK mechanisms and if they haven't the SDK will return an error at
	 * config time or when a request fails.
	 */

	// Prefer OAuth by looking for client ID in our config first
	databricksOAuthClientId := viper.GetString("databricks.oauthClientId")
	if databricksOAuthClientId != "" {
		/*
		 * If an OAuth client ID was in our config we will at this point tell
		 * the SDK to use OAuth M2M authentication. The secret may come from our
		 * config but can still come from any of the supported SDK mechanisms.
		 * So if we don't find the secret in our config file, it's not an error.
		 * Note that because we are forcing OAuth M2M authentication now, the
		 * SDK will not try other mechanisms if OAuth M2M authentication is
		 * unsuccessful.
		 */
		config.ClientID = databricksOAuthClientId
		config.Credentials = databricksSdkConfig.M2mCredentials{}

		databricksOAuthClientSecret := viper.GetString(
			"databricks.oauthClientSecret",
		)
		if databricksOAuthClientSecret != "" {
			config.ClientSecret = databricksOAuthClientSecret
		}

		return nil
	}

	// Check for a PAT in our config next
	databricksAccessToken := viper.GetString("databricks.accessToken")
	if databricksAccessToken != "" {
		/*
		* If the user didn't specify an OAuth client ID but does specify a PAT,
		* we will at this point tell the SDK to use PAT authentication. Note
		* that because we are forcing PAT authentication now, the SDK will not
		* try other mechanisms if PAT authentication is unsuccessful.
		*/
		config.Token = databricksAccessToken
		config.Credentials = databricksSdkConfig.PatCredentials{}

		return nil
	}

	/*
	 * At this point, it's up to the user to specify authentication via an
	 * SDK-supported mechanism. This does not preclude the user from using OAuth
	 * M2M authentication or PAT authentication. The user can still use these
	 * authentication types via SDK-supported mechanisms or any other
	 * SDK-supported authentication types via the corresponding SDK-supported
	 * mechanisms.
	 */

	return nil
}

func getSparkContextUiPathForCluster(
	ctx context.Context,
	w *databricksSdk.WorkspaceClient,
	c *databricksSdkCompute.ClusterDetails,
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
		"/driver-proxy-api/o/%s/%s/%s",
		vals[0],
		clusterId,
		vals[1],
	)

	return url, nil
}
