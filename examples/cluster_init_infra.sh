#!/bin/bash

# Define the config path

NEW_RELIC_CONFIG_FILE="/etc/newrelic-infra.yml"

# Add New Relic's Infrastructure Agent gpg key
curl -fsSL https://download.newrelic.com/infrastructure_agent/gpg/newrelic-infra.gpg | sudo gpg --dearmor -o /etc/apt/trusted.gpg.d/newrelic-infra.gpg

# Create a manifest file
echo "deb https://download.newrelic.com/infrastructure_agent/linux/apt/ jammy main" | sudo tee -a /etc/apt/sources.list.d/newrelic-infra.list

# Run an update
sudo apt-get update

# Install the New Relic Infrastructure agent
sudo apt-get install newrelic-infra -y

# Create the agent config file
sudo cat <<EOM >> $NEW_RELIC_CONFIG_FILE
license_key: $NEW_RELIC_LICENSE_KEY
custom_attributes:
  databricksClusterId: $DB_CLUSTER_ID
  databricksClusterName: $DB_CLUSTER_NAME
  databricksIsDriverNode: $DB_IS_DRIVER
  databricksIsJobCluster: $DB_IS_JOB_CLUSTER
EOM

# Start the agent
sudo systemctl start newrelic-infra
