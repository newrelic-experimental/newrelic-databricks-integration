<a href="https://opensource.newrelic.com/oss-category/#new-relic-experimental"><picture><source media="(prefers-color-scheme: dark)" srcset="https://github.com/newrelic/opensource-website/raw/main/src/images/categories/dark/Experimental.png"><source media="(prefers-color-scheme: light)" srcset="https://github.com/newrelic/opensource-website/raw/main/src/images/categories/Experimental.png"><img alt="New Relic Open Source experimental project banner." src="https://github.com/newrelic/opensource-website/raw/main/src/images/categories/Experimental.png"></picture></a>


# New Relic Databricks Integration

Welcome to the New Relic Integration for Databricks! This repository provides scripts and instructions for integrating New Relic with Databricks through New Relic Datbaricks Integration or through other Open Source tools like OpenTelemetry and Prometheus.

## Overview

This repository provides you various ways to utilize New Relic directly from your Databricks environment. With options to monitor via standalone integration, OpenTelemetry (OTel), or Prometheus, you flexibly choose what best fits your operational needs.

-  **New Relic Databricks Integration:** A direct connection between Databricks and New Relic, enabling seamless data transfer and analysis capabilities. This integration supports `spark metrics`, `databricks queries metrics`, `databricks job runs events`. This integration along with New Relic APM agent can pull `logs`, and cluster performance related data as well.


-  **OpenTelemetry (OTel) Integration:** An open-source observability framework, enabling you to generate and manage telemetry data, supports `spark metrics` from Databricks. Please follow the [instructions here](/opentelemetry/README.md) for a detailed guide on how to add initialization scripts for OpenTelemetry to your Databricks cluster.


-  **Prometheus Integration:** A powerful open-source systems monitoring and alerting toolkit which can process metrics from Databricks. Support `spark metrics` from Databricks. Please follow the [instructions here](link-to-the-instruction-page) for a detailed guide on how to add initialization scripts to your Databricks cluster.

Pick the option that suits your use-case and follow the associated guide to get started.


## Setup New Relic Databricks Integration

The Standalone environment runs the data pipelines as an independant service, either on-premises or cloud instances like AWS EC2. It can run on Linux, macOS, Windows, and any OS with support for GoLang.

### Prerequisites

- Go 1.20 or later.

### Build

Open a terminal, CD to `cmd/standalone`, and run:

```
$ go build
```

### Configuring the Pipeline

The standalone environment requieres a YAML file for pipeline configuration. The requiered keys are:

- `interval`: Integer. Time in seconds between requests.

Check `config/example_config.yaml` for a configuration example.


#### New Relic APIs exporter

- `nr_account_id`: String. Account ID.
- `nr_api_key`: String. Api key for writing.
- `nr_endpoint`: String. New Relic endpoint region. Either `US` or `EU`. Optional, default value is `US`.

### Running the Pipeline

Just run the following command from the build folder:

```
$ ./standalone path/to/config.yaml
```

To run the pipeline on system start, check your specific system init documentation.


### Adding Initialization Scripts to Databricks

Databricks Initialization Scripts are shell scripts that run when a cluster is starting. They are useful for setting up custom configurations or third-party integrations such as setting up monitoring agents. Here is how you add an init script to Databricks.

Based on the cloud Databricks is hosted on, you will be able to run the APM agent.

### Install New Relic APM on Databricks



1. **Add script to Databricks:** Create new file in workspace as nr-agent-installation.sh and add the below script to it.
   ```bash
   #!/bin/bash

   # Define the newrelic version and jar path
   
   NR_VERSION="8.10.0" # replace with the version you want
   NR_JAR_PATH="/databricks/jars/newrelic-agent-${NR_VERSION}.jar"
   NR_CONFIG_FILE="/databricks/jars/newrelic.yml"

   # Download the newrelic java agent
   curl -o ${NR_JAR_PATH} -L https://download.newrelic.com/newrelic/java-agent/newrelic-agent/${NR_VERSION}/newrelic-agent-${NR_VERSION}.jar

   # Create new relic yml file
   echo "common: &default_settings
   license_key: 'xxxxxx' # Replace with your License Key
   agent_enabled: true

   production:
   <<: *default_settings
   app_name: Databricks" > ${NR_CONFIG_FILE}
   ```

2. **Add the script to your Databricks cluster:** To add the initialization script to your cluster in Databricks, follow these steps:

   - Navigate to your Databricks workspace and go to the `Clusters` page.
   - Choose the cluster you want to add the script to and click `Edit`.
   - In the `Advanced Options` section, find the `Init Scripts` field.
   - Click on `Add`, then in the Script Path input, select workspace or cloud storage path where your script is stored.
   - Click `Confirm` and then `Update`.


3. **Add Spark configurations to attach the java agent:**

   * Navigate to your cluster `Advanced Options`, then `Spark`.
   - Add or update Spark configurations as key-value pairs. Here's an example:

    ```bash
    # Example jar path "/databricks/jars/newrelic-agent-8.10.0.jar"
   
    echo "spark.driver.extraJavaOptions -javaagent:${NR_JAR_PATH}"
    echo "spark.executor.extraJavaOptions -javaagent:${NR_JAR_PATH}"
    ```

4. **Verify the script was executed:** After your cluster starts/restarts, you should verify that the script was executed successfully. You can do this by checking the cluster logs via the `Logs` tab on your clusters page.

***Note***: Any changes to the script settings will apply only to new clusters or when existing clusters are restarted.

## Support

New Relic hosts and moderates an online forum where customers can interact with New Relic employees as well as other customers to get help and share best practices. If you're running into a problem, please raise an issue on this repository and we will try to help you ASAP. Please bear in mind this is an open source project and hence it isn't directly supported by New Relic.

## Contributing
We encourage your contributions to improve <strong>newrelic-databricks-integration</strong> . Keep in mind when you submit your pull request, you'll need to sign the CLA via the click-through using CLA-Assistant. You only have to sign the CLA one time per project.
If you have any questions, or to execute our corporate CLA, required if your contribution is on behalf of a company,  please drop us an email at opensource@newrelic.com.

**A note about vulnerabilities**

As noted in our [security policy](../../security/policy), New Relic is committed to the privacy and security of our customers and their data. We believe that providing coordinated disclosure by security researchers and engaging with the security community are important means to achieve our security goals.

If you believe you have found a security vulnerability in this project or any of New Relic's products or websites, we welcome and greatly appreciate you reporting it to New Relic through [HackerOne](https://hackerone.com/newrelic).

## License
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.