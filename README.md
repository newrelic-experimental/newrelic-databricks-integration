# New Relic Databricks Integration

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

Check `config/example_config.yaml` for a configuration example.


#### New Relic APIs exporter

- `nr_account_id`: String. Account ID. Only requiered for `nrevents` and `nrapi` exporters.
- `nr_api_key`: String. Api key for writing.
- `nr_endpoint`: String. New Relic endpoint region. Either `US` or `EU`. Optional, default value is `US`.

### Running the Pipeline

Just run the following command from the build folder:

```
$ ./standalone path/to/config.yaml
```

To run the pipeline on system start, check your specific system init documentation.

