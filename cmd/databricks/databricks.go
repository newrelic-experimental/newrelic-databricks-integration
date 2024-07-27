package main

import (
	"context"
	"fmt"

	"github.com/newrelic-experimental/newrelic-databricks-integration/internal/databricks"
	"github.com/newrelic/newrelic-labs-sdk/pkg/integration"
	"github.com/newrelic/newrelic-labs-sdk/pkg/integration/log"
	"github.com/spf13/viper"
)

var (
	/* Args below are populated via ldflags at build time */
	gIntegrationID      = "com.newrelic.labs.newrelic-databricks-integration"
	gIntegrationName    = "New Relic Databricks Integration"
	gIntegrationVersion = "2.0.0"
	gGitCommit          = ""
	gBuildDate          = ""
	gBuildInfo			= integration.BuildInfo{
		Id:        gIntegrationID,
		Name:      gIntegrationName,
		Version:   gIntegrationVersion,
		GitCommit: gGitCommit,
		BuildDate: gBuildDate,
	}
)

func main() {
	// Create a new background context to use
	ctx := context.Background()

	// Create the integration with options
	i, err := integration.NewStandaloneIntegration(
		&gBuildInfo,
		gBuildInfo.Name,
		integration.WithInterval(60),
		integration.WithLicenseKey(),
		integration.WithApiKey(),
		integration.WithEvents(ctx),
	)
	fatalIfErr(err)

	mode := viper.GetString("mode")
	if mode == "" {
		mode = "databricks"
	}

	switch mode {
	case "databricks":
		err = databricks.InitPipelines(ctx, i)
		fatalIfErr(err)

	// @TODO: support any spark context
	//case "spark":
	//	err = spark.InitPipelines(i)
	//	fatalIfErr(err)

	// @TODO: support other cluster providers/modes like yarn/k8s

	default:
		fatalIfErr(fmt.Errorf("unrecognized mode %s", mode))
	}

	// Run the integration
	defer i.Shutdown(ctx)
 	err = i.Run(ctx)
	fatalIfErr(err)
}

func fatalIfErr(err error) {
	if err != nil {
		log.Fatalf(err)
	}
}
