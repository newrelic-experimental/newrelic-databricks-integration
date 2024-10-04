<a href="https://opensource.newrelic.com/oss-category/#new-relic-experimental"><picture><source media="(prefers-color-scheme: dark)" srcset="https://github.com/newrelic/opensource-website/raw/main/src/images/categories/dark/Experimental.png"><source media="(prefers-color-scheme: light)" srcset="https://github.com/newrelic/opensource-website/raw/main/src/images/categories/Experimental.png"><img alt="New Relic Open Source experimental project banner." src="https://github.com/newrelic/opensource-website/raw/main/src/images/categories/Experimental.png"></picture></a>

![GitHub forks](https://img.shields.io/github/forks/newrelic-experimental/newrelic-databricks-integration?style=social)
![GitHub stars](https://img.shields.io/github/stars/newrelic-experimental/newrelic-databricks-integration?style=social)
![GitHub watchers](https://img.shields.io/github/watchers/newrelic-experimental/newrelic-databricks-integration?style=social)

![GitHub all releases](https://img.shields.io/github/downloads/newrelic-experimental/newrelic-databricks-integration/total)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/newrelic-experimental/newrelic-databricks-integration)
![GitHub last commit](https://img.shields.io/github/last-commit/newrelic-experimental/newrelic-databricks-integration)
![GitHub Release Date](https://img.shields.io/github/release-date/newrelic-experimental/newrelic-databricks-integration)

![GitHub issues](https://img.shields.io/github/issues/newrelic-experimental/newrelic-databricks-integration)
![GitHub issues closed](https://img.shields.io/github/issues-closed/newrelic-experimental/newrelic-databricks-integration)
![GitHub pull requests](https://img.shields.io/github/issues-pr/newrelic-experimental/newrelic-databricks-integration)
![GitHub pull requests closed](https://img.shields.io/github/issues-pr-closed/newrelic-experimental/newrelic-databricks-integration)

# New Relic Databricks Integration

This integration collects telemetry from Databricks (including Spark on
Databricks) and/or Spark telemetry from any Spark deployment. See the
[Features](#features) section for supported telemetry types.

![Apache Spark Dashboard Screenshot](./examples/spark-dashboard-jobs.png)

## Table of Contents

* [Important Notes](#important-notes)
* [Getting Started](#getting-started)
   * [On-host Deployment](#on-host)
* [Features](#features)
* [Usage](#usage)
   * [Command Line Options](#command-line-options)
   * [Configuration](#configuration)
      * [General configuration](#general-configuration)
      * [Pipeline configuration](#pipeline-configuration)
      * [Log configuration](#log-configuration)
      * [Databricks configuration](#databricks-configuration)
      * [Spark configuration](#spark-configuration)
   * [Authentication](#authentication)
   * [Billable Usage & List Pricing Data](#billable-usage--list-pricing-data)
* [Building](#building)
   * [Coding Conventions](#coding-conventions)
   * [Local Development](#local-development)
   * [Releases](#releases)
   * [Github Workflows](#github-workflows)
* [Appendix](#appendix)
   * [Monitoring Cluster Health](#monitoring-cluster-health)

## Important Notes

* All references within this document to Databricks documentation reference the
  [Databricks on AWS documentation](https://docs.databricks.com/en/index.html).
  Use the cloud switcher menu located in the upper right hand corner of the
  documentation to select corresponding documentation for a different cloud.
* [On-host deployment](#on-host) is currently the only supported deployment
  type. For Databricks and non-Databricks Spark deployments, the integration can
  be deployed on [any supported host platform](#deploy-the-integration-on-a-host).
  For Databricks, support is also provided to deploy the integration
  on the [driver node of a Databricks cluster](#deploy-the-integration-on-the-driver-node-of-a-databricks-cluster)
  using a [cluster-scoped init script](https://docs.databricks.com/en/init-scripts/cluster-scoped.html).

## Getting Started

To get started with the New Relic Databricks integration,
[deploy the integration](#on-host) using a supported deployment type,
[configure](#configuration) the integration using supported configuration
mechanisms, and then [import](https://docs.newrelic.com/docs/query-your-data/explore-query-data/dashboards/dashboards-charts-import-export-data/#import-json)
the [sample dashboard](./examples/spark-dashboard.json).

### On-host

The New Relic Databricks integration can be run on any supported host platform.
The integration will collect Databricks telemetry (including Spark on
Databricks) via the [Databricks ReST API](https://docs.databricks.com/api/workspace/introduction)
using the [Databricks SDK for Go](https://docs.databricks.com/en/dev-tools/sdk-go.html)
and/or Spark telemetry from a non-Databricks Spark deployment via the
[Spark ReST API](https://spark.apache.org/docs/3.5.2/monitoring.html#rest-api).

The New Relic Databricks integration can also be deployed on the driver node of
a Databricks [cluster](https://docs.databricks.com/en/getting-started/concepts.html#cluster)
using the provided [init script](./init/cluster_init_integration.sh) to install
and configure the integration at cluster startup time.

#### Deploy the integration on a host

The New Relic Databricks integration provides binaries for the following
host platforms.

* Linux amd64
* Windows amd64

To run the Databricks integration on a host, perform the following steps.

1. Download the appropriate archive for your platform from the [latest release](https://github.com/newrelic-experimental/newrelic-databricks-integration/releases).
1. Extract the archive to a new or existing directory.
1. Create a directory named `configs` in the same directory.
1. Create a file named `config.yml` in the `configs` directory and copy the
   contents of the file [`configs/config.template.yml`](./configs/config.template.yml)
   in this repository into it.
1. Edit the `config.yml` file to [configure](#configuration) the integration
   appropriately for your environment.
1. From the directory where the archive was extracted, execute the integration
   binary using the command `./newrelic-databricks-integration` (or
   `.\newrelic-databricks-integration.exe` on Windows) with the appropriate
   [Command Line Options](#command-line-options).

#### Deploy the integration on the driver node of a Databricks cluster

The New Relic Databricks integration can be deployed on the driver node of a
Databricks [cluster](https://docs.databricks.com/en/getting-started/concepts.html#cluster)
using a [cluster-scoped init script](https://docs.databricks.com/en/init-scripts/cluster-scoped.html).
The [init script](./init/cluster_init_integration.sh) uses custom
[environment variables](https://docs.databricks.com/en/compute/configure.html#env-var)
to specify configuration parameters necessary for the integration [configuration](#configuration).

To install the [init script](./init/cluster_init_integration.sh), perform the
following steps.

1. Login to your Databricks account and navigate to the desired
   [workspace](https://docs.databricks.com/en/getting-started/concepts.html#accounts-and-workspaces).
1. Follow the [recommendations for init scripts](https://docs.databricks.com/en/init-scripts/index.html#recommendations-for-init-scripts)
   to store the [`cluster_init_integration.sh`](./init/cluster_init_integration.sh)
   script within your workspace in the recommended manner. For example, if your
   workspace is [enabled for Unity Catalog](https://docs.databricks.com/en/data-governance/unity-catalog/get-started.html#step-1-confirm-that-your-workspace-is-enabled-for-unity-catalog),
   you should store the init script in a [Unity Catalog volume](https://docs.databricks.com/en/ingestion/file-upload/upload-to-volume.html).
1. Navigate to the [`Compute`](https://docs.databricks.com/en/compute/clusters-manage.html#view-compute)
   tab and select the desired all-purpose or job compute to open the compute
   details UI.
1. Click the button labeled `Edit` to [edit the compute's configuration](https://docs.databricks.com/en/compute/clusters-manage.html#edit-a-compute).
1. Follow the steps to [use the UI to configure a cluster-scoped init script](https://docs.databricks.com/en/init-scripts/cluster-scoped.html#configure-a-cluster-scoped-init-script-using-the-ui)
   and point to the location where you stored the init script in step 2 above.
1. If your cluster is not running, click on the button labeled `Confirm` to
   save your changes. Then, restart the cluster. If your cluster is already
   running, click on the button labeled `Confirm and restart` to save your
   changes and restart the cluster.

Additionally, follow the steps to [set environment variables](https://docs.databricks.com/en/compute/configure.html#environment-variables)
to add the following environment variables.

* `NEW_RELIC_API_KEY` - Your [New Relic User API Key](https://docs.newrelic.com/docs/apis/intro-apis/new-relic-api-keys/#user-key)
* `NEW_RELIC_LICENSE_KEY` - Your [New Relic License Key](https://docs.newrelic.com/docs/apis/intro-apis/new-relic-api-keys/#license-key)
* `NEW_RELIC_ACCOUNT_ID` - Your [New Relic Account ID](https://docs.newrelic.com/docs/accounts/accounts-billing/account-structure/account-id/)
* `NEW_RELIC_REGION` - The [region](https://docs.newrelic.com/docs/accounts/accounts-billing/account-setup/choose-your-data-center/#regions-availability)
   of your New Relic account; one of `US` or `EU`
* `NEW_RELIC_DATABRICKS_WORKSPACE_HOST` - The [instance name](https://docs.databricks.com/en/workspace/workspace-details.html#workspace-instance-names-urls-and-ids)
   of the target Databricks instance
* `NEW_RELIC_DATABRICKS_ACCESS_TOKEN` - To [authenticate](#authentication) with
   a [personal access token](https://docs.databricks.com/en/dev-tools/auth/pat.html#databricks-personal-access-tokens-for-workspace-users),
   your personal access token
* `NEW_RELIC_DATABRICKS_OAUTH_CLIENT_ID` - To [use a service principal to authenticate with Databricks (OAuth M2M)](https://docs.databricks.com/en/dev-tools/auth/oauth-m2m.html),
   the OAuth client ID for the service principal
* `NEW_RELIC_DATABRICKS_OAUTH_CLIENT_SECRET` - To [use a service principal to authenticate with Databricks (OAuth M2M)](https://docs.databricks.com/en/dev-tools/auth/oauth-m2m.html),
   an OAuth client secret associated with the service principal

Note that the `NEW_RELIC_API_KEY` and `NEW_RELIC_ACCOUNT_ID` are currently
unused but are required by the [new-relic-client-go](https://github.com/newrelic/newrelic-client-go)
module used by the integration. Additionally, note that only the personal access
token _or_ OAuth credentials need to be specified but not both. If both are
specified, the OAuth credentials take precedence. Finally, make sure to restart
the cluster following the configuration of the environment variables.

## Features

The New Relic Databricks integration supports the following capabilities.

* Collect Spark telemetry

  The New Relic Databricks integration can collect telemetry from Spark running
  on Databricks. By default, the integration will automatically connect to
  and collect telemetry from the Spark deployments in all clusters created via
  the UI or API in the specified workspace.

  The New Relic Databricks integration can also collect Spark telemetry from any
  non-Databricks Spark deployment.

* Collect Databricks billable usage and list pricing data

  The New Relic Databricks integration can collect Databricks
  [billable usage](https://docs.databricks.com/en/admin/system-tables/billing.html)
  and [list pricing](https://docs.databricks.com/en/admin/system-tables/pricing.html)
  data from the Databricks [system tables](https://docs.databricks.com/en/admin/system-tables/index.html).
  This data can be used to show basic Databricks DBU consumption and cost
  metrics directly within New Relic.

## Usage

### Command Line Options

| Option | Description | Default |
| --- | --- | --- |
| --config_path | path to the (#configyml) to use | `configs/config.yml` |
| --dry_run | flag to enable "dry run" mode | `false` |
| --env_prefix | prefix to use for environment variable lookup | `''` |
| --verbose | flag to enable "verbose" mode | `false` |
| --version | display version information only | N/a |

### Configuration

The Databricks integration is configured using the [`config.yml`](#configyml)
and/or environment variables. For Databricks, authentication related configuration
parameters may also be set in a [Databricks configuration profile](https://docs.databricks.com/en/dev-tools/auth/config-profiles.html).
In all cases, where applicable, environment variables always take precedence.

#### `config.yml`

All configuration parameters for the Databricks integration can be set using a
YAML file named [`config.yml`](#configyml). The default location for this file
is `configs/config.yml` relative to the current working directory when the
integration binary is executed. The supported configuration parameters are
listed below. See [`config.template.yml`](./configs/config.template.yml)
for a full configuration example.

##### General configuration

The parameters in this section are configured at the top level of the
[`config.yml`](#configyml).

###### `licenseKey`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| New Relic license key | string | Y | N/a |

This parameter specifies the New Relic License Key (INGEST) that should be used
to send generated metrics.

The license key can also be specified using the `NEW_RELIC_LICENSE_KEY`
environment variable.

###### `region`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| New Relic region identifier | `US` / `EU` | N | `US` |

This parameter specifies which New Relic region that generated metrics should be
sent to.

###### `interval`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| Polling interval (in _seconds_) | numeric | N | 60 |

This parameter specifies the interval (in _seconds_) at which the integration
should poll for data.

This parameter is only used when [`runAsService`](#runasservice) is set to
`true`.

###### `runAsService`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| Flag to enable running the integration as a "service" | `true` / `false` | N | `false` |

The integration can run either as a "service" or as a simple command line
utility which runs once and exits when it is complete.

When set to `true`, the integration process will run continuously and poll the
for data at the recurring interval specified by the [`interval`](#interval)
parameter. The process will only exit if it is explicitly stopped or a fatal
error or panic occurs.

When set to `false`, the integration will run once and exit. This is intended for
use with an external scheduling mechanism like [cron](https://man7.org/linux/man-pages/man8/cron.8.html).

###### `pipeline`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| The root node for the set of [pipeline configuration](#pipeline-configuration) parameters | YAML Mapping | N | N/a |

The integration retrieves, processes, and exports data to New Relic using
a data pipeline consisting of one or more receivers, a processing chain, and a
New Relic exporter. Various aspects of the pipeline are configurable. This
element groups together the configuration parameters related to
[pipeline configuration](#pipeline-configuration).

###### `log`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| The root node for the set of [log configuration](#log-configuration) parameters | YAML Mapping | N | N/a |

The integration uses the [logrus](https://pkg.go.dev/github.com/sirupsen/logrus)
package for application logging. This element groups together the configuration
parameters related to [log configuration](#log-configuration).

###### `mode`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| The integration execution mode | `databricks` | N | `databricks` |

The integration execution mode. Currently, the only supported execution mode is
`databricks`.

**Deprecated:** As of v2.3.0, this configuration parameter is no longer used.
The presence (or not) of the [`databricks`](#databricks) top-level node will be
used to enable (or disable) the Databricks collector. Likewise, the presence
(or not) of the [`spark`](#spark) top-level node will be used to enable (or
disable) the Spark collector separate from Databricks.

###### `databricks`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| The root node for the set of [Databricks configuration](#databricks-configuration) parameters | YAML Mapping | N | N/a |

This element groups together the configuration parameters to [configure](#databricks-configuration)
the Databricks collector. If this element is not specified, the Databricks
collector will not be run.

Note that this node is not required. It can be used with or without the
[`spark`](#spark) top-level node.

###### `spark`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| The root node for the set of [Spark configuration](#spark-configuration) parameters | YAML Mapping | N | N/a |

This element groups together the configuration parameters to [configure](#spark-configuration)
the Spark collector. If this element is not specified, the Spark collector will
not be run.

Note that this node is not required. It can be used with or without the
[`databricks`](#databricks) top-level node.

###### `tags`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| The root node for a set of custom [tags](https://docs.newrelic.com/docs/new-relic-solutions/new-relic-one/core-concepts/use-tags-help-organize-find-your-data/) to add to all telemetry sent to New Relic | YAML Mapping | N | N/a |

This element specifies a group of custom [tags](https://docs.newrelic.com/docs/new-relic-solutions/new-relic-one/core-concepts/use-tags-help-organize-find-your-data/)
that will be added to all telemetry sent to New Relic. The tags are specified as
a set of key-value pairs.

##### Pipeline configuration

###### `receiveBufferSize`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| Size of the buffer that holds items before processing | number | N | 500 |

This parameter specifies the size of the buffer that holds received items before
being flushed through the processing chain and on to the exporters. When this
size is reached, the items in the buffer will be flushed automatically.

###### `harvestInterval`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| Harvest interval (in _seconds_) | number | N | 60 |

This parameter specifies the interval (in _seconds_) at which the pipeline
should automatically flush received items through the processing chain and on
to the exporters. Each time this interval is reached, the pipeline will flush
items even if the item buffer has not reached the size specified by the
[`receiveBufferSize`](#receiveBufferSize) parameter.

###### `instances`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| Number of concurrent pipeline instances to run | number | N | 3 |

The integration retrieves, processes, and exports metrics to New Relic using
a data pipeline consisting of one or more receivers, a processing chain, and a
New Relic exporter. When [`runAsService`](#runasservice) is `true`, the
integration can launch one or more "instances" of this pipeline to receive,
process, and export data concurrently. Each "instance" will be configured with
the same processing chain and exporters and the receivers will be spread across
the available instances in a round-robin fashion.

This parameter specifies the number of pipeline instances to launch.

**NOTE:** When [`runAsService`](#runasservice) is `false`, only a single
pipeline instance is used.

##### Log configuration

###### `level`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| Log level | `panic` / `fatal` / `error` / `warn` / `info` / `debug` / `trace`  | N | `warn` |

This parameter specifies the maximum severity of log messages to output with
`trace` being the least severe and `panic` being the most severe. For example,
at the default log level (`warn`), all log messages with severities `warn`,
`error`, `fatal`, and `panic` will be output but `info`, `debug`, and `trace`
will not.

###### `fileName`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| Path to a file where log output will be written | string | N | `stderr` |

This parameter designates a file path where log output should be written. When
no path is specified, log output will be written to the standard error stream
(`stderr`).

##### Databricks configuration

The Databricks configuration parameters are used to configure the Databricks
collector.

###### `workspaceHost`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| Databricks workspace instance name | string | conditional | N/a |

This parameter specifies the [instance name](https://docs.databricks.com/en/workspace/workspace-details.html#workspace-instance-names-urls-and-ids)
of the target Databricks instance for which data should be collected. This is
used by the integration when constructing the URLs for API calls. Note that the
value of this parameter _must not_ include the `https://` prefix, e.g.
`https://my-databricks-instance-name.cloud.databricks.com`.

This parameter is required when the collection of Spark telemetry for Spark
running on Databricks is [enabled](#databricks-spark-enabled). Note that this
does not apply when the integration is [deployed directly on the driver node](#deploy-the-integration-on-the-driver-node-of-a-databricks-cluster)
via the provided [init script](./init/cluster_init_integration.sh). This
parameter is unused in that scenario.

The workspace host can also be specified using the `DATABRICKS_HOST`
environment variable.

**NOTE:** The `DATABRICKS_HOST` environment variable can not be used to specify
_both_ the [instance name](https://docs.databricks.com/en/workspace/workspace-details.html#workspace-instance-names-urls-and-ids)
and the accounts API endpoint. To account for this, the environment variables
`DATABRICKS_WORKSPACEHOST` and `DATABRICKS_ACCOUNTHOST` environment variables
can be alternately used either separately or in combination with the
`DATABRICKS_HOST` environment variable to specify the
[instance name](https://docs.databricks.com/en/workspace/workspace-details.html#workspace-instance-names-urls-and-ids)
and the accounts API endpoint, respectively.

###### `accountHost`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| Databricks accounts API endpoint | string | conditional | N/a |

This parameter specifies the accounts API endpoint. This is
used by the integration when constructing the URLs for account-level
[ReST API](https://docs.databricks.com/api/workspace/introduction) calls. Note
that unlike the value of [`workspaceHost`](#workspacehost), the value of this
parameter _must_ include the `https://` prefix, e.g.
`https://accounts.cloud.databricks.com`.

This parameter is required when the collection of Databricks billable usage and
list pricing data is [enabled](#databricks-usage-enabled).

The account host can also be specified using the `DATABRICKS_HOST`
environment variable.

**NOTE:** The `DATABRICKS_HOST` environment variable can not be used to specify
_both_ the [instance name](https://docs.databricks.com/en/workspace/workspace-details.html#workspace-instance-names-urls-and-ids)
and the accounts API endpoint. To account for this, the environment variables
`DATABRICKS_WORKSPACEHOST` and `DATABRICKS_ACCOUNTHOST` environment variables
can be alternately used either separately or in combination with the
`DATABRICKS_HOST` environment variable to specify the
[instance name](https://docs.databricks.com/en/workspace/workspace-details.html#workspace-instance-names-urls-and-ids)
and the accounts API endpoint, respectively.

###### `accountId`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| Databricks account ID for the accounts API | string | conditional | N/a |

This parameter specifies the Databricks account ID. This is used by the
integration when constructing the URLs for account-level
[ReST API](https://docs.databricks.com/api/workspace/introduction) calls.

This parameter is required when the collection of Databricks billable usage and
list pricing data is [enabled](#databricks-usage-enabled).

###### `accessToken`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| Databricks personal access token | string | N | N/a |

When set, the integration will use [Databricks personal access token authentication](https://docs.databricks.com/en/dev-tools/auth/pat.html)
to authenticate Databricks API calls with the value of this parameter as the
Databricks [personal access token](https://docs.databricks.com/en/dev-tools/auth/pat.html#databricks-personal-access-tokens-for-workspace-users).

The personal access token can also be specified using the `DATABRICKS_TOKEN`
environment variable or any other SDK-supported mechanism (e.g. the `token`
field in a Databricks [configuration profile](https://docs.databricks.com/en/dev-tools/auth/config-profiles.html)).

See the [authentication section](#authentication) for more details.

**NOTE:** Databricks personal access tokens can only be used to collect data
at the workspace level. To collect account level data such as
[billable usage and list pricing data](#billable-usage--list-pricing-data),
OAuth authentication must be used instead.

###### `oauthClientId`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| Databricks OAuth M2M client ID | string | N | N/a |

When set, the integration will [use a service principal to authenticate with Databricks (OAuth M2M)](https://docs.databricks.com/en/dev-tools/auth/oauth-m2m.html)
when making Databricks API calls. The value of this parameter will be used as
the OAuth client ID.

The OAuth client ID can also be specified using the `DATABRICKS_CLIENT_ID`
environment variable or any other SDK-supported mechanism (e.g. the `client_id`
field in a Databricks [configuration profile](https://docs.databricks.com/en/dev-tools/auth/config-profiles.html)).

See the [authentication section](#authentication) for more details.

###### `oauthClientSecret`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| Databricks OAuth M2M client secret | string | N | N/a |

When the [`oauthClientId`](#oauthclientid) is set, this parameter can be set to
specify the [OAuth secret](https://docs.databricks.com/en/dev-tools/auth/oauth-m2m.html#step-3-create-an-oauth-secret-for-a-service-principal)
associated with the [service principal](https://docs.databricks.com/en/admin/users-groups/service-principals.html#what-is-a-service-principal).

The OAuth client secret can also be specified using the
`DATABRICKS_CLIENT_SECRET` environment variable or any other SDK-supported
mechanism (e.g. the `client_secret` field in a Databricks
[configuration profile](https://docs.databricks.com/en/dev-tools/auth/config-profiles.html)).

See the [authentication section](#authentication) for more details.

###### `sparkMetrics`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| Flag to enable automatic collection of Spark metrics | `true` / `false` | N | `true` |

**Deprecated** This configuration parameter has been deprecated in favor of the
configuration parameter [`databricks.spark.enabled`](#databricks-spark-enabled).
Use that parameter instead.

###### `sparkMetricPrefix`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| A prefix to prepend to Spark metric names | string | N | N/a |

**Deprecated** This configuration parameter has been deprecated in favor of the
configuration parameter [`databricks.spark.metricPrefix`](#databricks-spark-metricprefix).
Use that parameter instead.

###### `sparkClusterSources`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| The root node for the [Databricks cluster source configuration](#databricks-cluster-source-configuration) | YAML Mapping | N | N/a |

**Deprecated** This configuration parameter has been deprecated in favor of the
configuration parameter [`databricks.spark.clusterSources`](#databricks-spark-clustersources).
Use that parameter instead.

###### `sqlStatementTimeout`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| Timeout (in seconds) to use when executing SQL statements on a SQL warehouse | number | N | 30 |

Certain telemetry and data collected by the Databricks collector requires the
collector to run Databricks SQL statements on a SQL warehouse. This
configuration parameter specifies the number of seconds to wait before timing
out a pending or running SQL query.

###### `spark`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| The root node for the set of [Databricks Spark configuration](#databricks-spark-configuration) parameters | YAML Mapping | N | N/a |

This element groups together the configuration parameters to [configure](#databricks-spark-configuration)
the Databricks collector settings related to the collection of telemetry from
Databricks running on Spark. The configuration parameters in this group replace
the configuration parameters [`sparkMetrics`](#sparkmetrics),
[`sparkMetricPrefix`](#sparkmetricprefix), and [`sparkClusterSources`](#sparkclustersources).

###### `usage`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| The root node for the set of [Databricks Usage configuration](#databricks-usage-configuration) parameters | YAML Mapping | N | N/a |

This element groups together the configuration parameters to [configure](#databricks-usage-configuration)
the Databricks collector settings related to the collection of
[billable usage and list pricing data](#billable-usage--list-pricing-data).

##### Databricks `spark` configuration

###### Databricks Spark `enabled`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| Flag to enable automatic collection of Spark metrics | `true` / `false` | N | `true` |

By default, when the Databricks collector is enabled, it will automatically
collect Spark telemetry from Spark running on Databricks.

This flag can be used to disable the collection of Spark telemetry by the
Databricks collector. This may be useful to control data ingest when business
requirements call for the collection of non-Spark related Databricks telemetry
and Spark telemetry is not used. This flag is also used by the integration when
it is [deployed directly on the driver node of a Databricks cluster](#deploy-the-integration-on-the-driver-node-of-a-databricks-cluster) using the
the provided [init script](./init/cluster_init_integration.sh) since Spark
telemetry is collected by the Spark collector in this scenario.

**NOTE:** This configuration parameter replaces the older [`sparkMetrics`](#sparkmetrics)
configuration parameter.

###### Databricks Spark `metricPrefix`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| A prefix to prepend to Spark metric names | string | N | N/a |

This parameter serves the same purpose as the [`metricPrefix`](#metricprefix)
parameter of the [Spark configuration](#spark-configuration) except that it
applies to Spark telemetry collected by the Databricks collector. See the
[`metricPrefix`](#metricprefix) parameter of the [Spark configuration](#spark-configuration)
for more details.

Note that this parameter has no effect on Spark telemetry collected by the Spark
collector. This includes the case when the integration is
[deployed directly on the driver node of a Databricks cluster](#deploy-the-integration-on-the-driver-node-of-a-databricks-cluster)
using the the provided [init script](./init/cluster_init_integration.sh)
since Spark telemetry is collected by the Spark collector in this scenario.

**NOTE:** This configuration parameter replaces the older [`sparkMetricPrefix`](#sparkmetricprefix)
configuration parameter.

###### Databricks Spark `clusterSources`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| The root node for the [Databricks cluster source configuration](#databricks-cluster-source-configuration) | YAML Mapping | N | N/a |

The mechanism used to create a cluster is referred to as a cluster "source". The
Databricks collector supports collecting Spark telemetry from all-purpose
clusters created via the UI or API and from job clusters created via the
Databricks Jobs Scheduler. This element groups together the flags used to
individually [enable or disable](#databricks-cluster-source-configuration) the
cluster sources from which the Databricks collector will collect Spark
telemetry.

**NOTE:** This configuration parameter replaces the older [`sparkClusterSources`](#sparkclustersources)
configuration parameter.

##### Databricks cluster source configuration

###### `ui`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| Flag to enable automatic collection of Spark telemetry from all-purpose clusters created via the UI | `true` / `false` | N | `true` |

By default, when the Databricks collector is enabled, it will automatically
collect Spark telemetry from all all-purpose clusters created via the UI.

This flag can be used to disable the collection of Spark telemetry from
all-purpose clusters created via the UI.

###### `job`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| Flag to enable automatic collection of Spark telemetry from job clusters created via the Databricks Jobs Scheduler | `true` / `false` | N | `true` |

By default, when the Databricks collector is enabled, it will automatically
collect Spark telemetry from job clusters created by the Databricks Jobs
Scheduler.

This flag can be used to disable the collection of Spark telemetry from job
clusters created via the Databricks Jobs Scheduler.

###### `api`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| Flag to enable automatic collection of Spark telemetry from all-purpose clusters created via the [Databricks ReST API](https://docs.databricks.com/api/workspace/introduction) | `true` / `false` | N | `true` |

By default, when the Databricks collector is enabled, it will automatically
collect Spark telemetry from all-purpose clusters created via the [Databricks ReST API](https://docs.databricks.com/api/workspace/introduction).

This flag can be used to disable the collection of Spark telemetry from
all-purpose clusters created via the [Databricks ReST API](https://docs.databricks.com/api/workspace/introduction).

##### Databricks Usage Configuration

The Databricks usage configuration parameters are used to configure Databricks
collector settings related to the collection of Databricks
[billable usage and list pricing data](#billable-usage--list-pricing-data).

###### Databricks Usage `enabled`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| Flag to enable automatic collection of billable usage and list pricing data | `true` / `false` | N | `true` |

By default, when the Databricks collector is enabled, it will automatically
collect [billable usage and list pricing data](#billable-usage--list-pricing-data).

This flag can be used to disable the collection of billable usage and list
pricing data by the Databricks collector. This may be useful when running
multiple instances of the New Relic Databricks integration. In this scenario,
Databricks billable usage and list pricing data collection should _only_ be
enabled on a single instance. Otherwise, billable usage data will be recorded
more than once in New Relic, affecting consumption and cost calculations.

###### `warehouseId`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| ID of a SQL warehouse on which to run usage-related SQL statements | string | Y | N/a |

The ID of a SQL warehouse on which to run the SQL statements used to collect
Databricks [billable usage and list pricing data](#billable-usage--list-pricing-data).

This parameter is required when the collection of Databricks
[billable usage and list pricing data](#billable-usage--list-pricing-data) is
[enabled](#databricks-usage-enabled).

###### `includeIdentityMetadata`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| Flag to enable inclusion of identity related metadata in billable usage data | `true` / `false` | N | `false` |

When the collection of Databricks [billable usage and list pricing data](#billable-usage--list-pricing-data)
is [enabled](#databricks-usage-enabled), the Databricks collector can include
several pieces of identifying information along with the billable usage data.

By default, when the collection of Databricks [billable usage and list pricing](#billable-usage--list-pricing-data)
data is [enabled](#databricks-usage-enabled), the Databricks collector will
_not_ collect such data as it may be personally identifiable. This flag can be
used to enable the inclusion of the identifying information.

When enabled, the following values are included.

* The identity of the user a serverless billing record is attributed to. This
  value is included in the [identity metadata](https://docs.databricks.com/en/admin/system-tables/billing.html#analyze-identity-metadata)
  returned from usage records in the [billable usage system table](https://docs.databricks.com/en/admin/system-tables/billing.html).
* The identity of the cluster creator for each usage record for billable usage
  attributed to all-purpose and job compute.
* The single user name for each usage record for billable usage attributed to
  all-purpose and job compute configured for [single-user access mode](https://docs.databricks.com/en/compute/configure.html#access-mode).
* The identity of the warehouse creator for each usage record for billable usage
  attributed to SQL warehouse compute.

###### `runTime`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| Time of day (as `HH:mm:ss`) at which to run usage data collection | string with format `HH:mm:ss` | N | `02:00:00` |

This parameter specifies the time of day at which the collection of
[billable usage and list pricing data](#billable-usage--list-pricing-data)
occur. The value must be of the form `HH:mm:ss` where `HH` is the `0`-padded
24-hour clock hour (`00` - `23`), `mm` is the `0`-padded minute (`00` - `59`)
and `ss` is the `0`-padded second (`00` - `59`). For example, `09:00:00` is the
time 9:00 AM and `23:30:00` is the time 11:30 PM.

The time will _always_ be interpreted according to the UTC time zone. The time
zone can not be configured. For example, to specify that the integration should
be run at 2:00 AM EST (-0500), the value `07:00:00` should be specified.

##### Spark configuration

The Spark configuration parameters are used to configure the Spark collector.

###### `webUiUrl`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| The [Web UI](https://spark.apache.org/docs/3.5.2/web-ui.html) URL of an application on the Spark deployment to monitor | string | N | N/a |

This parameter can be used to monitor a non-Databricks Spark deployment. It
specifes the URL of the [Web UI](https://spark.apache.org/docs/3.5.2/web-ui.html)
of an application running on the Spark deployment to monitor. The value should
be of the form `http[s]://<hostname>:<port>` where `<hostname>` is the hostname
of the Spark deployment to monitor and `<port>` is the port number of the
Spark application's Web UI (typically 4040 or 4041, 4042, etc if more than one
application is running on the same host).

Note that the value must not contain a path. The path of the [Spark ReST API](https://spark.apache.org/docs/3.5.2/monitoring.html#rest-api)
endpoints (mounted at `/api/v1`) will automatically be prepended.

###### `metricPrefix`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| A prefix to prepend to Spark metric names | string | N | N/a |

This parameter specifies a prefix that will be prepended to each Spark metric
name when the metric is exported to New Relic.

For example, if this parameter is set to `spark.`, then the full name of the
metric representing the value of the memory used on application executors
(`app.executor.memoryUsed`) will be `spark.app.executor.memoryUsed`.

Note that it is not recommended to leave this value empty as the metric names
without a prefix may be ambiguous. Additionally, note that this parameter has no
effect on Spark telemetry collected by the Databricks collector. In that case,
use the [`sparkMetricPrefix`](#sparkmetricprefix) instead.

### Authentication

The Databricks integration uses the [Databricks SDK for Go](https://docs.databricks.com/en/dev-tools/sdk-go.html)
to access the Databricks and Spark ReST APIs. The SDK performs authentication on
behalf of the integration and provides many options for configuring the
authentication type and credentials to be used. See the
[SDK documentation](https://github.com/databricks/databricks-sdk-go?tab=readme-ov-file#authentication)
and the [Databricks client unified authentication documentation](https://docs.databricks.com/en/dev-tools/auth/unified-auth.html)
for details.

For convenience purposes, the following parameters can be used in the
[Databricks configuration](#databricks-configuration) section of the
[`config.yml](#configyml) file.

- [`accessToken`](#accesstoken) - When set, the integration will instruct the
  SDK to explicitly use [Databricks personal access token authentication](https://docs.databricks.com/en/dev-tools/auth/pat.html).
  The SDK will _not_ attempt to try other authentication mechanisms and instead
  will fail immediately if personal access token authentication fails.

  **NOTE:** Databricks personal access tokens can only be used to collect data
  at the workspace level. To collect account level data such as
  [billable usage and list pricing data](#billable-usage--list-pricing-data),
  OAuth authentication must be used instead.
- [`oauthClientId`](#oauthclientid) - When set, the integration will instruct
  the SDK to explicitly [use a service principal to authenticate with Databricks (OAuth M2M)](https://docs.databricks.com/en/dev-tools/auth/oauth-m2m.html).
  The SDK will _not_ attempt to try other authentication mechanisms and instead
  will fail immediately if OAuth M2M authentication fails. The OAuth Client
  secret can be set using the [`oauthClientSecret`](#oauthclientsecret)
  configuration parameter or any of the other mechanisms supported by the SDK
  (e.g. the `client_secret` field in a Databricks [configuration profile](https://docs.databricks.com/en/dev-tools/auth/config-profiles.html)
  or the `DATABRICKS_CLIENT_SECRET` environment variable).
- [`oauthClientSecret`](#oauthclientsecret) - The OAuth client secret to use for
  [OAuth M2M authentication](https://docs.databricks.com/en/dev-tools/auth/oauth-m2m.html).
  This value is _only_ used when [`oauthClientId`](#oauthclientid) is set in the
  [`config.yml`](#configyml). The OAuth client secret can also be set using any
  of the other mechanisms supported by the SDK (e.g. the `client_secret` field
  in a Databricks [configuration profile](https://docs.databricks.com/en/dev-tools/auth/config-profiles.html)
  or the `DATABRICKS_CLIENT_SECRET` environment variable).

### Billable Usage & List Pricing Data

The New Relic Databricks integration can collect Databricks
[billable usage](https://docs.databricks.com/en/admin/system-tables/billing.html)
and [list pricing](https://docs.databricks.com/en/admin/system-tables/pricing.html)
data from the Databricks [system tables](https://docs.databricks.com/en/admin/system-tables/index.html).
This data can be used to show basic Databricks DBU consumption and cost
metrics directly within New Relic.

When the [Databricks usage `enabled`](#databricks-usage-enabled) flag is set to
`true`, the Databricks collector will import billable usage records from the
[`system.billing.usage` table](https://docs.databricks.com/en/admin/system-tables/billing.html)
and list pricing records from the
[`system.billing.list_prices` table](https://docs.databricks.com/en/admin/system-tables/pricing.html)
once a day at the time specified in the [`runTime`](#runtime) configuration
parameter.

**NOTE:** In order for the New Relic Databricks integration to collect billing
usage and list pricing data, OAuth [authentication](#authentication) _must_ be
used. This is required even when the integration is
[deployed on the driver node of a Databricks cluster](#deploy-the-integration-on-the-driver-node-of-a-databricks-cluster)
using the provided [init script](./init/cluster_init_integration.sh) because the
billing usage and list pricing data can only be acquired using the account-level
[ReST API](https://docs.databricks.com/api/workspace/introduction) calls.

#### Billable Usage Data

On each run, billable usage data is collected _for the previous day_. For each
billable usage record, a corresponding record is created as a
[New Relic event](https://docs.newrelic.com/docs/data-apis/understand-data/new-relic-data-types/#event-data)
with the event type `DatabricksUsage` and the following attributes.

**NOTE:** Not every attribute is included in every event. For example, the
`cluster_*` attributes are only included in events for usage records relevant to
all-purpose or job related compute. Similarly, the `warehouse_*` attributes are
only included in events for usage records relevant to SQL warehouse related
compute.

**NOTE:** Descriptions below are sourced from the
[billable usage system table reference](https://docs.databricks.com/en/admin/system-tables/billing.html).

| Name | Description |
|---|---|
| `account_id` | ID of the account this usage record was generated for |
| `workspace_id` | ID of the Workspace this usage record was associated with |
| `workspace_name` | Name of the Workspace this usage record was associated with |
| `record_id` | Unique ID for this usage record |
| `sku_name` | Name of the SKU associated with this usage record |
| `cloud` | Name of the Cloud this usage record is relevant for |
| `usage_start_time` | The start time relevant to this usage record |
| `usage_end_time` | The end time relevant to this usage record |
| `usage_date` | Date of this usage record |
| `custom_tags` | Tags applied by the users to this usage record. Includes compute resource tags and jobs tags. |
| `usage_unit` | Unit this usage record is measured in |
| `usage_quantity` | Number of units consumed for this usage record |
| `record_type` | Whether the usage record is original, a retraction, or a restatement. See the section ["Analyze Correction Records"](https://docs.databricks.com/en/admin/system-tables/billing.html#analyze-correction-records) in the Databricks documentation for more details. |
| `ingestion_date` | Date the usage record was ingested into the usage table |
| `billing_origin_product` | The product that originated this usage reocrd |
| `usage_type` | The type of usage attributed to the product or workload for billing purposes |
| `cluster_id` | ID of the cluster associated with this usage record |
| `cluster_creator` | Creator of the cluster associated with this usage record (only included if [`includeIdentityMetadata`](#includeidentitymetadata) is `true`) |
| `cluster_single_user_name` | Single user name of the cluster associated with this usage record if the access mode of the cluster is [single-user access mode](https://docs.databricks.com/en/compute/configure.html#access-mode) (only included if [`includeIdentityMetadata`](#includeidentitymetadata) is `true`) |
| `cluster_source` | Cluster source of the cluster associated with this usage record |
| `cluster_instance_pool_id` | Instance pool ID of the cluster associated with this usage record |
| `warehouse_id` | ID of the SQL warehouse associated with this usage record |
| `warehouse_name` | Name of the SQL warehouse associated with this usage record |
| `warehouse_creator` | Creator of the SQL warehouse associated with this usage record (only included if [`includeIdentityMetadata`](#includeidentitymetadata) is `true`) |
| `instance_pool_id` | ID of the instance pool associated with this usage record |
| `node_type` | The instance type of the compute resource associated with this usage record |
| `job_id` | ID of the job associated with this usage record for serverless compute or jobs compute usage |
| `job_run_id` | ID of the job run associated with this usage record for serverless compute or jobs compute usage |
| `job_name` | User-given name of the job associated with this usage record for serverless compute or jobs compute usage |
| `notebook_id` | ID of the notebook associated with this usage record for serverless compute for notebook usage  |
| `notebook_path` | Workspace storage path of the notebook associated with this usage for serverless compute for notebook usage |
| `dlt_pipeline_id` | ID of the Delta Live Tables pipeline associated with this usage record |
| `dlt_update_id` | ID of the Delta Live Tables pipeline update associated with this usage record |
| `dlt_maintenance_id` | ID of the Delta Live Tables pipeline maintenance tasks associated with this usage record |
| `run_name` | Unique user-facing identifier of the Mosaic AI Model Training fine-tuning run associated with this usage record  |
| `endpoint_name` | Name of the model serving endpoint or vector search endpoint associated with this usage record |
| `endpoint_id` | ID of the model serving endpoint or vector search endpoint associated with this usage record |
| `central_clean_room_id` | ID of the central clean room associated with this usage record |
| `run_as` | See the section ["Analyze Identity Metadata"](https://docs.databricks.com/en/admin/system-tables/billing.html#analyze-identity-metadata) in the Databricks documentation for more details (only included if [`includeIdentityMetadata`](#includeidentitymetadata) is `true`) |
| `jobs_tier` | Jobs tier product features for this usage record: values include `LIGHT`, `CLASSIC`, or `null` |
| `sql_tier` | SQL tier product features for this usage record: values include `CLASSIC`, `PRO`, or `null` |
| `dlt_tier` | DLT tier product features for this usage record: values include `CORE`, `PRO`, `ADVANCED`, or `null` |
| `is_serverless` | Flag indicating if this usage record is associated with serverless usage: values include `true` or `false`, or `null`  |
| `is_photon` | Flag indicating if this usage record is associated with Photon usage: values include `true` or `false`, or `null` |
| `serving_type` | Serving type associated with this usage record: values include `MODEL`, `GPU_MODEL`, `FOUNDATION_MODEL`, `FEATURE`, or `null` |

#### List Pricing Data

On each run, list pricing data is gathered into a
[New Relic lookup table](https://docs.newrelic.com/docs/logs/ui-data/lookup-tables-ui/)
named `DatabricksListPrices`. The entire lookup table is updated on each run.
For each pricing record, this table will contain a corresponding row with the
following columns.

**NOTE:** Descriptions below are sourced from the
[pricing system table reference](https://docs.databricks.com/en/admin/system-tables/pricing.html).

| Name | Description |
|---|---|
| account_id | ID of the account this pricing record was generated for |
| price_start_time | The time the price in this pricing record became effective in UTC |
| price_end_time | The time the price in this pricing record stopped being effective in UTC |
| sku_name | Name of the SKU associated with this pricing record |
| cloud | Name of the Cloud this pricing record is relevant for |
| currency_code | The currency the price in this pricing record is expressed in |
| usage_unit | The unit of measurement that is monetized in this pricing record |
| list_price | A single price that can be used for simple long-term estimates |
| promotional_price | A temporary promotional price that all customers get which could be used for cost estimation during the temporary period |
| effective_list_price | The effective list price used for calculating the cost |

## Building

### Coding Conventions

#### Style Guidelines

While not strictly enforced, the basic preferred editor settings are set in the
[.editorconfig](./.editorconfig). Other than this, no style guidelines are
currently imposed.

#### Static Analysis

This project uses both [`go vet`](https://pkg.go.dev/cmd/vet) and
[`staticcheck`](https://staticcheck.io/) to perform static code analysis. These
checks are run via [`precommit`](https://pre-commit.com) on all commits. Though
this can be bypassed on local commit, both tasks are also run during
[the `validate` workflow](./.github/workflows/validate.yml) and must have no
errors in order to be merged.

#### Commit Messages

Commit messages must follow [the conventional commit format](https://www.conventionalcommits.org/en/v1.0.0/).
Again, while this can be bypassed on local commit, it is strictly enforced in
[the `validate` workflow](./.github/workflows/validate.yml).

The basic commit message structure is as follows.

```
<type>[optional scope][!]: <description>

[optional body]

[optional footer(s)]
```

In addition to providing consistency, the commit message is used by
[svu](https://github.com/caarlos0/svu) during
[the release workflow](./.github/workflows/release.yml). The presence and values
of certain elements within the commit message affect auto-versioning. For
example, the `feat` type will bump the minor version. Therefore, it is important
to use the guidelines below and carefully consider the content of the commit
message.

Please use one of the types below.

- `feat` (bumps minor version)
- `fix` (bumps patch version)
- `chore`
- `build`
- `docs`
- `test`

Any type can be followed by the `!` character to indicate a breaking change.
Additionally, any commit that has the text `BREAKING CHANGE:` in the footer will
indicate a breaking change.

### Local Development

For local development, simply use `go build` and `go run`. For example,

```bash
go build cmd/databricks/databricks.go
```

Or

```bash
go run cmd/databricks/databricks.go
```

If you prefer, you can also use [`goreleaser`](https://goreleaser.com/) with
the `--single-target` option to build the binary for the local `GOOS` and
`GOARCH` only.

```bash
goreleaser build --single-target
```

### Releases

Releases are built and packaged using [`goreleaser`](https://goreleaser.com/).
By default, a new release will be built automatically on any push to the `main`
branch. For more details, review the [`.goreleaser.yaml`](./.goreleaser.yaml)
and [the `goreleaser` documentation](https://goreleaser.com/intro/).

The [svu](https://github.com/caarlos0/svu) utility is used to generate the next
tag value [based on commit messages](https://github.com/caarlos0/svu#commit-messages-vs-what-they-do).

### GitHub Workflows

This project utilizes GitHub workflows to perform actions in response to
certain GitHub events.

| Workflow | Events | Description
| --- | --- | --- |
| [validate](./.github/workflows/validate.yml) | `push`, `pull_request` to `main` branch | Runs [precommit](https://pre-commit.com) to perform static analysis and runs [commitlint](https://commitlint.js.org/#/) to validate the last commit message |
| [build](./.github/workflows/build.yml) | `push`, `pull_request` | Builds and tests code |
| [release](./.github/workflows/release.yml) | `push` to `main` branch | Generates a new tag using [svu](https://github.com/caarlos0/svu) and runs [`goreleaser`](https://goreleaser.com/) |
| [repolinter](./.github/workflows/repolinter.yml) | `pull_request` | Enforces repository content guidelines |

## Appendix

The sections below cover topics that are related to Databricks telemetry but
that are not specifically part of this integration. In particular, any assets
referenced in these sections must be installed and/or managed _separately_ from
the integration. For example, the init scripts provided to [monitor cluster health](#monitoring-cluster-health)
are not automatically installed or used by the integration.

### Monitoring Cluster Health

[New Relic Infrastructure](https://docs.newrelic.com/docs/infrastructure/infrastructure-monitoring/get-started/get-started-infrastructure-monitoring/)
can be used to collect system metrics like CPU and memory usage from the nodes in
a Databricks [cluster](https://docs.databricks.com/en/getting-started/concepts.html#cluster).
Additionally, [New Relic APM](https://docs.newrelic.com/docs/apm/new-relic-apm/getting-started/introduction-apm/)
can be used to collect application metrics like JVM heap size and GC cycle count
from the [Apache Spark](https://spark.apache.org/docs/latest/index.html) driver
and executor JVMs. Both are achieved using [cluster-scoped init scripts](https://docs.databricks.com/en/init-scripts/cluster-scoped.html).
The sections below cover the installation of these init scripts.

**NOTE:** Use of one or both init scripts will have a slight impact on cluster
startup time. Therefore, consideration should be given when using the init
scripts with a job cluster, particularly when using a job cluster scoped to a
single task.

#### Configure the New Relic license key

Both the [New Relic Infrastructure Agent init script](#install-the-new-relic-infrastructure-agent-init-script)
and the [New Relic APM Java Agent init script](#install-the-new-relic-apm-java-agent-init-script)
require a [New Relic license key](https://docs.newrelic.com/docs/apis/intro-apis/new-relic-api-keys/#overview-keys)
to be specified in a [custom environment variable](https://docs.databricks.com/en/compute/configure.html#environment-variables)
named `NEW_RELIC_LICENSE_KEY`. While the license key _can_ be specified by
hard-coding it in plain text in the compute configuration, this is not
recommended. Instead, it is recommended to create a [secret](https://docs.databricks.com/en/security/secrets/secrets.html).
using the [Databricks CLI](https://docs.databricks.com/en/dev-tools/cli/index.html)
and [reference the secret in the environment variable](https://docs.databricks.com/en/security/secrets/secrets.html#reference-a-secret-in-an-environment-variable).

To create the secret and set the environment variable, perform the following
steps.

1. Follow the steps to [install or update the Databricks CLI](https://docs.databricks.com/en/dev-tools/cli/install.html).
1. Use the Databricks CLI to create a [Databricks-backed secret scope](https://docs.databricks.com/en/security/secrets/secret-scopes.html#create-a-databricks-backed-secret-scope)
   with the name `newrelic`. For example,

   ```bash
   databricks secrets create-scope newrelic
   ```

   **NOTE:** Be sure to take note of the information in the referenced URL about
   the `MANAGE` scope permission and use the correct version of the command.
1. Use the Databricks CLI to create a secret for the license key in the new
   scope with the name `licenseKey`. For example,

   ```bash
   databricks secrets put-secret --json '{
      "scope": "newrelic",
      "key": "licenseKey",
      "string_value": "[YOUR_LICENSE_KEY]"
   }'
   ```

To set the custom environment variable named `NEW_RELIC_LICENSE_KEY` and
reference the value from the secret, follow the steps to
[configure custom environment variables](https://docs.databricks.com/en/compute/configure.html#environment-variables)
and add the following line after the last entry in the `Environment variables`
field.

`NEW_RELIC_LICENSE_KEY={{secrets/newrelic/licenseKey}}`

#### Install the New Relic Infrastructure Agent init script

The [`cluster_init_infra.sh`](./examples/cluster_init_infra.sh) script
automatically installs the latest version of the
[New Relic Infrastructure Agent](https://docs.newrelic.com/docs/infrastructure/infrastructure-monitoring/get-started/get-started-infrastructure-monitoring/)
on each node of the cluster.

To install the init script, perform the following steps.

1. Login to your Databricks account and navigate to the desired
   [workspace](https://docs.databricks.com/en/getting-started/concepts.html#accounts-and-workspaces).
1. Follow [the recommendations for init scripts](https://docs.databricks.com/en/init-scripts/index.html#recommendations-for-init-scripts)
   to store the [`cluster_init_infra.sh`](./examples/cluster_init_infra.sh)
   script within your workspace in the recommended manner. For example, if your
   workspace is [enabled for Unity Catalog](https://docs.databricks.com/en/data-governance/unity-catalog/get-started.html#step-1-confirm-that-your-workspace-is-enabled-for-unity-catalog),
   you should store the init script in a [Unity Catalog volume](https://docs.databricks.com/en/ingestion/file-upload/upload-to-volume.html).
1. Navigate to the [`Compute`](https://docs.databricks.com/en/compute/clusters-manage.html#view-compute)
   tab and select the desired all-purpose or job compute to open the compute
   details UI.
1. Click the button labeled `Edit` to [edit](https://docs.databricks.com/en/compute/clusters-manage.html#edit-a-compute)
   the compute's configuration.
1. Follow the steps to [use the UI to configure a cluster to run an init script](https://docs.databricks.com/en/init-scripts/cluster-scoped.html#configure-a-cluster-scoped-init-script-using-the-ui)
   and point to the location where you stored the init script in step 2.
1. If your cluster is not running, click on the button labeled `Confirm` to
   save your changes. Then, restart the cluster. If your cluster is already
   running, click on the button labeled `Confirm and restart` to save your
   changes and restart the cluster.

#### Install the New Relic APM Java Agent init script

The [`cluster_init_apm.sh`](./examples/cluster_init_apm.sh) script
automatically installs the latest version of the
[New Relic APM Java Agent](https://docs.newrelic.com/docs/apm/agents/java-agent/getting-started/introduction-new-relic-java/)
on each node of the cluster.

To install the init script, perform the same steps as outlined in the
[Install the New Relic Infrastructure Agent init script](#install-the-new-relic-infrastructure-agent-init-script)
section using the [`cluster_init_apm.sh`](./examples/cluster_init_apm.sh) script
instead of the [`cluster_init_infra.sh`](./examples/cluster_init_infra.sh)
script.

Additionally, perform the following steps.

1. Login to your Databricks account and navigate to the desired
   [workspace](https://docs.databricks.com/en/getting-started/concepts.html#accounts-and-workspaces).
1. Navigate to the [`Compute`](https://docs.databricks.com/en/compute/clusters-manage.html#view-compute)
   tab and select the desired all-purpose or job compute to open the compute
   details UI.
1. Click the button labeled `Edit` to [edit](https://docs.databricks.com/en/compute/clusters-manage.html#edit-a-compute)
   the compute's configuration.
1. Follow the steps to [configure custom Spark configuration properties](https://docs.databricks.com/en/compute/configure.html#spark-configuration)
   and add the following lines after the last entry in the `Spark Config` field.

   ```
   spark.driver.extraJavaOptions -javaagent:/databricks/jars/newrelic-agent.jar
   spark.executor.extraJavaOptions -javaagent:/databricks/jars/newrelic-agent.jar -Dnewrelic.tempdir=/tmp
   ```

1. If your cluster is not running, click on the button labeled `Confirm` to
   save your changes. Then, restart the cluster. If your cluster is already
   running, click on the button labeled `Confirm and restart` to save your
   changes and restart the cluster.

#### Viewing your cluster data

With the New Relic Infrastructure Agent init script installed, a host entity
will show up for each node in the cluster.

With the New Relic APM Java Agent init script installed, an APM application
entity named `Databricks Driver` will show up for the Spark driver JVM and an
APM application entity named `Databricks Executor` will show up for the
executor JVMs. Note that all executor JVMs will report to a single APM
application entity. Metrics for a specific executor can be viewed on many pages
of the [APM UI](https://docs.newrelic.com/docs/apm/apm-ui-pages/monitoring/response-time-chart-types-apm-browser/)
by selecting the instance from the `Instances` menu located below the time range
selector. On the [JVM Metrics page](https://docs.newrelic.com/docs/apm/agents/java-agent/features/jvms-page-java-view-app-server-metrics-jmx/),
the JVM metrics for a specific executor can be viewed by selecting an instance
from the `JVM instances` table.

Additionally, both the host entities and the APM entities are tagged with the
tags listed below to make it easy to filter down to the entities that make up
your cluster using the [entity filter bar](https://docs.newrelic.com/docs/new-relic-solutions/new-relic-one/core-concepts/search-filter-entities/)
that is available in many places in the UI.

* `databricksClusterId` - The ID of the Databricks cluster
* `databricksClusterName` - The name of the Databricks cluster
* `databricksIsDriverNode` - `true` if the entity is on the driver node,
   otherwise `false`
* `databricksIsJobCluster` - `true` if the entity is part of a [job cluster](https://docs.databricks.com/en/jobs/use-compute.html),
   otherwise `false`

Below is an example of using the `databricksClusterName` to filter down to the
host and APM entities for a single cluster using the [entity filter bar](https://docs.newrelic.com/docs/new-relic-solutions/new-relic-one/core-concepts/search-filter-entities/)
on the `All entities` view.

![infra and apm cluster filter example](./examples/cluster-infra-apm.png)

## Support

New Relic has open-sourced this project. This project is provided AS-IS WITHOUT
WARRANTY OR DEDICATED SUPPORT. Issues and contributions should be reported to
the project here on GitHub.

We encourage you to bring your experiences and questions to the
[Explorers Hub](https://discuss.newrelic.com/) where our community members
collaborate on solutions and new ideas.

### Privacy

At New Relic we take your privacy and the security of your information
seriously, and are committed to protecting your information. We must emphasize
the importance of not sharing personal data in public forums, and ask all users
to scrub logs and diagnostic information for sensitive information, whether
personal, proprietary, or otherwise.

We define “Personal Data” as any information relating to an identified or
identifiable individual, including, for example, your name, phone number, post
code or zip code, Device ID, IP address, and email address.

For more information, review [New Relic’s General Data Privacy Notice](https://newrelic.com/termsandconditions/privacy).

### Contribute

We encourage your contributions to improve this project! Keep in mind that
when you submit your pull request, you'll need to sign the CLA via the
click-through using CLA-Assistant. You only have to sign the CLA one time per
project.

If you have any questions, or to execute our corporate CLA (which is required
if your contribution is on behalf of a company), drop us an email at
opensource@newrelic.com.

**A note about vulnerabilities**

As noted in our [security policy](../../security/policy), New Relic is committed
to the privacy and security of our customers and their data. We believe that
providing coordinated disclosure by security researchers and engaging with the
security community are important means to achieve our security goals.

If you believe you have found a security vulnerability in this project or any of
New Relic's products or websites, we welcome and greatly appreciate you
reporting it to New Relic through [HackerOne](https://hackerone.com/newrelic).

If you would like to contribute to this project, review [these guidelines](./CONTRIBUTING.md).

To all contributors, we thank you!  Without your contribution, this project
would not be what it is today.

### License

The New Relic Databricks Integration project is licensed under the
[Apache 2.0](http://apache.org/licenses/LICENSE-2.0.txt) License.
