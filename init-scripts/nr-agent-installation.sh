#!/bin/bash

# Installs New Relic Java APM Agent

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