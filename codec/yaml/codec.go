package yaml

import (
	"github.com/izern/zf/types"
	"gopkg.in/yaml.v3"
)

func init() {

}

type YamlCodec struct {
}

func (y *YamlCodec) Marshal(m interface{}) (result []byte, err types.ZfError) {

	result, e := yaml.Marshal(&m)
	if e != nil {
		return nil, types.NewFormatError(e.Error(), "yaml")
	}
	return result, err
}

func (y *YamlCodec) Unmarshal(b []byte) (result interface{}, err types.ZfError) {
	var tmp = make(map[string]interface{})
	e := yaml.Unmarshal(b, &tmp)
	if e != nil {
		e = yaml.Unmarshal(b, &result)
		if e != nil {
			return nil, types.NewFormatError(e.Error(), "yaml")
		}
	} else {
		result = tmp
	}
	return result, err
}
