package utils

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

		var metricValue float64
		var metricName string
		mtime := time.Now()

		if e.Field(i).Kind() == reflect.Struct {
			log.Trace("encountered nested structure ")
			out = append(out, CreateMetricModels(prefix, e.Field(i), tags)...)
			continue
		}

		metricName = prefix + e.Type().Field(i).Name

		switch n := e.Field(i).Interface().(type) {
		case int:
			metricValue = float64(n)
		case int64:
			metricValue = float64(n)
		case uint64:
			metricValue = float64(n)
		case float64:
			metricValue = n
		case bool:
			metricValue = float64(0)
			if n {
				metricValue = float64(1)
			}
		default:
			log.Trace("setMetrics :skipping metric: ", n, metricName, metricValue)
		}

		meltMetric := model.MakeGaugeMetric(
			metricName, model.Numeric{FltVal: metricValue}, mtime)

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

	if e.Kind() == reflect.Interface || e.Kind() == reflect.Ptr {
		e = e.Elem()
	}

	log.Println(e)
	switch e.Kind() {
	case reflect.Struct:
		for i := 0; i < e.NumField(); i++ {
			var mname string
			mname = prefix + e.Type().Field(i).Name
			switch n := e.Field(i).Interface().(type) {
			case string:
				log.Trace("setTags : adding tags ", mname, "=", n)
				metricTags[mname] = n
			case []int:
				metricTags[mname] = SplitToString(n, ",")
			default:
				// Handle other cases if needed
			}
		}
	case reflect.Map:
		for _, key := range e.MapKeys() {
			var mname string
			mname = prefix + key.String()
			val := e.MapIndex(key)
			switch n := val.Interface().(type) {
			case string:
				log.Trace("setTags : adding tags ", mname, "=", n)
				metricTags[mname] = n
			case []int:
				metricTags[mname] = SplitToString(n, ",")
			default:
				// Handle other cases if needed
			}
		}
	default:
		log.Println("Unsupported kind:", e.Kind())
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
