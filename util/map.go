package util

import (
	"fmt"
	"strconv"
)

func GetMapString(m map[string]interface{}, key string, defaultValue string) string {
	if v, ok := m[key]; ok {
		return interfaceToString(v, defaultValue)
	}
	return defaultValue
}

func interfaceToString(v interface{}, defaultValue string) string {
	switch v.(type) {
	case string:
		return v.(string)
	case int:
		return strconv.Itoa(v.(int))
	case int32:
		return strconv.FormatInt(v.(int64), 32)
	case int64:
		return strconv.FormatInt(v.(int64), 64)
	case float32:
		return strconv.FormatFloat(v.(float64), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v.(float64), 'f', -1, 64)
	default:
		return defaultValue
	}
}

func GetMapStringSlice(m map[string]interface{}, key string, defaultValue []string) []string {
	if v, ok := m[key]; ok {
		switch v.(type) {
		case []interface{}:
			interfaces := v.([]interface{})
			res := make([]string, len(interfaces))
			for i, s := range interfaces {
				res[i] = interfaceToString(s, "")
			}
			return res
		case []string:
			return v.([]string)
		}
	}
	return defaultValue
}

func GetMapInt(m map[string]interface{}, key string, defaultValue int64) int64 {
	if v, ok := m[key]; ok {
		switch v.(type) {
		case int:
			return v.(int64)
		case int32:
			return v.(int64)
		case int64:
			return v.(int64)
		case float32:
			return v.(int64)
		case float64:
			return v.(int64)
		default:
			return defaultValue
		}
	}
	return defaultValue
}

// MergeNewMap 融合两个map,返回一个新map
func MergeNewMap(master, follow map[string]interface{}) map[string]interface{} {
	return MergeMap(master, MergeMap(follow, make(map[string]interface{}, len(master)+len(follow))))
}

func cleanupInterfaceArray(in []interface{}) []interface{} {
	res := make([]interface{}, len(in))
	for i, v := range in {
		res[i] = CleanupMapValue(v)
	}
	return res
}

func cleanupInterfaceMap(in map[interface{}]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range in {
		res[fmt.Sprintf("%v", k)] = CleanupMapValue(v)
	}
	return res
}

func CleanupMapValue(v interface{}) interface{} {
	if v == nil {
		return v
	}
	switch v := v.(type) {
	case []interface{}:
		return cleanupInterfaceArray(v)
	case map[interface{}]interface{}:
		return cleanupInterfaceMap(v)
	case map[string]interface{}:
		return CleanupStringMap(v)
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

func CleanupStringMap(in map[string]interface{}) map[string]interface{} {
	for k, v := range in {
		in[k] = CleanupMapValue(v)
	}
	return in
}
