## OpenTelemetry Collector Integration with Databricks

This Readme provides step-by-step instructions for setting up the OpenTelemetry Collector using `ApacheSpark` receiver for monitoring Databricks.

### Adding Initialization Scripts to Databricks

Databricks Initialization Scripts are shell scripts that run when a cluster is starting. They are useful for setting up custom configurations or third-party integrations such as setting up monitoring agents.

### Step 1: Download and Extract OpenTelemetry Collector

1. **Add script to Databricks:** Create new file in workspace as otel-installation.sh and add the below script to download and extract the OpenTelemetry Collector Contrib archive.

```bash
curl --proto '=https' --tlsv1.2 -fOL https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/v0.95.0/otelcol-contrib_0.95.0_linux_386.tar.gz
tar -xvf otelcol-contrib_0.95.0_linux_386.tar.gz
```

### Step 2: Create the Configuration File
Create a ﻿config.yaml file for the OpenTelemetry Collector with your desired configuration. Below is an example configuration, replace ﻿$DRIVER_HOST and ﻿$SPARKUIPORT with specific values related to your environment.

```bash
receivers:
  apachespark:
    collection_interval: 5s
    endpoint: http://$DRIVER_HOST:$SPARKUIPORT

exporters:
  otlp:
    endpoint: <nr-otlp-endpoint>
    headers:
      api-key: <nr-api-key>

service:
  pipelines:
    metrics:
      receivers: [apachespark]
      exporters: [otlp]
```

### Step 3: Configure Initialization Script
* Add the script to your Databricks cluster:** To add the initialization script to your cluster in Databricks, follow these steps.

  - Navigate to your Databricks workspace and go to the `Clusters` page.
  - Choose the cluster you want to add the script to and click `Edit`.
  - In the `Advanced Options` section, find the `Init Scripts` field.
  - Click on `Add`, then in the Script Path input, select workspace or cloud storage path where your script is stored.
  - Click `Confirm` and then `Update`.
