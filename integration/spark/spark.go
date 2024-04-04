package spark

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

var recv_interval = 0
var licenseKey = ""
var sparkEndpoint = ""

type ApplicationsResponse []struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Attempts []struct {
		StartTime        string `json:"startTime"`
		EndTime          string `json:"endTime"`
		LastUpdated      string `json:"lastUpdated"`
		Duration         int64  `json:"duration"`
		SparkUser        string `json:"sparkUser"`
		Completed        bool   `json:"completed"`
		AppSparkVersion  string `json:"appSparkVersion"`
		EndTimeEpoch     int64  `json:"endTimeEpoch"`
		StartTimeEpoch   int64  `json:"startTimeEpoch"`
		LastUpdatedEpoch int64  `json:"lastUpdatedEpoch"`
	} `json:"attempts"`
}

type SparkJob struct {
	JobId               int           `json:"jobId"`
	Name                string        `json:"name"`
	Description         string        `json:"description"`
	SubmissionTime      string        `json:"submissionTime"`
	CompletionTime      string        `json:"completionTime"`
	StageIds            []int         `json:"stageIds"`
	JobGroup            string        `json:"jobGroup"`
	JobTags             []interface{} `json:"jobTags"`
	Status              string        `json:"status"`
	NumTasks            int           `json:"numTasks"`
	NumActiveTasks      int           `json:"numActiveTasks"`
	NumCompletedTasks   int           `json:"numCompletedTasks"`
	NumSkippedTasks     int           `json:"numSkippedTasks"`
	NumFailedTasks      int           `json:"numFailedTasks"`
	NumKilledTasks      int           `json:"numKilledTasks"`
	NumCompletedIndices int           `json:"numCompletedIndices"`
	NumActiveStages     int           `json:"numActiveStages"`
	NumCompletedStages  int           `json:"numCompletedStages"`
	NumSkippedStages    int           `json:"numSkippedStages"`
	NumFailedStages     int           `json:"numFailedStages"`
	KilledTasksSummary  struct {
	} `json:"killedTasksSummary"`
}

type SparkExecutor struct {
	Id                string `json:"id"`
	HostPort          string `json:"hostPort"`
	IsActive          bool   `json:"isActive"`
	RddBlocks         int    `json:"rddBlocks"`
	MemoryUsed        int    `json:"memoryUsed"`
	DiskUsed          int    `json:"diskUsed"`
	TotalCores        int    `json:"totalCores"`
	MaxTasks          int    `json:"maxTasks"`
	ActiveTasks       int    `json:"activeTasks"`
	FailedTasks       int    `json:"failedTasks"`
	CompletedTasks    int    `json:"completedTasks"`
	TotalTasks        int    `json:"totalTasks"`
	TotalDuration     int    `json:"totalDuration"`
	TotalGCTime       int    `json:"totalGCTime"`
	TotalInputBytes   int    `json:"totalInputBytes"`
	TotalShuffleRead  int    `json:"totalShuffleRead"`
	TotalShuffleWrite int    `json:"totalShuffleWrite"`
	IsBlacklisted     bool   `json:"isBlacklisted"`
	MaxMemory         int    `json:"maxMemory"`
	AddTime           string `json:"addTime"`
	ExecutorLogs      struct {
		Stdout string `json:"stdout,omitempty"`
		Stderr string `json:"stderr,omitempty"`
	} `json:"executorLogs"`
	MemoryMetrics struct {
		UsedOnHeapStorageMemory   int `json:"usedOnHeapStorageMemory"`
		UsedOffHeapStorageMemory  int `json:"usedOffHeapStorageMemory"`
		TotalOnHeapStorageMemory  int `json:"totalOnHeapStorageMemory"`
		TotalOffHeapStorageMemory int `json:"totalOffHeapStorageMemory"`
	} `json:"memoryMetrics"`
	BlacklistedInStages []interface{} `json:"blacklistedInStages"`
	PeakMemoryMetrics   struct {
		JVMHeapMemory              int `json:"JVMHeapMemory"`
		JVMOffHeapMemory           int `json:"JVMOffHeapMemory"`
		OnHeapExecutionMemory      int `json:"OnHeapExecutionMemory"`
		OffHeapExecutionMemory     int `json:"OffHeapExecutionMemory"`
		OnHeapStorageMemory        int `json:"OnHeapStorageMemory"`
		OffHeapStorageMemory       int `json:"OffHeapStorageMemory"`
		OnHeapUnifiedMemory        int `json:"OnHeapUnifiedMemory"`
		OffHeapUnifiedMemory       int `json:"OffHeapUnifiedMemory"`
		DirectPoolMemory           int `json:"DirectPoolMemory"`
		MappedPoolMemory           int `json:"MappedPoolMemory"`
		NettyDirectMemory          int `json:"NettyDirectMemory"`
		JvmDirectMemory            int `json:"JvmDirectMemory"`
		SparkDirectMemoryOverLimit int `json:"SparkDirectMemoryOverLimit"`
		TotalOffHeapMemory         int `json:"TotalOffHeapMemory"`
		ProcessTreeJVMVMemory      int `json:"ProcessTreeJVMVMemory"`
		ProcessTreeJVMRSSMemory    int `json:"ProcessTreeJVMRSSMemory"`
		ProcessTreePythonVMemory   int `json:"ProcessTreePythonVMemory"`
		ProcessTreePythonRSSMemory int `json:"ProcessTreePythonRSSMemory"`
		ProcessTreeOtherVMemory    int `json:"ProcessTreeOtherVMemory"`
		ProcessTreeOtherRSSMemory  int `json:"ProcessTreeOtherRSSMemory"`
		MinorGCCount               int `json:"MinorGCCount"`
		MinorGCTime                int `json:"MinorGCTime"`
		MajorGCCount               int `json:"MajorGCCount"`
		MajorGCTime                int `json:"MajorGCTime"`
		TotalGCTime                int `json:"TotalGCTime"`
	} `json:"peakMemoryMetrics"`
	Attributes struct {
	} `json:"attributes"`
	Resources struct {
	} `json:"resources"`
	ResourceProfileId int           `json:"resourceProfileId"`
	IsExcluded        bool          `json:"isExcluded"`
	ExcludedInStages  []interface{} `json:"excludedInStages"`
}

type SparkStage struct {
	Status                           string        `json:"status"`
	StageId                          int           `json:"stageId"`
	AttemptId                        int           `json:"attemptId"`
	NumTasks                         int           `json:"numTasks"`
	NumActiveTasks                   int           `json:"numActiveTasks"`
	NumCompleteTasks                 int           `json:"numCompleteTasks"`
	NumFailedTasks                   int           `json:"numFailedTasks"`
	NumKilledTasks                   int           `json:"numKilledTasks"`
	NumCompletedIndices              int           `json:"numCompletedIndices"`
	PeakNettyDirectMemory            int           `json:"peakNettyDirectMemory"`
	PeakJvmDirectMemory              int           `json:"peakJvmDirectMemory"`
	PeakSparkDirectMemoryOverLimit   int           `json:"peakSparkDirectMemoryOverLimit"`
	PeakTotalOffHeapMemory           int           `json:"peakTotalOffHeapMemory"`
	SubmissionTime                   string        `json:"submissionTime"`
	FirstTaskLaunchedTime            string        `json:"firstTaskLaunchedTime"`
	CompletionTime                   string        `json:"completionTime"`
	ExecutorDeserializeTime          int           `json:"executorDeserializeTime"`
	ExecutorDeserializeCpuTime       int           `json:"executorDeserializeCpuTime"`
	ExecutorRunTime                  int           `json:"executorRunTime"`
	ExecutorCpuTime                  int           `json:"executorCpuTime"`
	ResultSize                       int           `json:"resultSize"`
	JvmGcTime                        int           `json:"jvmGcTime"`
	ResultSerializationTime          int           `json:"resultSerializationTime"`
	MemoryBytesSpilled               int           `json:"memoryBytesSpilled"`
	DiskBytesSpilled                 int           `json:"diskBytesSpilled"`
	PeakExecutionMemory              int           `json:"peakExecutionMemory"`
	InputBytes                       int           `json:"inputBytes"`
	InputRecords                     int           `json:"inputRecords"`
	OutputBytes                      int           `json:"outputBytes"`
	OutputRecords                    int           `json:"outputRecords"`
	ShuffleRemoteBlocksFetched       int           `json:"shuffleRemoteBlocksFetched"`
	ShuffleLocalBlocksFetched        int           `json:"shuffleLocalBlocksFetched"`
	ShuffleFetchWaitTime             int           `json:"shuffleFetchWaitTime"`
	ShuffleRemoteBytesRead           int           `json:"shuffleRemoteBytesRead"`
	ShuffleRemoteBytesReadToDisk     int           `json:"shuffleRemoteBytesReadToDisk"`
	ShuffleLocalBytesRead            int           `json:"shuffleLocalBytesRead"`
	ShuffleReadBytes                 int           `json:"shuffleReadBytes"`
	ShuffleReadRecords               int           `json:"shuffleReadRecords"`
	ShuffleCorruptMergedBlockChunks  int           `json:"shuffleCorruptMergedBlockChunks"`
	ShuffleMergedFetchFallbackCount  int           `json:"shuffleMergedFetchFallbackCount"`
	ShuffleMergedRemoteBlocksFetched int           `json:"shuffleMergedRemoteBlocksFetched"`
	ShuffleMergedLocalBlocksFetched  int           `json:"shuffleMergedLocalBlocksFetched"`
	ShuffleMergedRemoteChunksFetched int           `json:"shuffleMergedRemoteChunksFetched"`
	ShuffleMergedLocalChunksFetched  int           `json:"shuffleMergedLocalChunksFetched"`
	ShuffleMergedRemoteBytesRead     int           `json:"shuffleMergedRemoteBytesRead"`
	ShuffleMergedLocalBytesRead      int           `json:"shuffleMergedLocalBytesRead"`
	ShuffleRemoteReqsDuration        int           `json:"shuffleRemoteReqsDuration"`
	ShuffleMergedRemoteReqsDuration  int           `json:"shuffleMergedRemoteReqsDuration"`
	ShuffleWriteBytes                int           `json:"shuffleWriteBytes"`
	ShuffleWriteTime                 int           `json:"shuffleWriteTime"`
	ShuffleWriteRecords              int           `json:"shuffleWriteRecords"`
	Name                             string        `json:"name"`
	Description                      string        `json:"description"`
	Details                          string        `json:"details"`
	SchedulingPool                   string        `json:"schedulingPool"`
	RddIds                           []int         `json:"rddIds"`
	AccumulatorUpdates               []interface{} `json:"accumulatorUpdates"`
	KilledTasksSummary               struct {
	} `json:"killedTasksSummary"`
	ResourceProfileId   int `json:"resourceProfileId"`
	PeakExecutorMetrics struct {
		JVMHeapMemory              int `json:"JVMHeapMemory"`
		JVMOffHeapMemory           int `json:"JVMOffHeapMemory"`
		OnHeapExecutionMemory      int `json:"OnHeapExecutionMemory"`
		OffHeapExecutionMemory     int `json:"OffHeapExecutionMemory"`
		OnHeapStorageMemory        int `json:"OnHeapStorageMemory"`
		OffHeapStorageMemory       int `json:"OffHeapStorageMemory"`
		OnHeapUnifiedMemory        int `json:"OnHeapUnifiedMemory"`
		OffHeapUnifiedMemory       int `json:"OffHeapUnifiedMemory"`
		DirectPoolMemory           int `json:"DirectPoolMemory"`
		MappedPoolMemory           int `json:"MappedPoolMemory"`
		NettyDirectMemory          int `json:"NettyDirectMemory"`
		JvmDirectMemory            int `json:"JvmDirectMemory"`
		SparkDirectMemoryOverLimit int `json:"SparkDirectMemoryOverLimit"`
		TotalOffHeapMemory         int `json:"TotalOffHeapMemory"`
		ProcessTreeJVMVMemory      int `json:"ProcessTreeJVMVMemory"`
		ProcessTreeJVMRSSMemory    int `json:"ProcessTreeJVMRSSMemory"`
		ProcessTreePythonVMemory   int `json:"ProcessTreePythonVMemory"`
		ProcessTreePythonRSSMemory int `json:"ProcessTreePythonRSSMemory"`
		ProcessTreeOtherVMemory    int `json:"ProcessTreeOtherVMemory"`
		ProcessTreeOtherRSSMemory  int `json:"ProcessTreeOtherRSSMemory"`
		MinorGCCount               int `json:"MinorGCCount"`
		MinorGCTime                int `json:"MinorGCTime"`
		MajorGCCount               int `json:"MajorGCCount"`
		MajorGCTime                int `json:"MajorGCTime"`
		TotalGCTime                int `json:"TotalGCTime"`
	} `json:"peakExecutorMetrics"`
	IsShufflePushEnabled bool `json:"isShufflePushEnabled"`
	ShuffleMergersCount  int  `json:"shuffleMergersCount"`
}

func GetSparkApplications(url string, headers map[string]string) (ApplicationsResponse, error) {
	req, _ := http.NewRequest("GET", url+"/api/v1/applications", nil)

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}

	body, _ := io.ReadAll(resp.Body)

	var apps ApplicationsResponse
	if err := json.Unmarshal(body, &apps); err != nil {
		log.Printf("Error: %s", err)
	}
	return apps, nil
}
