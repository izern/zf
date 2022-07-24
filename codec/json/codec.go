package json

import (
	"encoding/json"
	"github.com/izern/zf/types"
)

func init() {

}

type JSONCodec struct {
}

func (y *JSONCodec) Marshal(m interface{}) (result []byte, err types.ZfError) {
	result, e := json.Marshal(&m)
	if e != nil {
		return nil, types.NewFormatError(e.Error(), "json")
	}
	return result, err
}

func (y *JSONCodec) Unmarshal(b []byte) (result interface{}, err types.ZfError) {
	if len(b) == 0 || b == nil {
		return nil, nil
	}
	e := json.Unmarshal(b, &result)
	if e != nil {
		return nil, types.NewFormatError(e.Error(), "json")
	}
	return result, err
}
