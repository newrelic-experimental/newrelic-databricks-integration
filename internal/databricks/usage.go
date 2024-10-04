package databricks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"time"

	databricksSdk "github.com/databricks/databricks-sdk-go"
	databricksSdkCompute "github.com/databricks/databricks-sdk-go/service/compute"
	databricksSql "github.com/databricks/databricks-sdk-go/service/sql"
	"github.com/spf13/viper"

	"github.com/newrelic/newrelic-client-go/v2/pkg/region"
	"github.com/newrelic/newrelic-labs-sdk/v2/pkg/integration"
	"github.com/newrelic/newrelic-labs-sdk/v2/pkg/integration/connectors"
	"github.com/newrelic/newrelic-labs-sdk/v2/pkg/integration/log"
	"github.com/newrelic/newrelic-labs-sdk/v2/pkg/integration/model"
)

const (
	DATABRICKS_USAGE_EVENT_TYPE = "DatabricksUsage"
	DATABRICKS_LIST_PRICES_LOOKUP_TABLE_NAME = "DatabricksListPrices"
)

type newRelicNrqlLookupsApiRequest struct {
	Table struct {
		Headers []string 				`json:"headers"`
		Rows [][]interface{}			`json:"rows"`
	} 									`json:"table"`
}

type clusterInfo struct {
	name			string
	source			string
	creator			string
	instancePoolId	string
	singleUserName	string
}

type warehouseInfo struct {
	name			string
	creator			string
}

type DatabricksUsageReceiver struct {
	w					*databricksSdk.WorkspaceClient
	a					*databricksSdk.AccountClient
	warehouseId 		string
}

func NewDatabricksUsageReceiver(
	w *databricksSdk.WorkspaceClient,
	a *databricksSdk.AccountClient,
	warehouseId string,
) *DatabricksUsageReceiver {
	return &DatabricksUsageReceiver{
		w,
		a,
		warehouseId,
	}
}

func (d *DatabricksUsageReceiver) GetId() string {
	return "databricks-usage-receiver"
}

func (d *DatabricksUsageReceiver) PollEvents(
	ctx context.Context,
	writer chan <- model.Event,
) error {
	// retrieve results for previous day by setting usageDate to now - 24 hours
	year, month, day := time.Now().Add(-24 * time.Hour).Date()
	date := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

	log.Debugf("collecting billing usage data")
	err := d.importBillingUsageTable(ctx, date, writer)
	if err != nil {
		return err
	}

	return nil
}

func NewDatabricksListPricesCollector(
	li *integration.LabsIntegration,
	w *databricksSdk.WorkspaceClient,
	warehouseId string,
) *integration.SimpleComponent {
	return integration.NewSimpleComponent(
		"databricks-list-prices-collector",
		li,
		func(ctx context.Context, i *integration.LabsIntegration) error {
			return li.UseLastUpdate(ctx, func(lastUpdate time.Time) error {
				log.Debugf("collecting billing list prices data")
				err := importBillingListPricesTable(
					ctx,
					w,
					warehouseId,
					i.GetRegion(),
					i.GetAccountId(),
					i.GetApiKey(),
					i.DryRun,
					lastUpdate,
				)
				if err != nil {
					return err
				}

				return nil
			}, "databricks-list-prices-collector")
		},
	)
}

func (d *DatabricksUsageReceiver) importBillingUsageTable(
	ctx context.Context,
	usageDate time.Time,
	writer chan <- model.Event,
) error {
	workspaceNameById, err := buildWorkspaceNameByIdMap(ctx, d.a)
	if err != nil {
		return err
	}

	clusterInfoById, err := buildClusterInfoByIdMap(ctx, d.a)
	if err != nil {
		return err
	}

	warehouseInfoById, err := buildWarehouseInfoByIdMap(ctx, d.a)
	if err != nil {
		return err
	}

	stmt := fmt.Sprintf(`
SELECT
	t1.account_id AS account_id, t1.workspace_id AS workspace_id,
	record_id, sku_name, cloud,
	usage_start_time, usage_end_time, usage_date, custom_tags,
	usage_unit, usage_quantity,	record_type, ingestion_date,
	billing_origin_product, usage_type,
	usage_metadata.cluster_id,
	usage_metadata.warehouse_id, usage_metadata.instance_pool_id,
	usage_metadata.node_type, usage_metadata.job_id,
	usage_metadata.job_run_id, usage_metadata.job_name,
	usage_metadata.notebook_id, usage_metadata.notebook_path,
	usage_metadata.dlt_pipeline_id, usage_metadata.dlt_update_id,
	usage_metadata.dlt_maintenance_id, usage_metadata.run_name,
	usage_metadata.endpoint_name, usage_metadata.endpoint_id,
	usage_metadata.central_clean_room_id, identity_metadata.run_as,
	product_features.jobs_tier, product_features.sql_tier,
	product_features.dlt_tier, product_features.is_serverless,
	product_features.is_photon, product_features.serving_type
FROM
	system.billing.usage t1
WHERE t1.usage_date = '%s'
LIMIT 10000`,
		usageDate.Format(time.DateOnly),
	)

	rows, err := executeStatementOnWarehouse(
		ctx,
		d.w,
		d.warehouseId,
		stmt,
	)
	if err != nil {
		return err
	}

	includeIdentityMetadata := viper.GetBool(
		"databricks.usage.includeIdentityMetadata",
	)

	for _, row := range rows {
		attrs := map[string]interface{}{}

		attrs["account_id"] = row[0]
		attrs["workspace_id"] = row[1]
		attrs["record_id"] = row[2]
		attrs["sku_name"] = row[3]
		attrs["cloud"] = row[4]

		start, err := time.Parse(time.RFC3339Nano, row[5])
		if err != nil {
			log.Warnf("invalid start date \"%s\" for record %s", row[5], row[2])
			continue
		}

		end, err := time.Parse(time.RFC3339Nano, row[6])
		if err != nil {
			log.Warnf("invalid end date \"%s\" for record %s", row[6], row[2])
			continue
		}

		attrs["usage_start_date"] = start.UnixMilli()
		attrs["usage_end_date"] = end.UnixMilli()
		attrs["usage_date"] = row[7]

		// custom_tags
		tags, err := unmarshalStruct(row[8])
		if err != nil {
			log.Warnf("invalid json for custom tags for record %s", row[2])
			continue
		}

		for k, v := range tags {
			attrs[k] = v
		}

		attrs["usage_unit"] = row[9]

		quantity, err := strconv.ParseFloat(row[10], 64)
		if err != nil {
			log.Warnf("invalid quantity \"%s\" for record %s", row[10], row[2])
			continue
		}

		attrs["usage_quantity"] = quantity
		attrs["record_type"] = row[11]
		attrs["ingestion_date"] = row[12]
		attrs["billing_origin_product"] = row[13]
		attrs["usage_type"] = row[14]

		// usage_metadata
		attrs["cluster_id"] = row[15]

		if row[15] != "" {
			clusterInfo, err := getClusterInfoById(
				row[1],
				row[15],
				clusterInfoById,
			)
			if err != nil {
				log.Warnf(
					"invalid or missing cluster id %v found, cluster info will be empty: %v",
					attrs["cluster_id"],
					err,
				)
			} else if clusterInfo != nil {
				attrs["cluster_name"] = clusterInfo.name

				if includeIdentityMetadata {
					attrs["cluster_single_user_name"] = clusterInfo.singleUserName
					attrs["cluster_creator"] = clusterInfo.creator
				}

				attrs["cluster_source"] = clusterInfo.source
				attrs["cluster_instance_pool_id"] = clusterInfo.instancePoolId
			}
		}

		attrs["warehouse_id"] = row[16]

		if row[16] != "" {
			warehouseInfo, err := getWarehouseInfoById(
				row[1],
				row[16],
				warehouseInfoById,
			)
			if err != nil {
				log.Warnf(
					"invalid or missing warehouse id %v found, warehouse info will be empty: %v",
					attrs["warehouse_id"],
					err,
				)
			} else if warehouseInfo != nil {
				attrs["warehouse_name"] = warehouseInfo.name

				if includeIdentityMetadata {
					attrs["warehouse_creator"] = warehouseInfo.creator
				}
			}
		}

		attrs["instance_pool_id"] = row[17]
		attrs["node_type"] = row[18]
		attrs["job_id"] = row[19]
		attrs["job_run_id"] = row[20]
		attrs["job_name"] = row[21]
		attrs["notebook_id"] = row[22]
		attrs["notebook_path"] = row[23]
		attrs["dlt_pipeline_id"] = row[24]
		attrs["dlt_update_id"] = row[25]
		attrs["dlt_maintenance_id"] = row[26]
		attrs["run_name"] = row[27]
		attrs["endpoint_name"] = row[28]
		attrs["endpoint_id"] = row[29]
		attrs["central_clean_room_id"] = row[30]

		// identity_metadata
		if includeIdentityMetadata {
			attrs["run_as"] = row[31]
		}

		// product features
		attrs["jobs_tier"] = row[32]
		attrs["sql_tier"] = row[33]
		attrs["dlt_tier"] = row[34]
		attrs["is_serverless"] = row[35]
		attrs["is_photon"] = row[36]
		attrs["serving_type"] = row[37]

		workspaceName, err := getWorkspaceNameById(
			row[1],
			workspaceNameById,
		)
		if err != nil {
			log.Warnf(
				"invalid or missing workspace id %v found, workspace_name will be empty: %v",
				attrs["workspace_id"],
				err,
			)
		} else if workspaceName != "" {
			attrs["workspace_name"] = workspaceName
		}

		// write the event to the event channel with a timestamp of 23 hours ago
		// so that it occurs around the start time of the data we queried from
		// the system.billing.usage table.
		writer <- model.NewEvent(
			DATABRICKS_USAGE_EVENT_TYPE,
			attrs,
			time.Now().Add(-23 * time.Hour),
		)
	}

	return nil
}

func buildListPricesNrqlLookupsApiRequest(
	rows [][]string,
) *newRelicNrqlLookupsApiRequest {
	body := &newRelicNrqlLookupsApiRequest{}

	body.Table.Headers = []string{
		"account_id",
		"price_start_time",
		"price_end_time",
		"sku_name",
		"cloud",
		"currency_code",
		"usage_unit",
		"list_price",
		"promotional_price",
		"effective_list_price",
	}

	tableRows := [][]interface{}{}

	for _, row := range rows {
		// price_start_time is always populated and is the time the price
		// went into effect
		start, err := time.Parse(time.RFC3339Nano, row[1])
		if err != nil {
			log.Warnf(
				"invalid price start date \"%s\" for sku %s",
				row[1],
				row[3],
			)
			continue
		}

		// price_end_time can sometimes be null if the price is still in
		// effect
		end := time.Time{}
		if row[2] != "" {
			end, err = time.Parse(time.RFC3339Nano, row[2])
			if err != nil {
				log.Warnf(
					"invalid price end date \"%s\" for sku %s",
					row[2],
					row[3],
				)
				continue
			}
		}

		list, promo, effectiveList, err := unmarshalPricingStruct(row[7])
		if err != nil {
			log.Warnf("invalid json for pricing data for sku %s", row[3])
			continue
		}

		endMilli := int64(0)
		if !end.IsZero() {
			endMilli = end.UnixMilli()
		}

		tableRows = append(tableRows, []interface{}{
			row[0],
			start.UnixMilli(),
			endMilli,
			row[3],
			row[4],
			row[5],
			row[6],
			list,
			promo,
			effectiveList,
		})
	}

	body.Table.Rows = tableRows

	return body
}

func listPricesPostBodyBuilder(rows [][]string) connectors.HttpBodyBuilder {
	return func() (any, error) {
		jsonValue, err := json.Marshal(
			buildListPricesNrqlLookupsApiRequest(rows),
		)
		if err != nil {
			return nil, err
		}

		return bytes.NewReader(jsonValue), nil
	}
}

func importBillingListPricesTable(
	ctx context.Context,
	w *databricksSdk.WorkspaceClient,
	warehouseId string,
	newRelicRegion region.Name,
	newRelicAccountId int,
	newRelicApiKey string,
	dryRun bool,
	lastUpdate time.Time,
) error {
	stmt := `
SELECT
	account_id, price_start_time, price_end_time, sku_name,
	cloud, currency_code, usage_unit, pricing
FROM
	system.billing.list_prices
LIMIT 15000`

	rows, err := executeStatementOnWarehouse(
		ctx,
		w,
		warehouseId,
		stmt,
	)
	if err != nil {
		return err
	}

	client := connectors.NewHttpPostConnector(
		getNrqlLookupsApiUrl(
			newRelicRegion,
			newRelicAccountId,
		),
		listPricesPostBodyBuilder(rows),
	)

	// On original creation, method must be POST but on updates, must be PUT
	// so if the last update time is not the zero time, we have already created
	// the lookup table so use PUT. Otherwise use POST.
	if !lastUpdate.IsZero() {
		client.SetMethod("PUT")
	}

	if dryRun || log.IsDebugEnabled() {
		log.Debugf("list prices lookup table payload JSON follows")
		log.PrettyPrintJson(buildListPricesNrqlLookupsApiRequest(rows))

		if dryRun {
			return nil
		}
	}

	headers := map[string]string{
		"Accept": "application/json",
		"Content-Type": "application/json",
		"Api-Key": newRelicApiKey,
	}

	client.SetHeaders(headers)

	readCloser, err := client.Request()
	if err != nil {
		return err
	}

	defer readCloser.Close()

	if log.IsDebugEnabled() {
		log.Debugf("table api response payload follows")

		bytes, err := io.ReadAll(readCloser)
		if err != nil {
			log.Warnf("error reading table api response: %v", err)
		} else {
			log.Debugf(string(bytes))
		}
	}

	return nil
}

func unmarshalStruct(data string) (map[string]interface{}, error) {
	var (
		m map[string]interface{}
	)

	attrs := map[string]interface{}{}

	err := json.Unmarshal([]byte(data), &m)
	if err != nil {
		return nil, err
	}

	for k, v := range m {
		switch val := v.(type) {
		case string, float64, bool, nil:
			attrs[k] = val
		default:
			log.Warnf(
				"skipping unsupported value type for key %s when converting struct to event attributes",
				k,
			)
		}
	}

	return attrs, nil
}

func unmarshalPricingStruct(data string) (float64, float64, float64, error) {
	var (
		m map[string]interface{}
	)

	err := json.Unmarshal([]byte(data), &m)
	if err != nil {
		return 0.0, 0.0, 0.0, err
	}

	listPrice := safeGetFloatField(m, "default", 0.0)
	promoPrice := 0.0
	effectiveListPrice := 0.0

	p, ok := m["promotional"]
	if ok {
		promoMap, ok := p.(map[string]interface{})
		if ok {
			promoPrice = safeGetFloatField(promoMap, "default", 0.0)
		}
	}

	e, ok := m["effective_list"]
	if ok {
		elMap, ok := e.(map[string]interface{})
		if ok {
			effectiveListPrice = safeGetFloatField(elMap, "default", 0.0)
		}
	}

	return listPrice, promoPrice, effectiveListPrice, nil
}

func safeGetFloatField(
	m map[string]interface{},
	name string,
	def float64,
) float64 {
	v, ok := m[name]
	if !ok {
		return def
	}

	vs, ok := v.(string)
	if !ok {
		return def
	}

	f, err := strconv.ParseFloat(vs, 64)
	if err != nil {
		return def
	}

	return f
}

func getNrqlLookupsApiUrl(
	newRelicRegion region.Name,
	newRelicAccountId int,
) string {
	if newRelicRegion == region.EU {
		return fmt.Sprintf(
			"https://nrql-lookup.service.eu.newrelic.com/v1/accounts/%d/%s",
			newRelicAccountId,
			DATABRICKS_LIST_PRICES_LOOKUP_TABLE_NAME,
		)
	}

	return fmt.Sprintf(
		"https://nrql-lookup.service.newrelic.com/v1/accounts/%d/%s",
		newRelicAccountId,
		DATABRICKS_LIST_PRICES_LOOKUP_TABLE_NAME,
	)
}

func buildWorkspaceNameByIdMap(
	ctx context.Context,
	a *databricksSdk.AccountClient,
) (map[int64]string, error) {
	log.Debugf("building workspace name by ID map...")

	workspaces, err := a.Workspaces.List(ctx)
	if err != nil {
		return nil, err
	}

	m := make(map[int64]string)

	for _, workspace := range workspaces {
		m[workspace.WorkspaceId] = workspace.WorkspaceName
	}

	return m, nil
}

func getWorkspaceNameById(
	workspaceIdStr string,
	workspaceIdToName map[int64]string,
) (string, error) {
	workspaceName := ""

	workspaceId, err := strconv.Atoi(workspaceIdStr)
	if err != nil {
		return "", fmt.Errorf("workspace ID is not an integer")
	}

	name, ok := workspaceIdToName[int64(workspaceId)]
	if ok {
		workspaceName = name
	}

	return workspaceName, nil
}

func buildClusterInfoByIdMap(
	ctx context.Context,
	a *databricksSdk.AccountClient,
) (map[string]*clusterInfo, error) {
	log.Debugf("building cluster ID to info map...")

	workspaces, err := a.Workspaces.List(ctx)
	if err != nil {
		return nil, err
	}

	m := map[string]*clusterInfo{}

	for _, workspace := range workspaces {
		w, err := a.GetWorkspaceClient(workspace)
		if err != nil {
			return nil, err
		}

		log.Debugf("listing clusters for workspace %s", workspace.WorkspaceName)

		all := w.Clusters.List(
			ctx,
			databricksSdkCompute.ListClustersRequest{ PageSize: 100 },
		)

		for ; all.HasNext(ctx);  {
			c, err := all.Next(ctx)
			if err != nil {
				return nil, err
			}

			log.Debugf(
				"cluster ID: %s ; cluster name: %s",
				c.ClusterId,
				c.ClusterName,
			)

			// namespace cluster ids with workspace id, just in case cluster ids
			// can be the same across workspaces
			id := fmt.Sprintf(
				"%d.%s",
				workspace.WorkspaceId,
				c.ClusterId,
			)

			clusterInfo := &clusterInfo{}
			clusterInfo.name = c.ClusterName
			clusterInfo.source = string(c.ClusterSource)
			clusterInfo.creator = c.CreatorUserName
			clusterInfo.singleUserName = c.SingleUserName
			clusterInfo.instancePoolId = c.InstancePoolId

			m[id] = clusterInfo
		}
	}

	return m, nil
}

func getClusterInfoById(
	workspaceIdStr string,
	clusterIdStr string,
	clusterInfoById map[string]*clusterInfo,
) (*clusterInfo, error) {
	workspaceId, err := strconv.Atoi(workspaceIdStr)
	if err != nil {
		return nil, fmt.Errorf("workspace ID is not an integer")
	}

	// namespace cluster ids with workspace id, just in case cluster ids can be
	// the same across workspaces
	id := fmt.Sprintf(
		"%d.%s",
		workspaceId,
		clusterIdStr,
	)

	clusterInfo, ok := clusterInfoById[id]
	if ok {
		return clusterInfo, nil
	}

	return nil, nil
}

func buildWarehouseInfoByIdMap(
	ctx context.Context,
	a *databricksSdk.AccountClient,
) (map[string]*warehouseInfo, error) {
	log.Debugf("building warehouse info by ID map...")

	workspaces, err := a.Workspaces.List(ctx)
	if err != nil {
		return nil, err
	}

	m := map[string]*warehouseInfo{}

	for _, workspace := range workspaces {
		w, err := a.GetWorkspaceClient(workspace)
		if err != nil {
			return nil, err
		}

		log.Debugf(
			"listing warehouses for workspace %s",
			workspace.WorkspaceName,
		)

		all := w.Warehouses.List(
			ctx,
			databricksSql.ListWarehousesRequest{},
		)

		for ; all.HasNext(ctx); {
			warehouse, err := all.Next(ctx)
			if err != nil {
				return nil, err
			}

			log.Debugf(
				"warehouse ID: %s ; warehouse name: %s",
				warehouse.Id,
				warehouse.Name,
			)

			// namespace warehouse ids with workspace id, just in case warehouse
			// ids can be the same across workspaces
			id := fmt.Sprintf(
				"%d.%s",
				workspace.WorkspaceId,
				warehouse.Id,
			)

			warehouseInfo := &warehouseInfo{}
			warehouseInfo.name = warehouse.Name
			warehouseInfo.creator = warehouse.CreatorName

			m[id] = warehouseInfo
		}
	}

	return m, nil
}

func getWarehouseInfoById(
	workspaceIdStr string,
	warehouseIdStr string,
	warehouseInfoById map[string]*warehouseInfo,
) (*warehouseInfo, error) {
	workspaceId, err := strconv.Atoi(workspaceIdStr)
	if err != nil {
		return nil, fmt.Errorf("workspace ID is not an integer")
	}

	// namespace warehouse ids with workspace id, just in case warehouse ids can
	// be the same across workspaces
	id := fmt.Sprintf(
		"%d.%s",
		workspaceId,
		warehouseIdStr,
	)

	warehouseInfo, ok := warehouseInfoById[id]
	if ok {
		return warehouseInfo, nil
	}

	return nil, nil
}
