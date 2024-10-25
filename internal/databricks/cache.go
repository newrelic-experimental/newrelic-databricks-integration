package databricks

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	databricksSdk "github.com/databricks/databricks-sdk-go"
	databricksSdkCompute "github.com/databricks/databricks-sdk-go/service/compute"
	databricksSql "github.com/databricks/databricks-sdk-go/service/sql"
	"github.com/newrelic/newrelic-labs-sdk/v2/pkg/integration/log"
)

type cacheLoaderFunc[T interface{}] func(ctx context.Context) (*T, error)

type memoryCache[T interface{}] struct {
	mu				sync.Mutex
	expiry			time.Duration
	value			*T
	expiration		time.Time
	loader			cacheLoaderFunc[T]
}

type workspaceInfo struct {
	name			string
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

var (
	workspaceInfoCache			*memoryCache[map[int64]*workspaceInfo]
	clusterInfoCache			*memoryCache[map[string]*clusterInfo]
	warehouseInfoCache			*memoryCache[map[string]*warehouseInfo]
)

func newMemoryCache[T interface{}](
	expiry time.Duration,
	loader cacheLoaderFunc[T],
) *memoryCache[T] {
	return &memoryCache[T] {
		expiry: expiry,
		loader: loader,
	}
}

func (m *memoryCache[T]) get(ctx context.Context) (*T, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.value == nil || time.Now().After(m.expiration) {
		val, err := m.loader(ctx)
		if err != nil {
			return nil, err
		}

		m.value = val
		m.expiration = time.Now().Add(m.expiry * time.Second)
	}

	return m.value, nil
}

/*
func (m *memoryCache[T]) invalidate() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.expiration = time.Time{}
	m.value = nil
}
*/

// @todo: allow cache expiry values to be configured

func initInfoByIdCaches(
	a *databricksSdk.AccountClient,
) {
	workspaceInfoCache = newMemoryCache(
		5 * time.Minute,
		func(ctx context.Context) (*map[int64]*workspaceInfo, error) {
			m, err := buildWorkspaceInfoByIdMap(ctx, a)
			if err != nil {
				return nil, err
			}

			return &m, nil
		},
	)

	clusterInfoCache = newMemoryCache(
		5 * time.Minute,
		func(ctx context.Context) (*map[string]*clusterInfo, error) {
			m, err := buildClusterInfoByIdMap(ctx, a)
			if err != nil {
				return nil, err
			}

			return &m, nil
		},
	)

	warehouseInfoCache = newMemoryCache(
		5 * time.Minute,
		func(ctx context.Context) (*map[string]*warehouseInfo, error) {
			m, err := buildWarehouseInfoByIdMap(ctx, a)
			if err != nil {
				return nil, err
			}

			return &m, nil
		},
	)
}

func buildWorkspaceInfoByIdMap(
	ctx context.Context,
	a *databricksSdk.AccountClient,
) (map[int64]*workspaceInfo, error) {
	log.Debugf("building workspace info by ID map...")

	workspaces, err := a.Workspaces.List(ctx)
	if err != nil {
		return nil, err
	}

	m := make(map[int64]*workspaceInfo)

	for _, workspace := range workspaces {
		workspaceInfo := &workspaceInfo{}
		workspaceInfo.name = workspace.WorkspaceName

		m[workspace.WorkspaceId] = workspaceInfo
	}

	return m, nil
}

func getWorkspaceInfoById(
	ctx context.Context,
	workspaceIdStr string,
) (*workspaceInfo, error) {
	workspaceId, err := strconv.ParseInt(workspaceIdStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("workspace ID is not an integer")
	}

	workspaceInfoMap, err := workspaceInfoCache.get(ctx)
	if err != nil {
		return nil, err
	}

	workspaceInfo, ok := (*workspaceInfoMap)[workspaceId]
	if ok {
		return workspaceInfo, nil
	}

	return nil, nil
}

func buildClusterInfoByIdMap(
	ctx context.Context,
	a *databricksSdk.AccountClient,
) (map[string]*clusterInfo, error) {
	log.Debugf("building cluster info by ID map...")

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
			// can be the same in different workspaces
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
	ctx context.Context,
	workspaceIdStr string,
	clusterIdStr string,
) (*clusterInfo, error) {
	workspaceId, err := strconv.ParseInt(workspaceIdStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("workspace ID is not an integer")
	}

	// namespace cluster ids with workspace id, just in case cluster ids can be
	// the same in different workspaces
	id := fmt.Sprintf(
		"%d.%s",
		workspaceId,
		clusterIdStr,
	)

	clusterInfoMap, err := clusterInfoCache.get(ctx)
	if err != nil {
		return nil, err
	}

	clusterInfo, ok := (*clusterInfoMap)[id]
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
			// ids can be the same in different workspaces
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
	ctx context.Context,
	workspaceIdStr string,
	warehouseIdStr string,
) (*warehouseInfo, error) {
	workspaceId, err := strconv.ParseInt(workspaceIdStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("workspace ID is not an integer")
	}

	// namespace warehouse ids with workspace id, just in case warehouse ids can
	// be the same in different workspaces
	id := fmt.Sprintf(
		"%d.%s",
		workspaceId,
		warehouseIdStr,
	)

	warehouseInfoMap, err := warehouseInfoCache.get(ctx)
	if err != nil {
		return nil, err
	}

	warehouseInfo, ok := (*warehouseInfoMap)[id]
	if ok {
		return warehouseInfo, nil
	}

	return nil, nil
}
