package deser

import (
	"encoding/json"
	"fmt"
)

// JSON deserializer
//func DeserJson(data []byte, v *map[string]any) error {
//	err := json.Unmarshal(data, v)
//	if err != nil {
//		log.Error("Deser error = ", err.Error())
//	}
//	return err
//}

func DeserJson(data []byte) (any, string, error) {
	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err == nil {
		return v, "object", nil
	}

	var a []map[string]interface{}
	if err := json.Unmarshal(data, &a); err == nil {
		return a, "array", nil
	}

	return nil, "", fmt.Errorf("unknown type")

}
