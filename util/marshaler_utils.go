package util

import (
	"fmt"
	"github.com/izern/zf/types"
)

func init() {

}

func ConvertArray2String(param []interface{}) (result []interface{}) {

	result = make([]interface{}, 0)
	for _, item := range param {
		itemType, _ := types.GetType(item)
		switch itemType {
		case types.Object:
			switch item.(type) {
			case map[interface{}]interface{}:
				result = append(result, ConvertMap2String(item.(map[interface{}]interface{})))
			case map[string]interface{}:
				result = append(result, item.(map[string]interface{}))
			}
		case types.Array:
			res := ConvertArray2String(item.([]interface{}))
			result = append(result, res)
		default:
			result = append(result, item)
		}
	}
	return result
}

func ConvertMap2String(m map[interface{}]interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	for k, v := range m {
		vType, _ := types.GetType(v)
		switch vType {
		case types.Object:
			switch v.(type) {
			case map[interface{}]interface{}:
				res[fmt.Sprint(k)] = ConvertMap2String(v.(map[interface{}]interface{}))
			case map[string]interface{}:
				res[fmt.Sprint(k)] = v.(map[string]interface{})
			}
		case types.Array:
			switch v.(type) {
			case []interface{}:
				v2 := v.([]interface{})
				res[fmt.Sprint(k)] = ConvertArray2String(v2)
			case []map[interface{}]interface{}:
				v2 := v.([]map[interface{}]interface{})
				var v2Res []map[string]interface{}
				for _, tmpV := range v2 {
					v2Res = append(v2Res, ConvertMap2String(tmpV))
				}
				res[fmt.Sprint(k)] = v2Res
			case []map[string]interface{}:
				v2 := v.([]map[string]interface{})
				res[fmt.Sprint(k)] = v2
			}
		default:
			res[fmt.Sprint(k)] = v
		}
	}
	return res
}

//func ConvertArray2Interface(param []interface{}) (result []interface{}) {
//
//	result = make([]interface{}, 0)
//	for _, item := range param {
//		itemType, _ := types.GetType(item)
//		switch itemType {
//		case types.Object:
//			switch item.(type) {
//			case map[interface{}]interface{}:
//				result = append(result, ConvertMap2String(item.(map[interface{}]interface{})))
//			case map[string]interface{}:
//				result = append(result, item.(map[string]interface{}))
//			}
//		case types.Array:
//			res := ConvertArray2Interface(item.([]interface{}))
//			result = append(result, res)
//		default:
//			result = append(result, item)
//		}
//	}
//	return result
//}

//func ConvertMap2Interface(m map[string]interface{}) map[interface{}]interface{} {
//	res := map[interface{}]interface{}{}
//	for k, v := range m {
//		vType, _ := types.GetType(v)
//		switch vType {
//		case types.Object:
//			res[k] = ConvertMap2Interface(v.(map[string]interface{}))
//		case types.Array:
//			switch v.(type) {
//			case []map[interface{}]interface{}:
//				v2 := v.([]map[interface{}]interface{})
//				res[k] = v2
//			case []interface{}:
//				v2 := v.([]interface{})
//				v2Res := ConvertArray2Interface(v2)
//				res[k] = v2Res
//			case []map[string]interface{}:
//				v2 := v.([]map[string]interface{})
//				var v2Res []map[interface{}]interface{}
//				for _, tmpV := range v2 {
//					v2Res = append(v2Res, ConvertMap2Interface(tmpV))
//				}
//				res[k] = v2Res
//
//			}
//		default:
//			res[k] = v
//		}
//	}
//	return res
//}
