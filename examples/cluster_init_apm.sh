#!/bin/bash

# Define the newrelic version, jar path, and config path

NEW_RELIC_VERSION="current" # update if you want to use a specific version
NEW_RELIC_JAR_PATH="/databricks/jars/newrelic-agent.jar"
NEW_RELIC_CONFIG_FILE="/databricks/jars/newrelic.yml"

# Download the newrelic java agent
curl -o $NEW_RELIC_JAR_PATH -L https://download.newrelic.com/newrelic/java-agent/newrelic-agent/$NEW_RELIC_VERSION/newrelic-agent.jar

NEW_RELIC_APP_NAME="Databricks"

if [ "$DB_IS_DRIVER" = "TRUE" ]; then
  NEW_RELIC_APP_NAME+=" Driver"
else
  NEW_RELIC_APP_NAME+=" Executor"
fi

# Create agent config file
sudo cat <<EOM > $NEW_RELIC_CONFIG_FILE
common: &default_settings
  app_name: $NEW_RELIC_APP_NAME
  license_key: $NEW_RELIC_LICENSE_KEY
  agent_enabled: true
  labels: "databricksClusterId:$DB_CLUSTER_ID;databricksClusterName:$DB_CLUSTER_NAME;databricksIsDriverNode:$DB_IS_DRIVER;databricksIsJobCluster:$DB_IS_JOB_CLUSTER"
production:
  <<: *default_settings
EOM
