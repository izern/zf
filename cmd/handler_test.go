package cmd

import (
	"github.com/izern/zf/codec/json"
	"github.com/izern/zf/codec/toml"
	yaml2 "github.com/izern/zf/codec/yaml"
	"github.com/izern/zf/test"
	"github.com/izern/zf/util"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func init() {

}

var text = test.TestYaml

var yamlHandler = NewHandler(&yaml2.YamlCodec{}, &yaml2.YamlCodec{}, "yaml")

var handlers = [3]*Handler{
	NewHandler(&yaml2.YamlCodec{}, &yaml2.YamlCodec{}, "yaml"),
	NewHandler(&json.JSONCodec{}, &json.JSONCodec{}, "json"),
	NewHandler(&toml.TomlCodec{}, &toml.TomlCodec{}, "toml"),
}

// 统一提前转换成map格式
func Before(t *testing.T) interface{} {
	res, e := yamlHandler.GetValues(0, math.MaxUint32, ".", text)
	assert.Nil(t, e, "parse text to yaml failed")
	assert.NotNil(t, res, "parse text to yaml failed")
	switch res.(type) {
	case map[string]interface{}:
		res = res.(map[string]interface{})
	case map[interface{}]interface{}:
		res = util.ConvertMap2String(res.(map[interface{}]interface{}))
	default:
		assert.Fail(t, "not map", res)
	}
	return res
}

func TestAppend(t *testing.T) {
	param := Before(t)

	var m = `- cipher: aes-256-gcm
  name: '[07-23]|oslook|斯洛伐克(SK)Slovakia/Bratislava_1'
  password: TEzjfAYq2IjtuoS
  port: 6697
  server: 141.164.39.146
  type: ss
`
	tmpValue, e := yamlHandler.Unmarshaler.Unmarshal([]byte(m))
	assert.Nil(t, e, "parse text to yaml failed")
	assert.NotNil(t, tmpValue, "parse text to yaml failed")

	for _, handler := range handlers {
		str, e := handler.Marshal(param)
		assert.Nil(t, e, "marshal param failed.", param)
		result, err := handler.Append(".rules", "", 0, "1234", str)
		assert.Nil(t, err, ".rules append value.", handler, str)
		assert.NotNil(t, result, ".rules append value.", handler, str)

		result, err = handler.Append(".proxies[0]", "password", 0, "1234", str)
		assert.Nil(t, err, ".proxies[0] append value.", handler, str)
		assert.NotNil(t, result, ".proxies[0] append value.", handler, str)

		marshal, zfError := handler.Marshal(tmpValue)
		assert.Nil(t, zfError, "Marshal map", handler, str)
		assert.NotNil(t, marshal, "Marshal map", handler, str)

		result, err = handler.Append(".proxies", "", 0, marshal, str)
		assert.Nil(t, err, ".proxies append value.", handler, str)
		assert.NotNil(t, result, ".proxies append value.", handler, str)

	}
}

func Test_GetValues(t *testing.T) {

	param := Before(t)

	for _, handler := range handlers {
		str, e := handler.Marshal(param)
		assert.Nil(t, e, "marshal param failed.", param)
		result, err := handler.GetValues(0, math.MaxUint32, ".proxies", str)
		assert.Nil(t, err, ".proxies getValues error", handler, str)
		assert.NotNil(t, result, ".proxies getValues error", handler, str)

		result, err = handler.GetValues(0, math.MaxUint32, ".proxies[0]", str)
		assert.Nil(t, err, ".proxies[0] getValues error", handler, str)
		assert.NotNil(t, result, ".proxies[0] getValues error", handler, str)

	}
}

func Test_SetValues(t *testing.T) {
	param := Before(t)
	var m = `- cipher: aes-256-gcm
  name: '[07-23]|oslook|斯洛伐克(SK)Slovakia/Bratislava_1'
  password: TEzjfAYq2IjtuoS
  port: 6697
  server: 141.164.39.146
  type: ss
`
	tmpValue, e := yamlHandler.Unmarshaler.Unmarshal([]byte(m))
	assert.Nil(t, e, "parse text to yaml failed")
	assert.NotNil(t, tmpValue, "parse text to yaml failed")

	for _, handler := range handlers {
		str, e := handler.Marshal(param)
		assert.Nil(t, e, "marshal param failed.", param)

		result, err := handler.SetValue(".rules[1,2]", "12345", str)
		assert.Nil(t, err, ".rules[1,2] getValues error", handler, str)
		assert.NotNil(t, result, ".rules[1,2] getValues error", handler, str)

		marshal, zfError := handler.Marshal(tmpValue)
		assert.Nil(t, zfError, "Marshal map", handler, str)
		assert.NotNil(t, marshal, "Marshal map", handler, str)
		result, err = handler.SetValue(".proxies[1,2]", marshal, str)
		assert.Nil(t, err, ".proxies[1,2] getValues error", handler, str)
		assert.NotNil(t, result, ".proxies[1,2] getValues error", handler, str)

		result, err = handler.SetValue(".proxies[1,2].name", "54321", str)
		assert.Nil(t, err, ".proxies[1,2] getValues error", handler, str)
		assert.NotNil(t, result, ".proxies[1,2] getValues error", handler, str)

	}

}
