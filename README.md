# New Relic Databricks Integration

Welcome to the New Relic Integration for Databricks! This repository provides scripts and instructions for integrating New Relic with Databricks through New Relic Datbaricks Integration or through other Open Source tools like OpenTelemetry and Prometheus.

## Overview

This integration enables you to utilize New Relic directly within your Databricks environment. With options to monitor via standalone integration, OpenTelemetry (OTel), or Prometheus, you flexibly choose what best fits your operational needs.

-  **New Relic Databricks Integration:** A direct connection between Databricks and New Relic, enabling seamless data transfer and analysis capabilities. This integration supports `spark metrics`, `databricks queries metrics`, `databricks job runs events`. This integration along with New Relic APM agent can pull `logs`, and cluster performance related data as well. Depending on the cloud you are using for Databricks, you will be able to install New Relic Infra Agent.

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

Based on the cloud Databricks is hosted on, you will be able to run the APM agent and Infra agent accordingly.

### Install New Relic APM on Databricks



1. **Add script to Databricks:** Create new file in workspace same as /init-scripts/nr-agent-installation.sh.
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
2. **Add Spark configurations to attach the java agent:**
    ```bash
   # Example jar path "/databricks/jars/newrelic-agent-8.10.0.jar"
   
    echo "spark.driver.extraJavaOptions -javaagent:${NR_JAR_PATH}"
    echo "spark.executor.extraJavaOptions -javaagent:${NR_JAR_PATH}"
    ```
3. **Add the script to your Databricks cluster:** To add the initialization script to your cluster in Databricks, follow these steps:

    - Navigate to your Databricks workspace and go to the `Clusters` page.
    - Choose the cluster you want to add the script to and click `Edit`.
    - In the `Advanced Options` section, find the `Init Scripts` field.
    - Click on `Add`, then in the Script Path input, select workspace or cloud storage path where your script is stored.
    - Click `Confirm` and then `Update`.
   

4. **Verify the script was executed:** After your cluster starts/restarts, you should verify that the script was executed successfully. You can do this by checking the cluster logs via the `Logs` tab on your clusters page.

***Note***: Any changes to the script settings will apply only to new clusters or when existing clusters are restarted.

