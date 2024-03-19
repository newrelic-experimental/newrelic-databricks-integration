# Multi Environment NRI Framework

## Standalone

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
- `exporter`: String. Exporter type, possible values are: `nrapi`, `nrmetrics`, `nrevents`, `nrlogs`, `nrtraces`, `otel` and `prom`.

Check `config/example_config.yaml` for a configuration example.

Each one of the exporters also require specific configuration.

#### New Relic APIs exporter

For exporter types `nrapi`, `nrmetrics`, `nrevents`, `nrlogs`, and `nrtraces`, the following keys are used:

- `nr_account_id`: String. Account ID. Only requiered for `nrevents` and `nrapi` exporters.
- `nr_api_key`: String. Api key for writing.
- `nr_endpoint`: String. New Relic endpoint region. Either `US` or `EU`. Optional, default value is `US`.

#### OpenTelemetry exporter

For exporter type `otel` the following keys are used:

- `otel_endpoint`: String. Domain and (optionally) port only, no scheme or path. Example: `my.otel.example.com:4318`
- `otel_scheme`: String. Either `http` or `https`. If not specified, default value is `https`.
- `otel_headers`: Map. List of key-value pairs sent as HTTP request headers. Optional.
- `otel_metric_endpoint`: String. Full URL of an endpoint to send metrics. If specified, `otel_endpoint` and `otel_scheme` will be ignored.
- `otel_logs_endpoint`: String. Full URL of an endpoint to send logs. If specified, `otel_endpoint` and `otel_scheme` will be ignored.

#### Prometheus exporter

For exporter type `prom` the following keys are used:

- `prom_endpoint`: String. Prometheus remote write URL.
- `prom_credentials`: String. Authorization credentials. Optional.
- `prom_headers`: Map. List of key-value pairs sent as HTTP request headers. Optional.

### Running the Pipeline

Just run the following command from the build folder:

```
$ ./standalone path/to/config.yaml
```

To run the pipeline on system start, check your specific system init documentation.

## Lambda

The Lambda environment runs the data pipeline in AWS Lambda instances. It's located in the `lambda` folder, and is divided into 3 binaries: `lambda/receiver`, `lambda/processor` and `lambda/exporter`.

### Prerequisites

- An AWS account.
- [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html) tool.
- Go 1.20 or later.
- GNU Make.

### Setting Up AWS

1. Create 3 lambdas for **Receiver**, **Processor** and **Exporter**. Runtime `Go 1.x`, arch `x86_64`, and handler names `receiver`, `processor` and `exporter`.
2. Create an SQS for **ProcessorToExporter**, type Standard, condition OnSuccess.
3. Open **Receiver** lambda config->permissions->execution role. Add another permission->create inline policy, add Lambda write permissions `InvokeAsync` and `InvokeFunction`.
4. Edit **Receiver** lambda config, add as a destination another lambda, the **Processor**, with async invocation.
5. Open **Processor** lambda config->permissions->execution role. Add another permission->create inline policy, add SQS write permissions.
6. Edit **Processor** lambda config, add as a destination the SQS **ProcessorToExporter**.
7. Open **Exporter** lambda config->permissions->execution role. Add another permission->create inline policy, add SQS read and write permissions.
8. Edit **Exporter** lambda config, add as a trigger the SQS **ProcessorToExporter**.

Note: when creating and configuring the SQS service and trigger, make sure to set the timing and batching options you will need. For example, a time interval of 5 minutes and batching of 50 events.

### Build & Deploy

Open a terminal, CD to `lambda`, and run:

```
$ make recv=RECEIVER proc=PROCESSOR expt=EXPORTER
```

Where *RECEIVER*, *PROCESSOR*, and *EXPORTER*, are the AWS Lambda functions you just created in the previous step.

### Configuring the Pipeline

A Lambda pipeline requieres some configuration keys to be set as **environment variables**. To set up environment variables, go to AWS console, Lambda->Functions, click your function, Configuration->Environment variables. Currently, only the Exporter lambda requires some variables to be set:

- `exporter`: String. Exporter type, possible values are: `nrapi`, `nrmetrics`, `nrevents`, `nrlogs`, `nrtraces`, `otel` and `prom`.

Each one of the exporters also require specific configuration.

#### New Relic APIs exporter

For exporter types `nrapi`, `nrmetrics`, `nrevents`, `nrlogs`, and `nrtraces`, the following keys are used:

- `nr_account_id`: String. Account ID. Only requiered for `nrevents` and `nrapi` exporters.
- `nr_api_key`: String. Api key for writing.
- `nr_endpoint`: String. New Relic endpoint region. Either `US` or `EU`. Optional, default value is `US`.

#### OpenTelemetry exporter

For exporter type `otel` the following keys are used:

- `otel_endpoint`: String. Domain and (optionally) port only, no scheme or path. Example: `my.otel.example.com:4318`
- `otel_scheme`: String. Either `http` or `https`. If not specified, default value is `https`.
- `otel_headers`: Map. List of key-value pairs sent as HTTP request headers. Optional.
- `otel_metric_endpoint`: String. Full URL of an endpoint to send metrics. If specified, `otel_endpoint` and `otel_scheme` will be ignored.
- `otel_logs_endpoint`: String. Full URL of an endpoint to send logs. If specified, `otel_endpoint` and `otel_scheme` will be ignored.

#### Prometheus exporter

For exporter type `prom` the following keys are used:

- `prom_endpoint`: String. Prometheus remote write URL.
- `prom_credentials`: String. Optional authorization credentials.
- `prom_headers`: Map. List of key-value pairs sent as HTTP request headers. Optional.

### Running the Pipeline

Finally, to start running the pipeline you will need an EventBridge rule. Add a trigger for the **Receiver** lambda, select EventBridge as the source, create new rule, schedule expression `rate(1 minute)` (or the time you desire).

### Testing

Instead of running the pipeline with an EventBridge rule, you can just send async invocations to the **Receiver** lambda from the command line, using the following command:

```
$ aws lambda invoke-async --function-name RECEIVER --invoke-args INPUT.json
```

Where *RECEIVER* is the **Receiver** lambda name and *INPUT.json* is a file containing any JSON (the input event will be ignored by the receiver).

This will simulate a timer event and trigger the pipeline.

## Infra

The Infra environment runs the data pipeline as a New Relic Infrastructure Agent plugin. It's located in the `cmd/infra` folder.

### Prerequisites

- [New Relic Infrastructure Agent](https://docs.newrelic.com/docs/infrastructure/install-infrastructure-agent/get-started/install-infrastructure-agent/).
- Go 1.20 or later.

### Build & Desploy

Open a terminal, CD to `cmd/infra`, and run:

```
$ go build
```

It will generate a binary named `infra` in the same folder. Move it to the folder where you have all your NR Infra integrations, for Linux it is `/var/db/newrelic-infra/custom-integrations/`, and rename it as you wish, usually `nri-SOMETHING`.

### Configuring the Pipeline

As any other NR Infra integration, it's configured through a YAML file in `/etc/newrelic-infra/integrations.d/` (Linux), that should look like:

```
integrations:
  - name: nri-myintegration    
    config:
      key_one: value_one
      key_two: value_two
    interval: 15s
```

Whatever you put under `config` will be passed to the integration. The following keys are used:

- `name`: Integration name. Default value `InfraIntegration`.
- `version`: Integration version. Default value `0.1.0`.
- `entity_name`: Entity name. Must be unique for the NR account, otherwise it will cause conflicts. Default value `EntityName`.
- `entity_type`: Entity type. Default value `EntityType`.
- `entity_display`: Entity display name. Default value `EntityDisplay`.

Unlike Standalone and Lambda environments, there is no need to set the exporter, because it will be always `nrinfra`.

For further information on NR Infra configuration, check out the following links:

- https://docs.newrelic.com/docs/infrastructure/install-infrastructure-agent/configuration/configure-infrastructure-agent/
- https://docs.newrelic.com/docs/infrastructure/install-infrastructure-agent/configuration/infrastructure-agent-configuration-settings/
- https://github.com/newrelic/infrastructure-agent/blob/master/assets/examples/infrastructure/newrelic-infra-template.yml.example

### Running the Pipeline

It will be automatically executed by the New Relic Infrastructure Agent based on the time interval configuration.

### Testing

Instead of running it automatically, you can run it manually for testing and debugging:

```
$ go build && CONFIG_PATH=/path/to/config.yaml ./infra
```

Where `/path/to/config.yaml` is a YAML file containing the contents of your `config` key.
