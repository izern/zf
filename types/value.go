package types

import (
	"reflect"
)

type ValueType string

const (
	Array  ValueType = "array"
	Object           = "object"
	Number           = "number"
	String           = "string"
	Bool             = "bool"
	Null             = "null"
)

func GetType(v interface{}) (ValueType, ZfError) {
	switch v.(type) {
	case bool:
		return Bool, nil
	case nil:
		return Null, nil
	case map[string]interface{}:
		return Object, nil
	case []interface{}, []map[string]interface{}:
		return Array, nil
	case string:
		return String, nil
	case byte, int, uint, int8, int16, uint16, int32, uint32, int64, uint64, float32, float64, complex64, complex128:
		return Number, nil
	default:
		return Null, NewUnSupportError(reflect.TypeOf(v).Name())
	}
}
