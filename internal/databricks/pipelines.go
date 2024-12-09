package databricks

import (
	"context"
	"fmt"
	"time"

	databricksSdk "github.com/databricks/databricks-sdk-go"
	databricksSdkPipelines "github.com/databricks/databricks-sdk-go/service/pipelines"
	"github.com/newrelic/newrelic-labs-sdk/v2/pkg/integration"
	"github.com/newrelic/newrelic-labs-sdk/v2/pkg/integration/log"
	"github.com/newrelic/newrelic-labs-sdk/v2/pkg/integration/model"
)

const (
	RFC_3339_MILLI_LAYOUT = "2006-01-02T15:04:05.000Z07:00"
)

type DatabricksPipelineEventsReceiver struct {
	i							*integration.LabsIntegration
	w							*databricksSdk.WorkspaceClient
	tags 						map[string]string
}

func NewDatabricksPipelineEventsReceiver(
	i *integration.LabsIntegration,
	w *databricksSdk.WorkspaceClient,
	tags map[string]string,
) *DatabricksPipelineEventsReceiver {
	return &DatabricksPipelineEventsReceiver{
		i,
		w,
		tags,
	}
}

func (d *DatabricksPipelineEventsReceiver) GetId() string {
	return "databricks-pipeline-events-receiver"
}

func (d *DatabricksPipelineEventsReceiver) PollLogs(
	ctx context.Context,
	writer chan <- model.Log,
) error {
	lastRun := time.Now().Add(-d.i.Interval * time.Second)

	all := d.w.Pipelines.ListPipelines(
		ctx,
		databricksSdkPipelines.ListPipelinesRequest{},
	)

	LOOP:

	for ; all.HasNext(ctx);  {
		pipelineStateInfo, err := all.Next(ctx)
		if err != nil {
			return err
		}

		log.Debugf(
			"processing pipeline events for pipeline %s (%s) with state %s",
			pipelineStateInfo.PipelineId,
			pipelineStateInfo.Name,
			pipelineStateInfo.State,
		)

		allEvents := d.w.Pipelines.ListPipelineEvents(
			ctx,
			databricksSdkPipelines.ListPipelineEventsRequest{
				Filter: "timestamp >= '" +
					lastRun.UTC().Format(RFC_3339_MILLI_LAYOUT) +
					"'",
				PipelineId: pipelineStateInfo.PipelineId,
			},
		)

		count := 0

		for ; allEvents.HasNext(ctx);  {
			pipelineEvent, err := allEvents.Next(ctx)
			if err != nil {
				log.Warnf(
					"unexpected no more items error or failed to fetch next page of pipeline events for pipeline %s (%s): %v",
					pipelineStateInfo.PipelineId,
					pipelineStateInfo.Name,
					err,
				)
				continue LOOP
			}

			// eventTimestamp is not used until the end of the loop but by
			// parsing here we can abort if an error occurs and skip allocating
			// and populating the attributes map.

			eventTimestamp, err := time.Parse(
				RFC_3339_MILLI_LAYOUT,
				pipelineEvent.Timestamp,
			)
			if err != nil {
				log.Warnf(
					"ignoring event with ID %s with invalid timestamp %s for pipeline %s (%s): %v",
					pipelineEvent.Id,
					pipelineEvent.Timestamp,
					pipelineStateInfo.PipelineId,
					pipelineStateInfo.Name,
					err,
				)
				continue
			}

			attrs := makeAttributesMap(d.tags)

			attrs["databricksPipelineEventId"] = pipelineEvent.Id
			attrs["databricksPipelineEventType"] = pipelineEvent.EventType
			attrs["level"] = pipelineEvent.Level
			attrs["databricksPipelineEventLevel"] = pipelineEvent.Level
			attrs["databricksPipelineEventMaturityLevel"] =
				pipelineEvent.MaturityLevel

			isError := pipelineEvent.Error != nil

			attrs["databricksPipelineEventError"] = isError

			if isError {
				attrs["databricksPipelineEventErrorFatal"] =
					pipelineEvent.Error.Fatal

				if len(pipelineEvent.Error.Exceptions) > 0 {
					for i, e := range pipelineEvent.Error.Exceptions {
						attrName := fmt.Sprintf(
							"databricksPipelineEventErrorException%dClassName",
							i + 1,
						)
						attrs[attrName] = e.ClassName

						attrName = fmt.Sprintf(
							"databricksPipelineEventErrorException%dMessage",
							i + 1,
						)
						attrs[attrName] = e.Message
					}
				}
			}

			if pipelineEvent.Origin != nil {
				origin := pipelineEvent.Origin

				attrs["databricksPipelineEventBatchId"] = origin.BatchId
				attrs["databricksPipelineEventCloud"] = origin.Cloud
				attrs["databricksPipelineEventClusterId"] = origin.ClusterId
				attrs["databricksPipelineEventDatasetName"] = origin.DatasetName
				attrs["databricksPipelineEventFlowId"] = origin.FlowId
				attrs["databricksPipelineEventFlowName"] = origin.FlowName
				attrs["databricksPipelineEventHost"] = origin.Host
				attrs["databricksPipelineEventMaintenanceId"] =
					origin.MaintenanceId
				attrs["databricksPipelineEventMaterializationName"] =
					origin.MaterializationName
				attrs["databricksPipelineEventOrgId"] = origin.OrgId
				attrs["databricksPipelineEventPipelineId"] = origin.PipelineId
				attrs["databricksPipelineEventPipelineName"] =
					origin.PipelineName
				attrs["databricksPipelineEventRegion"] = origin.Region
				attrs["databricksPipelineEventRequestId"] = origin.RequestId
				attrs["databricksPipelineEventTableId"] = origin.TableId
				attrs["databricksPipelineEventUcResourceId"] =
					origin.UcResourceId
				attrs["databricksPipelineEventUpdateId"] = origin.UpdateId
			}

			writer <- model.NewLog(
				pipelineEvent.Message,
				attrs,
				eventTimestamp,
			)

			count += 1
		}

		log.Debugf("received %d pipeline events", count)
	}

	return nil
}
