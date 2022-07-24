package toml

import (
	"encoding/json"
	"github.com/izern/zf/types"
	"github.com/pelletier/go-toml/v2"
)

func init() {

}

type TomlCodec struct {
	json.Decoder
}

func (y *TomlCodec) Marshal(m interface{}) (result []byte, err types.ZfError) {
	mType, err := types.GetType(m)
	if err != nil {
		return nil, err
	}
	var e error
	switch mType {
	case types.Object:
		result, e = toml.Marshal(m)
	default:
		result, e = json.Marshal(m)
	}
	if e != nil {
		return nil, types.NewFormatError(e.Error(), "toml")
	}
	return result, err
}

func (y *TomlCodec) Unmarshal(b []byte) (result interface{}, err types.ZfError) {
	result = make(map[interface{}]interface{})
	e := toml.Unmarshal(b, &result)
	if e != nil {
		result = make([]interface{}, 0)
		e = json.Unmarshal(b, &result)
		if e != nil {
			return string(b), nil
		}
	}
	return result, err
}
