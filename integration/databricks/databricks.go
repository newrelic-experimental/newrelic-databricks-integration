package databricks

var recv_interval = 0

type AWSJobRuns struct {
	Runs []struct {
		JobId                int    `json:"job_id"`
		RunId                int    `json:"run_id"`
		NumberInJob          int    `json:"number_in_job"`
		CreatorUserName      string `json:"creator_user_name"`
		OriginalAttemptRunId int    `json:"original_attempt_run_id"`
		State                struct {
			LifeCycleState          string `json:"life_cycle_state"`
			QueueReason             string `json:"queue_reason"`
			ResultState             string `json:"result_state"`
			UserCancelledOrTimedout bool   `json:"user_cancelled_or_timedout"`
			StateMessage            string `json:"state_message"`
		} `json:"state"`
		Schedule struct {
			QuartzCronExpression string `json:"quartz_cron_expression"`
			TimezoneId           string `json:"timezone_id"`
			PauseStatus          string `json:"pause_status"`
		} `json:"schedule"`
		Tasks []struct {
			SetupDuration int    `json:"setup_duration"`
			StartTime     int64  `json:"start_time"`
			TaskKey       string `json:"task_key"`
			State         struct {
				LifeCycleState          string `json:"life_cycle_state"`
				ResultState             string `json:"result_state,omitempty"`
				StateMessage            string `json:"state_message"`
				UserCancelledOrTimedout bool   `json:"user_cancelled_or_timedout"`
			} `json:"state"`
			Description     string `json:"description"`
			JobClusterKey   string `json:"job_cluster_key,omitempty"`
			EndTime         int64  `json:"end_time"`
			RunPageUrl      string `json:"run_page_url"`
			RunId           int    `json:"run_id"`
			ClusterInstance struct {
				ClusterId      string `json:"cluster_id"`
				SparkContextId string `json:"spark_context_id,omitempty"`
			} `json:"cluster_instance"`
			SparkJarTask struct {
				MainClassName string `json:"main_class_name"`
			} `json:"spark_jar_task,omitempty"`
			Libraries []struct {
				Jar string `json:"jar"`
			} `json:"libraries,omitempty"`
			AttemptNumber     int    `json:"attempt_number"`
			CleanupDuration   int    `json:"cleanup_duration"`
			ExecutionDuration int    `json:"execution_duration"`
			RunIf             string `json:"run_if"`
			NotebookTask      struct {
				NotebookPath string `json:"notebook_path"`
				Source       string `json:"source"`
			} `json:"notebook_task,omitempty"`
			DependsOn []struct {
				TaskKey string `json:"task_key"`
			} `json:"depends_on,omitempty"`
			NewCluster struct {
				SparkVersion string      `json:"spark_version"`
				NodeTypeId   interface{} `json:"node_type_id"`
				SparkConf    struct {
					SparkSpeculation bool `json:"spark.speculation"`
				} `json:"spark_conf"`
				Autoscale struct {
					MinWorkers int `json:"min_workers"`
					MaxWorkers int `json:"max_workers"`
				} `json:"autoscale"`
			} `json:"new_cluster,omitempty"`
			ExistingClusterId string `json:"existing_cluster_id,omitempty"`
		} `json:"tasks"`
		JobClusters []struct {
			JobClusterKey string `json:"job_cluster_key"`
			NewCluster    struct {
				SparkVersion string      `json:"spark_version"`
				NodeTypeId   interface{} `json:"node_type_id"`
				SparkConf    struct {
					SparkSpeculation bool `json:"spark.speculation"`
				} `json:"spark_conf"`
				Autoscale struct {
					MinWorkers int `json:"min_workers"`
					MaxWorkers int `json:"max_workers"`
				} `json:"autoscale"`
			} `json:"new_cluster"`
		} `json:"job_clusters"`
		ClusterSpec struct {
			ExistingClusterId string `json:"existing_cluster_id"`
			NewCluster        struct {
				NumWorkers int `json:"num_workers"`
				Autoscale  struct {
					MinWorkers int `json:"min_workers"`
					MaxWorkers int `json:"max_workers"`
				} `json:"autoscale"`
				ClusterName  string `json:"cluster_name"`
				SparkVersion string `json:"spark_version"`
				SparkConf    struct {
					Property1 string `json:"property1"`
					Property2 string `json:"property2"`
				} `json:"spark_conf"`
				AwsAttributes struct {
					FirstOnDemand       string `json:"first_on_demand"`
					Availability        string `json:"availability"`
					ZoneId              string `json:"zone_id"`
					InstanceProfileArn  string `json:"instance_profile_arn"`
					SpotBidPricePercent string `json:"spot_bid_price_percent"`
					EbsVolumeType       string `json:"ebs_volume_type"`
					EbsVolumeCount      string `json:"ebs_volume_count"`
					EbsVolumeSize       int    `json:"ebs_volume_size"`
					EbsVolumeIops       int    `json:"ebs_volume_iops"`
					EbsVolumeThroughput int    `json:"ebs_volume_throughput"`
				} `json:"aws_attributes"`
				NodeTypeId       string   `json:"node_type_id"`
				DriverNodeTypeId string   `json:"driver_node_type_id"`
				SshPublicKeys    []string `json:"ssh_public_keys"`
				CustomTags       struct {
					Property1 string `json:"property1"`
					Property2 string `json:"property2"`
				} `json:"custom_tags"`
				ClusterLogConf struct {
					Dbfs struct {
						Destination string `json:"destination"`
					} `json:"dbfs"`
					S3 struct {
						Destination      string `json:"destination"`
						Region           string `json:"region"`
						Endpoint         string `json:"endpoint"`
						EnableEncryption bool   `json:"enable_encryption"`
						EncryptionType   string `json:"encryption_type"`
						KmsKey           string `json:"kms_key"`
						CannedAcl        string `json:"canned_acl"`
					} `json:"s3"`
				} `json:"cluster_log_conf"`
				InitScripts []struct {
					Workspace struct {
						Destination string `json:"destination"`
					} `json:"workspace"`
					Volumes struct {
						Destination string `json:"destination"`
					} `json:"volumes"`
					S3 struct {
						Destination      string `json:"destination"`
						Region           string `json:"region"`
						Endpoint         string `json:"endpoint"`
						EnableEncryption bool   `json:"enable_encryption"`
						EncryptionType   string `json:"encryption_type"`
						KmsKey           string `json:"kms_key"`
						CannedAcl        string `json:"canned_acl"`
					} `json:"s3"`
					File struct {
						Destination string `json:"destination"`
					} `json:"file"`
					Dbfs struct {
						Destination string `json:"destination"`
					} `json:"dbfs"`
					Abfss struct {
						Destination string `json:"destination"`
					} `json:"abfss"`
					Gcs struct {
						Destination string `json:"destination"`
					} `json:"gcs"`
				} `json:"init_scripts"`
				SparkEnvVars struct {
					Property1 string `json:"property1"`
					Property2 string `json:"property2"`
				} `json:"spark_env_vars"`
				AutoterminationMinutes    int    `json:"autotermination_minutes"`
				EnableElasticDisk         bool   `json:"enable_elastic_disk"`
				ClusterSource             string `json:"cluster_source"`
				InstancePoolId            string `json:"instance_pool_id"`
				PolicyId                  string `json:"policy_id"`
				EnableLocalDiskEncryption bool   `json:"enable_local_disk_encryption"`
				DriverInstancePoolId      string `json:"driver_instance_pool_id"`
				WorkloadType              struct {
					Clients struct {
						Notebooks string `json:"notebooks"`
						Jobs      string `json:"jobs"`
					} `json:"clients"`
				} `json:"workload_type"`
				RuntimeEngine string `json:"runtime_engine"`
				DockerImage   struct {
					Url       string `json:"url"`
					BasicAuth struct {
						Username string `json:"username"`
						Password string `json:"password"`
					} `json:"basic_auth"`
				} `json:"docker_image"`
				DataSecurityMode         string `json:"data_security_mode"`
				SingleUserName           string `json:"single_user_name"`
				ApplyPolicyDefaultValues string `json:"apply_policy_default_values"`
			} `json:"new_cluster"`
			Libraries []struct {
				Jar  string `json:"jar"`
				Egg  string `json:"egg"`
				Pypi struct {
					Package string `json:"package"`
					Repo    string `json:"repo"`
				} `json:"pypi"`
				Maven struct {
					Coordinates string   `json:"coordinates"`
					Repo        string   `json:"repo"`
					Exclusions  []string `json:"exclusions"`
				} `json:"maven"`
				Cran struct {
					Package string `json:"package"`
					Repo    string `json:"repo"`
				} `json:"cran"`
				Whl string `json:"whl"`
			} `json:"libraries"`
		} `json:"cluster_spec"`
		ClusterInstance struct {
			ClusterId      string `json:"cluster_id"`
			SparkContextId string `json:"spark_context_id"`
		} `json:"cluster_instance"`
		GitSource struct {
			GitUrl      string `json:"git_url"`
			GitBranch   string `json:"git_branch"`
			GitProvider string `json:"git_provider"`
		} `json:"git_source"`
		OverridingParameters struct {
			JarParams      []string `json:"jar_params"`
			NotebookParams struct {
				Name string `json:"name"`
				Age  string `json:"age"`
			} `json:"notebook_params"`
			PythonParams      []string `json:"python_params"`
			SparkSubmitParams []string `json:"spark_submit_params"`
			PythonNamedParams struct {
				Name string `json:"name"`
				Data string `json:"data"`
			} `json:"python_named_params"`
			PipelineParams struct {
				FullRefresh bool `json:"full_refresh"`
			} `json:"pipeline_params"`
			SqlParams struct {
				Name string `json:"name"`
				Age  string `json:"age"`
			} `json:"sql_params"`
			DbtCommands   []string `json:"dbt_commands"`
			JobParameters struct {
				Name string `json:"name"`
				Age  string `json:"age"`
			} `json:"job_parameters"`
		} `json:"overriding_parameters"`
		StartTime         int64 `json:"start_time"`
		SetupDuration     int   `json:"setup_duration"`
		ExecutionDuration int   `json:"execution_duration"`
		CleanupDuration   int   `json:"cleanup_duration"`
		EndTime           int64 `json:"end_time"`
		TriggerInfo       struct {
			RunId int `json:"run_id"`
		} `json:"trigger_info"`
		RunDuration   int    `json:"run_duration"`
		Trigger       string `json:"trigger"`
		RunName       string `json:"run_name"`
		RunPageUrl    string `json:"run_page_url"`
		RunType       string `json:"run_type"`
		AttemptNumber int    `json:"attempt_number"`
		JobParameters []struct {
			Name    string `json:"name"`
			Default string `json:"default"`
			Value   string `json:"value"`
		} `json:"job_parameters"`
	} `json:"runs"`
	HasMore       bool   `json:"has_more"`
	NextPageToken string `json:"next_page_token"`
	PrevPageToken string `json:"prev_page_token"`
}

type AWSQueriesList struct {
	NextPageToken string `json:"next_page_token"`
	HasNextPage   bool   `json:"has_next_page"`
	Res           []struct {
		QueryId                 string `json:"query_id"`
		Status                  string `json:"status"`
		QueryText               string `json:"query_text"`
		QueryStartTimeMs        int64  `json:"query_start_time_ms"`
		ExecutionEndTimeMs      int64  `json:"execution_end_time_ms"`
		QueryEndTimeMs          int64  `json:"query_end_time_ms"`
		UserId                  int64  `json:"user_id"`
		UserName                string `json:"user_name"`
		SparkUiUrl              string `json:"spark_ui_url"`
		EndpointId              string `json:"endpoint_id"`
		WarehouseId             string `json:"warehouse_id"`
		LookupKey               string `json:"lookup_key"`
		ErrorMessage            string `json:"error_message"`
		RowsProduced            int    `json:"rows_produced"`
		CanSubscribeToLiveQuery bool   `json:"canSubscribeToLiveQuery"`
		Metrics                 struct {
			TotalTimeMs                     int   `json:"total_time_ms"`
			ReadBytes                       int   `json:"read_bytes"`
			RowsProducedCount               int   `json:"rows_produced_count"`
			CompilationTimeMs               int   `json:"compilation_time_ms"`
			ExecutionTimeMs                 int   `json:"execution_time_ms"`
			ReadRemoteBytes                 int   `json:"read_remote_bytes"`
			WriteRemoteBytes                int   `json:"write_remote_bytes"`
			ReadCacheBytes                  int   `json:"read_cache_bytes"`
			SpillToDiskBytes                int   `json:"spill_to_disk_bytes"`
			TaskTotalTimeMs                 int   `json:"task_total_time_ms"`
			ReadFilesCount                  int   `json:"read_files_count"`
			ReadPartitionsCount             int   `json:"read_partitions_count"`
			PhotonTotalTimeMs               int   `json:"photon_total_time_ms"`
			RowsReadCount                   int   `json:"rows_read_count"`
			ResultFetchTimeMs               int   `json:"result_fetch_time_ms"`
			NetworkSentBytes                int   `json:"network_sent_bytes"`
			ResultFromCache                 bool  `json:"result_from_cache"`
			PrunedBytes                     int   `json:"pruned_bytes"`
			PrunedFilesCount                int   `json:"pruned_files_count"`
			ProvisioningQueueStartTimestamp int64 `json:"provisioning_queue_start_timestamp"`
			OverloadingQueueStartTimestamp  int64 `json:"overloading_queue_start_timestamp"`
			QueryCompilationStartTimestamp  int64 `json:"query_compilation_start_timestamp"`
			MetadataTimeMs                  int   `json:"metadata_time_ms"`
			PlanningTimeMs                  int   `json:"planning_time_ms"`
			QueryExecutionTimeMs            int   `json:"query_execution_time_ms"`
		} `json:"metrics"`
		IsFinal     bool `json:"is_final"`
		ChannelUsed struct {
			Name         string `json:"name"`
			DbsqlVersion string `json:"dbsql_version"`
		} `json:"channel_used"`
		Duration           int    `json:"duration"`
		ExecutedAsUserId   int64  `json:"executed_as_user_id"`
		ExecutedAsUserName string `json:"executed_as_user_name"`
		PlansState         string `json:"plans_state"`
		StatementType      string `json:"statement_type"`
	} `json:"res"`
}

type GCPQueriesList struct {
	NextPageToken string `json:"next_page_token"`
	HasNextPage   bool   `json:"has_next_page"`
	Res           []struct {
		QueryId                 string `json:"query_id"`
		Status                  string `json:"status"`
		QueryText               string `json:"query_text"`
		QueryStartTimeMs        int64  `json:"query_start_time_ms"`
		ExecutionEndTimeMs      int64  `json:"execution_end_time_ms"`
		QueryEndTimeMs          int64  `json:"query_end_time_ms"`
		UserId                  int64  `json:"user_id"`
		UserName                string `json:"user_name"`
		SparkUiUrl              string `json:"spark_ui_url"`
		EndpointId              string `json:"endpoint_id"`
		WarehouseId             string `json:"warehouse_id"`
		LookupKey               string `json:"lookup_key"`
		ErrorMessage            string `json:"error_message"`
		RowsProduced            int    `json:"rows_produced"`
		CanSubscribeToLiveQuery bool   `json:"canSubscribeToLiveQuery"`
		Metrics                 struct {
			TotalTimeMs                     int   `json:"total_time_ms"`
			ReadBytes                       int   `json:"read_bytes"`
			RowsProducedCount               int   `json:"rows_produced_count"`
			CompilationTimeMs               int   `json:"compilation_time_ms"`
			ExecutionTimeMs                 int   `json:"execution_time_ms"`
			ReadRemoteBytes                 int   `json:"read_remote_bytes"`
			WriteRemoteBytes                int   `json:"write_remote_bytes"`
			ReadCacheBytes                  int   `json:"read_cache_bytes"`
			SpillToDiskBytes                int   `json:"spill_to_disk_bytes"`
			TaskTotalTimeMs                 int   `json:"task_total_time_ms"`
			ReadFilesCount                  int   `json:"read_files_count"`
			ReadPartitionsCount             int   `json:"read_partitions_count"`
			PhotonTotalTimeMs               int   `json:"photon_total_time_ms"`
			RowsReadCount                   int   `json:"rows_read_count"`
			ResultFetchTimeMs               int   `json:"result_fetch_time_ms"`
			NetworkSentBytes                int   `json:"network_sent_bytes"`
			ResultFromCache                 bool  `json:"result_from_cache"`
			PrunedBytes                     int   `json:"pruned_bytes"`
			PrunedFilesCount                int   `json:"pruned_files_count"`
			ProvisioningQueueStartTimestamp int64 `json:"provisioning_queue_start_timestamp"`
			OverloadingQueueStartTimestamp  int64 `json:"overloading_queue_start_timestamp"`
			QueryCompilationStartTimestamp  int64 `json:"query_compilation_start_timestamp"`
			MetadataTimeMs                  int   `json:"metadata_time_ms"`
			PlanningTimeMs                  int   `json:"planning_time_ms"`
			QueryExecutionTimeMs            int   `json:"query_execution_time_ms"`
		} `json:"metrics"`
		IsFinal     bool `json:"is_final"`
		ChannelUsed struct {
			Name         string `json:"name"`
			DbsqlVersion string `json:"dbsql_version"`
		} `json:"channel_used"`
		Duration           int    `json:"duration"`
		ExecutedAsUserId   int64  `json:"executed_as_user_id"`
		ExecutedAsUserName string `json:"executed_as_user_name"`
		PlansState         string `json:"plans_state"`
		StatementType      string `json:"statement_type"`
	} `json:"res"`
}

type AzureQueriesList struct {
	NextPageToken string `json:"next_page_token"`
	HasNextPage   bool   `json:"has_next_page"`
	Res           []struct {
		QueryId                 string `json:"query_id"`
		Status                  string `json:"status"`
		QueryText               string `json:"query_text"`
		QueryStartTimeMs        int64  `json:"query_start_time_ms"`
		ExecutionEndTimeMs      int64  `json:"execution_end_time_ms"`
		QueryEndTimeMs          int64  `json:"query_end_time_ms"`
		UserId                  int64  `json:"user_id"`
		UserName                string `json:"user_name"`
		SparkUiUrl              string `json:"spark_ui_url"`
		EndpointId              string `json:"endpoint_id"`
		WarehouseId             string `json:"warehouse_id"`
		LookupKey               string `json:"lookup_key"`
		ErrorMessage            string `json:"error_message"`
		RowsProduced            int    `json:"rows_produced"`
		CanSubscribeToLiveQuery bool   `json:"canSubscribeToLiveQuery"`
		Metrics                 struct {
			TotalTimeMs                     int   `json:"total_time_ms"`
			ReadBytes                       int   `json:"read_bytes"`
			RowsProducedCount               int   `json:"rows_produced_count"`
			CompilationTimeMs               int   `json:"compilation_time_ms"`
			ExecutionTimeMs                 int   `json:"execution_time_ms"`
			ReadRemoteBytes                 int   `json:"read_remote_bytes"`
			WriteRemoteBytes                int   `json:"write_remote_bytes"`
			ReadCacheBytes                  int   `json:"read_cache_bytes"`
			SpillToDiskBytes                int   `json:"spill_to_disk_bytes"`
			TaskTotalTimeMs                 int   `json:"task_total_time_ms"`
			ReadFilesCount                  int   `json:"read_files_count"`
			ReadPartitionsCount             int   `json:"read_partitions_count"`
			PhotonTotalTimeMs               int   `json:"photon_total_time_ms"`
			RowsReadCount                   int   `json:"rows_read_count"`
			ResultFetchTimeMs               int   `json:"result_fetch_time_ms"`
			NetworkSentBytes                int   `json:"network_sent_bytes"`
			ResultFromCache                 bool  `json:"result_from_cache"`
			PrunedBytes                     int   `json:"pruned_bytes"`
			PrunedFilesCount                int   `json:"pruned_files_count"`
			ProvisioningQueueStartTimestamp int64 `json:"provisioning_queue_start_timestamp"`
			OverloadingQueueStartTimestamp  int64 `json:"overloading_queue_start_timestamp"`
			QueryCompilationStartTimestamp  int64 `json:"query_compilation_start_timestamp"`
			MetadataTimeMs                  int   `json:"metadata_time_ms"`
			PlanningTimeMs                  int   `json:"planning_time_ms"`
			QueryExecutionTimeMs            int   `json:"query_execution_time_ms"`
		} `json:"metrics"`
		IsFinal     bool `json:"is_final"`
		ChannelUsed struct {
			Name         string `json:"name"`
			DbsqlVersion string `json:"dbsql_version"`
		} `json:"channel_used"`
		Duration           int    `json:"duration"`
		ExecutedAsUserId   int64  `json:"executed_as_user_id"`
		ExecutedAsUserName string `json:"executed_as_user_name"`
		PlansState         string `json:"plans_state"`
		StatementType      string `json:"statement_type"`
	} `json:"res"`
}

type GCPJobRuns struct {
	Runs []struct {
		JobId                int    `json:"job_id"`
		RunId                int    `json:"run_id"`
		CreatorUserName      string `json:"creator_user_name"`
		NumberInJob          int    `json:"number_in_job"`
		OriginalAttemptRunId int    `json:"original_attempt_run_id"`
		State                struct {
			LifeCycleState          string `json:"life_cycle_state"`
			ResultState             string `json:"result_state"`
			StateMessage            string `json:"state_message"`
			UserCancelledOrTimedout bool   `json:"user_cancelled_or_timedout"`
			QueueReason             string `json:"queue_reason"`
		} `json:"state"`
		Schedule struct {
			QuartzCronExpression string `json:"quartz_cron_expression"`
			TimezoneId           string `json:"timezone_id"`
			PauseStatus          string `json:"pause_status"`
		} `json:"schedule"`
		ClusterSpec struct {
			ExistingClusterId string `json:"existing_cluster_id"`
			NewCluster        struct {
				NumWorkers int `json:"num_workers"`
				Autoscale  struct {
					MinWorkers int `json:"min_workers"`
					MaxWorkers int `json:"max_workers"`
				} `json:"autoscale"`
				ClusterName  string `json:"cluster_name"`
				SparkVersion string `json:"spark_version"`
				SparkConf    struct {
					Property1 string `json:"property1"`
					Property2 string `json:"property2"`
				} `json:"spark_conf"`
				GcpAttributes struct {
					GoogleServiceAccount    string `json:"google_service_account"`
					BootDiskSize            int    `json:"boot_disk_size"`
					Availability            string `json:"availability"`
					LocalSsdCount           int    `json:"local_ssd_count"`
					UsePreemptibleExecutors bool   `json:"use_preemptible_executors"`
					ZoneId                  string `json:"zone_id"`
				} `json:"gcp_attributes"`
				NodeTypeId       string   `json:"node_type_id"`
				DriverNodeTypeId string   `json:"driver_node_type_id"`
				SshPublicKeys    []string `json:"ssh_public_keys"`
				CustomTags       struct {
					Property1 string `json:"property1"`
					Property2 string `json:"property2"`
				} `json:"custom_tags"`
				ClusterLogConf struct {
					Dbfs struct {
						Destination string `json:"destination"`
					} `json:"dbfs"`
				} `json:"cluster_log_conf"`
				InitScripts []struct {
					Workspace struct {
						Destination string `json:"destination"`
					} `json:"workspace"`
					Volumes struct {
						Destination string `json:"destination"`
					} `json:"volumes"`
					File struct {
						Destination string `json:"destination"`
					} `json:"file"`
					Dbfs struct {
						Destination string `json:"destination"`
					} `json:"dbfs"`
					Abfss struct {
						Destination string `json:"destination"`
					} `json:"abfss"`
					Gcs struct {
						Destination string `json:"destination"`
					} `json:"gcs"`
				} `json:"init_scripts"`
				SparkEnvVars struct {
					Property1 string `json:"property1"`
					Property2 string `json:"property2"`
				} `json:"spark_env_vars"`
				AutoterminationMinutes    int    `json:"autotermination_minutes"`
				EnableElasticDisk         bool   `json:"enable_elastic_disk"`
				ClusterSource             string `json:"cluster_source"`
				InstancePoolId            string `json:"instance_pool_id"`
				PolicyId                  string `json:"policy_id"`
				EnableLocalDiskEncryption bool   `json:"enable_local_disk_encryption"`
				DriverInstancePoolId      string `json:"driver_instance_pool_id"`
				WorkloadType              struct {
					Clients struct {
						Notebooks string `json:"notebooks"`
						Jobs      string `json:"jobs"`
					} `json:"clients"`
				} `json:"workload_type"`
				RuntimeEngine string `json:"runtime_engine"`
				DockerImage   struct {
					Url       string `json:"url"`
					BasicAuth struct {
						Username string `json:"username"`
						Password string `json:"password"`
					} `json:"basic_auth"`
				} `json:"docker_image"`
				DataSecurityMode         string `json:"data_security_mode"`
				SingleUserName           string `json:"single_user_name"`
				ApplyPolicyDefaultValues string `json:"apply_policy_default_values"`
			} `json:"new_cluster"`
			JobClusterKey string `json:"job_cluster_key"`
			Libraries     []struct {
				Jar  string `json:"jar"`
				Egg  string `json:"egg"`
				Pypi struct {
					Package string `json:"package"`
					Repo    string `json:"repo"`
				} `json:"pypi"`
				Maven struct {
					Coordinates string   `json:"coordinates"`
					Repo        string   `json:"repo"`
					Exclusions  []string `json:"exclusions"`
				} `json:"maven"`
				Cran struct {
					Package string `json:"package"`
					Repo    string `json:"repo"`
				} `json:"cran"`
				Whl string `json:"whl"`
			} `json:"libraries"`
		} `json:"cluster_spec"`
		ClusterInstance struct {
			ClusterId      string `json:"cluster_id"`
			SparkContextId string `json:"spark_context_id"`
		} `json:"cluster_instance"`
		OverridingParameters struct {
			JarParams      []string `json:"jar_params"`
			NotebookParams struct {
				Age  string `json:"age"`
				Name string `json:"name"`
			} `json:"notebook_params"`
			PythonParams      []string `json:"python_params"`
			SparkSubmitParams []string `json:"spark_submit_params"`
			PythonNamedParams struct {
				Data string `json:"data"`
				Name string `json:"name"`
			} `json:"python_named_params"`
			DbtCommands    []string `json:"dbt_commands"`
			PipelineParams struct {
				FullRefresh bool `json:"full_refresh"`
			} `json:"pipeline_params"`
		} `json:"overriding_parameters"`
		StartTime         int64  `json:"start_time"`
		SetupDuration     int    `json:"setup_duration"`
		ExecutionDuration int    `json:"execution_duration"`
		CleanupDuration   int    `json:"cleanup_duration"`
		EndTime           int64  `json:"end_time"`
		RunDuration       int    `json:"run_duration"`
		QueueDuration     int64  `json:"queue_duration"`
		Trigger           string `json:"trigger"`
		TriggerInfo       struct {
			RunId int `json:"run_id"`
		} `json:"trigger_info"`
		RunName    string `json:"run_name"`
		RunPageUrl string `json:"run_page_url"`
		RunType    string `json:"run_type"`
		Tasks      []struct {
			SetupDuration int    `json:"setup_duration"`
			StartTime     int64  `json:"start_time"`
			TaskKey       string `json:"task_key"`
			State         struct {
				LifeCycleState          string `json:"life_cycle_state"`
				ResultState             string `json:"result_state,omitempty"`
				StateMessage            string `json:"state_message"`
				UserCancelledOrTimedout bool   `json:"user_cancelled_or_timedout"`
			} `json:"state"`
			Description     string `json:"description"`
			JobClusterKey   string `json:"job_cluster_key,omitempty"`
			EndTime         int64  `json:"end_time"`
			RunPageUrl      string `json:"run_page_url"`
			RunId           int    `json:"run_id"`
			ClusterInstance struct {
				ClusterId      string `json:"cluster_id"`
				SparkContextId string `json:"spark_context_id,omitempty"`
			} `json:"cluster_instance"`
			SparkJarTask struct {
				MainClassName string `json:"main_class_name"`
			} `json:"spark_jar_task,omitempty"`
			Libraries []struct {
				Jar string `json:"jar"`
			} `json:"libraries,omitempty"`
			AttemptNumber     int    `json:"attempt_number"`
			CleanupDuration   int    `json:"cleanup_duration"`
			ExecutionDuration int    `json:"execution_duration"`
			RunIf             string `json:"run_if"`
			NotebookTask      struct {
				NotebookPath string `json:"notebook_path"`
				Source       string `json:"source"`
			} `json:"notebook_task,omitempty"`
			DependsOn []struct {
				TaskKey string `json:"task_key"`
			} `json:"depends_on,omitempty"`
			NewCluster struct {
				Autoscale struct {
					MaxWorkers int `json:"max_workers"`
					MinWorkers int `json:"min_workers"`
				} `json:"autoscale"`
				NodeTypeId interface{} `json:"node_type_id"`
				SparkConf  struct {
					SparkSpeculation bool `json:"spark.speculation"`
				} `json:"spark_conf"`
				SparkVersion string `json:"spark_version"`
			} `json:"new_cluster,omitempty"`
			ExistingClusterId string `json:"existing_cluster_id,omitempty"`
		} `json:"tasks"`
		Description   string `json:"description"`
		AttemptNumber int    `json:"attempt_number"`
		JobClusters   []struct {
			JobClusterKey string `json:"job_cluster_key"`
			NewCluster    struct {
				Autoscale struct {
					MaxWorkers int `json:"max_workers"`
					MinWorkers int `json:"min_workers"`
				} `json:"autoscale"`
				NodeTypeId interface{} `json:"node_type_id"`
				SparkConf  struct {
					SparkSpeculation bool `json:"spark.speculation"`
				} `json:"spark_conf"`
				SparkVersion string `json:"spark_version"`
			} `json:"new_cluster"`
		} `json:"job_clusters"`
		GitSource struct {
			GitBranch   string `json:"git_branch"`
			GitProvider string `json:"git_provider"`
			GitUrl      string `json:"git_url"`
		} `json:"git_source"`
		RepairHistory []struct {
			Type      string `json:"type"`
			StartTime int64  `json:"start_time"`
			EndTime   int64  `json:"end_time"`
			State     struct {
				LifeCycleState          string `json:"life_cycle_state"`
				ResultState             string `json:"result_state"`
				StateMessage            string `json:"state_message"`
				UserCancelledOrTimedout bool   `json:"user_cancelled_or_timedout"`
				QueueReason             string `json:"queue_reason"`
			} `json:"state"`
			Id         int64   `json:"id"`
			TaskRunIds []int64 `json:"task_run_ids"`
		} `json:"repair_history"`
		JobParameters []struct {
			Default string `json:"default"`
			Name    string `json:"name"`
			Value   string `json:"value"`
		} `json:"job_parameters"`
	} `json:"runs"`
	HasMore       bool   `json:"has_more"`
	NextPageToken string `json:"next_page_token"`
	PrevPageToken string `json:"prev_page_token"`
}

type AzureJobRuns struct {
	Runs []struct {
		JobId                int    `json:"job_id"`
		RunId                int    `json:"run_id"`
		CreatorUserName      string `json:"creator_user_name"`
		NumberInJob          int    `json:"number_in_job"`
		OriginalAttemptRunId int    `json:"original_attempt_run_id"`
		State                struct {
			LifeCycleState          string `json:"life_cycle_state"`
			ResultState             string `json:"result_state"`
			StateMessage            string `json:"state_message"`
			UserCancelledOrTimedout bool   `json:"user_cancelled_or_timedout"`
			QueueReason             string `json:"queue_reason"`
		} `json:"state"`
		Schedule struct {
			QuartzCronExpression string `json:"quartz_cron_expression"`
			TimezoneId           string `json:"timezone_id"`
			PauseStatus          string `json:"pause_status"`
		} `json:"schedule"`
		ClusterSpec struct {
			ExistingClusterId string `json:"existing_cluster_id"`
			NewCluster        struct {
				NumWorkers int `json:"num_workers"`
				Autoscale  struct {
					MinWorkers int `json:"min_workers"`
					MaxWorkers int `json:"max_workers"`
				} `json:"autoscale"`
				ClusterName  string `json:"cluster_name"`
				SparkVersion string `json:"spark_version"`
				SparkConf    struct {
					Property1 string `json:"property1"`
					Property2 string `json:"property2"`
				} `json:"spark_conf"`
				AzureAttributes struct {
					LogAnalyticsInfo struct {
						LogAnalyticsWorkspaceId string `json:"log_analytics_workspace_id"`
						LogAnalyticsPrimaryKey  string `json:"log_analytics_primary_key"`
					} `json:"log_analytics_info"`
					FirstOnDemand   string `json:"first_on_demand"`
					Availability    string `json:"availability"`
					SpotBidMaxPrice string `json:"spot_bid_max_price"`
				} `json:"azure_attributes"`
				NodeTypeId       string   `json:"node_type_id"`
				DriverNodeTypeId string   `json:"driver_node_type_id"`
				SshPublicKeys    []string `json:"ssh_public_keys"`
				CustomTags       struct {
					Property1 string `json:"property1"`
					Property2 string `json:"property2"`
				} `json:"custom_tags"`
				ClusterLogConf struct {
					Dbfs struct {
						Destination string `json:"destination"`
					} `json:"dbfs"`
				} `json:"cluster_log_conf"`
				InitScripts []struct {
					Workspace struct {
						Destination string `json:"destination"`
					} `json:"workspace"`
					Volumes struct {
						Destination string `json:"destination"`
					} `json:"volumes"`
					File struct {
						Destination string `json:"destination"`
					} `json:"file"`
					Dbfs struct {
						Destination string `json:"destination"`
					} `json:"dbfs"`
					Abfss struct {
						Destination string `json:"destination"`
					} `json:"abfss"`
					Gcs struct {
						Destination string `json:"destination"`
					} `json:"gcs"`
				} `json:"init_scripts"`
				SparkEnvVars struct {
					Property1 string `json:"property1"`
					Property2 string `json:"property2"`
				} `json:"spark_env_vars"`
				AutoterminationMinutes    int    `json:"autotermination_minutes"`
				EnableElasticDisk         bool   `json:"enable_elastic_disk"`
				ClusterSource             string `json:"cluster_source"`
				InstancePoolId            string `json:"instance_pool_id"`
				PolicyId                  string `json:"policy_id"`
				EnableLocalDiskEncryption bool   `json:"enable_local_disk_encryption"`
				DriverInstancePoolId      string `json:"driver_instance_pool_id"`
				WorkloadType              struct {
					Clients struct {
						Notebooks string `json:"notebooks"`
						Jobs      string `json:"jobs"`
					} `json:"clients"`
				} `json:"workload_type"`
				RuntimeEngine string `json:"runtime_engine"`
				DockerImage   struct {
					Url       string `json:"url"`
					BasicAuth struct {
						Username string `json:"username"`
						Password string `json:"password"`
					} `json:"basic_auth"`
				} `json:"docker_image"`
				DataSecurityMode         string `json:"data_security_mode"`
				SingleUserName           string `json:"single_user_name"`
				ApplyPolicyDefaultValues string `json:"apply_policy_default_values"`
			} `json:"new_cluster"`
			JobClusterKey string `json:"job_cluster_key"`
			Libraries     []struct {
				Jar  string `json:"jar"`
				Egg  string `json:"egg"`
				Pypi struct {
					Package string `json:"package"`
					Repo    string `json:"repo"`
				} `json:"pypi"`
				Maven struct {
					Coordinates string   `json:"coordinates"`
					Repo        string   `json:"repo"`
					Exclusions  []string `json:"exclusions"`
				} `json:"maven"`
				Cran struct {
					Package string `json:"package"`
					Repo    string `json:"repo"`
				} `json:"cran"`
				Whl string `json:"whl"`
			} `json:"libraries"`
		} `json:"cluster_spec"`
		ClusterInstance struct {
			ClusterId      string `json:"cluster_id"`
			SparkContextId string `json:"spark_context_id"`
		} `json:"cluster_instance"`
		OverridingParameters struct {
			JarParams      []string `json:"jar_params"`
			NotebookParams struct {
				Age  string `json:"age"`
				Name string `json:"name"`
			} `json:"notebook_params"`
			PythonParams      []string `json:"python_params"`
			SparkSubmitParams []string `json:"spark_submit_params"`
			PythonNamedParams struct {
				Data string `json:"data"`
				Name string `json:"name"`
			} `json:"python_named_params"`
			SqlParams struct {
				Age  string `json:"age"`
				Name string `json:"name"`
			} `json:"sql_params"`
			DbtCommands    []string `json:"dbt_commands"`
			PipelineParams struct {
				FullRefresh bool `json:"full_refresh"`
			} `json:"pipeline_params"`
		} `json:"overriding_parameters"`
		StartTime         int64  `json:"start_time"`
		SetupDuration     int    `json:"setup_duration"`
		ExecutionDuration int    `json:"execution_duration"`
		CleanupDuration   int    `json:"cleanup_duration"`
		EndTime           int64  `json:"end_time"`
		RunDuration       int    `json:"run_duration"`
		QueueDuration     int64  `json:"queue_duration"`
		Trigger           string `json:"trigger"`
		TriggerInfo       struct {
			RunId int `json:"run_id"`
		} `json:"trigger_info"`
		RunName    string `json:"run_name"`
		RunPageUrl string `json:"run_page_url"`
		RunType    string `json:"run_type"`
		Tasks      []struct {
			SetupDuration int    `json:"setup_duration"`
			StartTime     int64  `json:"start_time"`
			TaskKey       string `json:"task_key"`
			State         struct {
				LifeCycleState          string `json:"life_cycle_state"`
				ResultState             string `json:"result_state,omitempty"`
				StateMessage            string `json:"state_message"`
				UserCancelledOrTimedout bool   `json:"user_cancelled_or_timedout"`
			} `json:"state"`
			Description     string `json:"description"`
			JobClusterKey   string `json:"job_cluster_key,omitempty"`
			EndTime         int64  `json:"end_time"`
			RunPageUrl      string `json:"run_page_url"`
			RunId           int    `json:"run_id"`
			ClusterInstance struct {
				ClusterId      string `json:"cluster_id"`
				SparkContextId string `json:"spark_context_id,omitempty"`
			} `json:"cluster_instance"`
			SparkJarTask struct {
				MainClassName string `json:"main_class_name"`
			} `json:"spark_jar_task,omitempty"`
			Libraries []struct {
				Jar string `json:"jar"`
			} `json:"libraries,omitempty"`
			AttemptNumber     int    `json:"attempt_number"`
			CleanupDuration   int    `json:"cleanup_duration"`
			ExecutionDuration int    `json:"execution_duration"`
			RunIf             string `json:"run_if"`
			NotebookTask      struct {
				NotebookPath string `json:"notebook_path"`
				Source       string `json:"source"`
			} `json:"notebook_task,omitempty"`
			DependsOn []struct {
				TaskKey string `json:"task_key"`
			} `json:"depends_on,omitempty"`
			NewCluster struct {
				Autoscale struct {
					MaxWorkers int `json:"max_workers"`
					MinWorkers int `json:"min_workers"`
				} `json:"autoscale"`
				NodeTypeId interface{} `json:"node_type_id"`
				SparkConf  struct {
					SparkSpeculation bool `json:"spark.speculation"`
				} `json:"spark_conf"`
				SparkVersion string `json:"spark_version"`
			} `json:"new_cluster,omitempty"`
			ExistingClusterId string `json:"existing_cluster_id,omitempty"`
		} `json:"tasks"`
		Description   string `json:"description"`
		AttemptNumber int    `json:"attempt_number"`
		JobClusters   []struct {
			JobClusterKey string `json:"job_cluster_key"`
			NewCluster    struct {
				Autoscale struct {
					MaxWorkers int `json:"max_workers"`
					MinWorkers int `json:"min_workers"`
				} `json:"autoscale"`
				NodeTypeId interface{} `json:"node_type_id"`
				SparkConf  struct {
					SparkSpeculation bool `json:"spark.speculation"`
				} `json:"spark_conf"`
				SparkVersion string `json:"spark_version"`
			} `json:"new_cluster"`
		} `json:"job_clusters"`
		GitSource struct {
			GitBranch   string `json:"git_branch"`
			GitProvider string `json:"git_provider"`
			GitUrl      string `json:"git_url"`
		} `json:"git_source"`
		RepairHistory []struct {
			Type      string `json:"type"`
			StartTime int64  `json:"start_time"`
			EndTime   int64  `json:"end_time"`
			State     struct {
				LifeCycleState          string `json:"life_cycle_state"`
				ResultState             string `json:"result_state"`
				StateMessage            string `json:"state_message"`
				UserCancelledOrTimedout bool   `json:"user_cancelled_or_timedout"`
				QueueReason             string `json:"queue_reason"`
			} `json:"state"`
			Id         int64   `json:"id"`
			TaskRunIds []int64 `json:"task_run_ids"`
		} `json:"repair_history"`
		JobParameters []struct {
			Default string `json:"default"`
			Name    string `json:"name"`
			Value   string `json:"value"`
		} `json:"job_parameters"`
	} `json:"runs"`
	HasMore       bool   `json:"has_more"`
	NextPageToken string `json:"next_page_token"`
	PrevPageToken string `json:"prev_page_token"`
}
