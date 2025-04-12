package utils

import (
	"reflect"
	"strconv"
)

func GetTypeName(t interface{}) string {
	return reflect.TypeOf(t).String()
}

// GetStringFromMap safely gets a string value from a map
func GetStringFromMap(data map[string]interface{}, key string, defaultValue string) string {
	if val, ok := data[key]; ok {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	return defaultValue
}

// GetIntFromMap safely gets an int value from a map
func GetIntFromMap(data map[string]interface{}, key string, defaultValue int) int {
	if val, ok := data[key]; ok {
		switch v := val.(type) {
		case int:
			return v
		case int64:
			return int(v)
		case float64:
			return int(v)
		case string:
			i, err := strconv.Atoi(v)
			if err == nil {
				return i
			}
		}
	}
	return defaultValue
}

// GetFloatFromMap safely gets a float32 value from a map
func GetFloatFromMap(data map[string]interface{}, key string, defaultValue float32) float32 {
	if val, ok := data[key]; ok {
		switch v := val.(type) {
		case float32:
			return v
		case float64:
			return float32(v)
		case int:
			return float32(v)
		case int64:
			return float32(v)
		case string:
			f, err := strconv.ParseFloat(v, 32)
			if err == nil {
				return float32(f)
			}
		}
	}
	return defaultValue
}