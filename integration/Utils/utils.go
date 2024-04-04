package Utils

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"newrelic/multienv/pkg/model"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func StringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func Flatten(input map[string]interface{}, currentKey string, result map[string]interface{}) map[string]interface{} {
	for key, value := range input {
		newKey := key
		if currentKey != "" {
			newKey = currentKey + "." + key
		}

		// check if the value is another nested map
		if nestedMap, ok := value.(map[string]interface{}); ok {
			result = Flatten(nestedMap, newKey, result)
		} else {
			result[newKey] = value
		}
	}
	return result
}

func ConvertJsonToMap(jsonData []byte) (map[string]interface{}, error) {
	var JSONObject map[string]interface{}
	err := json.Unmarshal(jsonData, &JSONObject)
	if err != nil {
		return nil, err
	}
	return JSONObject, nil
}

func CreateMetricModels(prefix string, e reflect.Value, tags map[string]interface{}) []model.MeltModel {

	out := make([]model.MeltModel, 0)

	for i := 0; i < e.NumField(); i++ {

		var mvalue float64
		var mname string
		mtime := time.Now()

		if e.Field(i).Kind() == reflect.Struct {
			log.Trace("encountered nested structure ")
			out = append(out, CreateMetricModels(prefix, e.Field(i), tags)...)
			continue
		}

		mname = prefix + e.Type().Field(i).Name

		switch n := e.Field(i).Interface().(type) {
		case int:
			mvalue = float64(n)
		case int64:
			mvalue = float64(n)
		case uint64:
			mvalue = float64(n)
		case float64:
			mvalue = n
		case bool:
			mvalue = float64(0)
			if n {
				mvalue = float64(1)
			}
		default:
			log.Trace("setMetrics :skipping metric: ", n, mname, mvalue)
		}

		meltMetric := model.MakeGaugeMetric(
			mname, model.Numeric{FltVal: mvalue}, mtime)

		tags["instrumentation.name"] = "newrelic-databricks-integration"
		meltMetric.Attributes = tags

		out = append(out, meltMetric)
	}

	return out

}

func SetTags(prefix string, e reflect.Value, tags map[string]interface{}, metricTags map[string]interface{}) {

	for k, v := range tags {
		metricTags[k] = v
	}

	//log.Println(e)
	for i := 0; i < e.NumField(); i++ {

		var mname string
		mname = prefix + e.Type().Field(i).Name
		//log.Println("MNAME ", mname)
		// populate string metrics as tags
		switch n := e.Field(i).Interface().(type) {
		case string:
			// Add this in tags
			log.Trace("setTags : adding tags ", mname, "=", n)
			metricTags[mname] = n
		case []int:
			metricTags[mname] = SplitToString(n, ",")
		default:
			//log.Debug("setTags :Skipping tags")
		}

	}
}

func SplitToString(a []int, sep string) string {
	if len(a) == 0 {
		return ""
	}

	b := make([]string, len(a))
	for i, v := range a {
		b[i] = strconv.Itoa(v)
	}
	return strings.Join(b, sep)
}
