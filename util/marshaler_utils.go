package util

import (
	"fmt"
	"github.com/izern/zf/types"
)

func init() {

}

// ConvertArray2String converts an array with interface{} keys to string keys
// Optimized to reduce type checking overhead
func ConvertArray2String(param []interface{}) []interface{} {
	if len(param) == 0 {
		return param
	}

	result := make([]interface{}, len(param))
	for i, item := range param {
		result[i] = convertValue(item)
	}
	return result
}

// ConvertMap2String converts a map with interface{} keys to string keys  
// Optimized to avoid redundant type checking and allocations
func ConvertMap2String(m map[interface{}]interface{}) map[string]interface{} {
	if len(m) == 0 {
		return make(map[string]interface{})
	}

	res := make(map[string]interface{}, len(m))
	for k, v := range m {
		key := fmt.Sprint(k)
		res[key] = convertValue(v)
	}
	return res
}

// convertValue is a helper function to convert a single value
// This reduces code duplication between array and map conversion
func convertValue(v interface{}) interface{} {
	switch val := v.(type) {
	case map[interface{}]interface{}:
		return ConvertMap2String(val)
	case map[string]interface{}:
		return val // Already correct type
	case []interface{}:
		return ConvertArray2String(val)
	case []map[interface{}]interface{}:
		result := make([]map[string]interface{}, len(val))
		for i, item := range val {
			result[i] = ConvertMap2String(item)
		}
		return result
	case []map[string]interface{}:
		return val // Already correct type
	default:
		return v // Primitive types don't need conversion
	}
}

// BatchConvertMaps efficiently converts multiple maps at once
func BatchConvertMaps(maps []map[interface{}]interface{}) []map[string]interface{} {
	if len(maps) == 0 {
		return nil
	}
	
	result := make([]map[string]interface{}, len(maps))
	for i, m := range maps {
		result[i] = ConvertMap2String(m)
	}
	return result
}

// IsConversionNeeded checks if a value needs type conversion
func IsConversionNeeded(v interface{}) bool {
	switch v.(type) {
	case map[interface{}]interface{}:
		return true
	case []map[interface{}]interface{}:
		return true
	case []interface{}:
		// Check if any element needs conversion
		if arr, ok := v.([]interface{}); ok {
			for _, item := range arr {
				if IsConversionNeeded(item) {
					return true
				}
			}
		}
		return false
	default:
		return false
	}
}
