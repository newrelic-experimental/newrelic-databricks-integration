package databricks

import (
	"context"

	databricksSdk "github.com/databricks/databricks-sdk-go"
	databricksSdkConfig "github.com/databricks/databricks-sdk-go/config"

	"github.com/newrelic/newrelic-labs-sdk/v2/pkg/integration"
	"github.com/newrelic/newrelic-labs-sdk/v2/pkg/integration/exporters"
	"github.com/newrelic/newrelic-labs-sdk/v2/pkg/integration/log"
	"github.com/newrelic/newrelic-labs-sdk/v2/pkg/integration/pipeline"
	"github.com/spf13/viper"
)

func InitPipelines(
	ctx context.Context,
	i *integration.LabsIntegration,
	tags map[string]string,
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
		// Create the newrelic exporter
		newRelicExporter := exporters.NewNewRelicExporter(
			"newrelic-api",
			i.Name,
			i.Id,
			i.NrClient,
			i.GetLicenseKey(),
			i.GetRegion(),
			i.DryRun,
		)

		// Create a metrics pipeline
		mp := pipeline.NewMetricsPipeline()
		mp.AddExporter(newRelicExporter)

		// Create the receiver
		databricksSparkReceiver := NewDatabricksSparkReceiver(w, tags)
		mp.AddReceiver(databricksSparkReceiver)

		log.Debugf("initializing Databricks Spark pipeline")

		i.AddPipeline(mp)
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
