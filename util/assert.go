package util

import "fmt"

func AssertMapHas(m map[string]interface{}, key string) string {
	if v, ok := m[key]; !ok || v == nil {
		return fmt.Sprintf(" key `%s` must has; ", key)
	}
	return ""
}
