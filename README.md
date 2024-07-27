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

This integration collects Spark telemetry, Workflow telemetry, and Cost and
Billing information from Databricks.

![Apache Spark Dashboard Screenshot](./examples/spark-dashboard-jobs.png)

## Table of Contents

* [Getting Started](#getting-started)
   * [On-host Deployment](#on-host)
* [Features](#features)
* [Usage](#usage)
   * [Command Line Options](#command-line-options)
   * [Configuration](#configuration)
      * [General configuration](#general-configuration)
      * [Pipeline configuration](#pipeline-configuration)
      * [Log configuration](#log-configuration)
      * [Query configuration](#query-configuration)

## Getting Started

To get started with the New Relic Databricks integration, [deploy the integration on a supported host environment](#on-host) and then [import](https://docs.newrelic.com/docs/query-your-data/explore-query-data/dashboards/dashboards-charts-import-export-data/#import-json) the [sample dashboard](./examples/spark-daskboard.json).

**NOTE:** [On-host deployment](#on-host) is currently the only supported deployment type.

### On-host

The New Relic Databricks integration provides binaries for the following
host platforms.

* Linux amd64
* Windows amd64

#### Deploy the integration on host

To run the Databricks integration on a host, perform the following steps.

1. Download the appropriate binary for your platform from the [latest release](https://github.com/newrelic-experimental/newrelic-databricks-integration/releases).
1. Create a directory named `configs` in the directory containing the binary.
1. Create a file named `config.yml` in the `configs` directory and copy the
   contents of the file [`configs/config.template.yml`](./configs/config.template.yml)
   in this repository into it.
1. Edit the `config.yml` file to [configure](#configuration) the integration
   appropriately for your environment.
1. Using the directory containing the binary as the current working directory,
   execute the binary using the appropriate [Command Line Options](#command-line-options).

## Features

The New Relic Databricks integration supports the following capabilities.

* [Collect Spark telemetry](#event-log-files)

  The New Relic Databricks integration can collect telemetry from Spark running
  on Databricks. By default, the integration will automatically connect to
  and collect telemetry from the Spark deployments in all clusters created via
  the UI or API in the specified workspace.

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

#### `config.yml`

The Databricks integration is configured using a YAML file named [`config.yml`](#configyml).
The default location for this file is `configs/config.yml` relative to the
current working directory when the integration binary is executed. The supported
configuration parameters are listed below. See [`config.template.yml`](./configs/config.template.yml)
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
| The root node for the set of [pipeline configuration](#pipeline-configuration) parameters | YAML Sequence | N | N/a |

The integration retrieves, processes, and exports data to New Relic using
a data pipeline consisting of one or more receivers, a processing chain, and a
New Relic exporter. Various aspects of the pipeline are configurable. This
element groups together the configuration parameters related to
[pipeline configuration](#pipeline-configuration).

###### `log`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| The root node for the set of [log configuration](#log-configuration) parameters | YAML Sequence | N | N/a |

The integration uses the [logrus](https://pkg.go.dev/github.com/sirupsen/logrus)
package for application logging. This element groups together the configuration
parameters related to [log configuration](#log-configuration).

###### `mode`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| The integration execution mode | `databricks` | N | `databricks` |

The integration execution mode. Currently, the only supported execution mode is
`databricks`.

###### `databricks`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| The root node for the set of [Databricks configuration](#databricks-configuration) parameters | YAML Sequence | N | N/a |

This element groups together the configuration parameters related to
[Databricks configuration](#databricks-configuration).

The configuration parameters in this section are only needed when collecting
Databricks specific data such as workflow and job telemetry or Databricks
account information or when collecting telemetry from Spark running on a
Databricks cluster.

###### `spark`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| The root node for the set of [Spark configuration](#spark-configuration) parameters | YAML Sequence | N | N/a |

This element groups together the configuration parameters related to
[Spark configuration](#spark-configuration).

The configuration parameters in this section can be used with or without the
[`databricks`](#databricks) configuration parameters depending on whether the
target Spark instance(s) are running on Databricks or not. See the
[Databricks configuration](#databricks-configuration) section for details on
collecting Spark telemetry from Spark running on Databricks.

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

###### `accessToken`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| Databricks access token | string | Y | N/a |

This parameter specifies the Databricks [personal access token](https://docs.databricks.com/en/dev-tools/auth/pat.html)
used by the integration to authenticate Databricks API calls.

The access token can also be specified using the `DATABRICKS_ACCESSTOKEN`
environment variable.

###### `workspaceHost`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| Databricks workspace instance name | string | Y | N/a |

This parameter specifies the [instance name](https://docs.databricks.com/en/workspace/workspace-details.html)
of the target Databricks instance for which data should be collected. This is
used by the integration when constructing the URLs for API calls. Note that the
value of this parameter _must not_ include the `https://` prefix, e.g.
`https://my-databricks-instance-name.cloud.databricks.com`.

The workspace host can also be specified using the `DATABRICKS_WORKSPACEHOST`
environment variable.

##### Spark configuration

###### `metricPrefix`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| A prefix to prepend to Spark metric names | string | N | `spark.` |

This parameter specifies a prefix that will be prepended to each Spark metric
name when the metric is exported to New Relic.

For example, if this parameter is set to `spark.` (the default), then the full
name of the metric representing the value of the memory used on application
executors (`app.executor.memoryUsed`) will be `spark.app.executor.memoryUsed`.

**NOTE:** It is not recommended to leave this value empty as the metric names
without a prefix may be ambiguous.

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
