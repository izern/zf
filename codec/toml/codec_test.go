package toml

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {

}

func Test_Unmarshal(t *testing.T) {
	var data = `
[owner]
name = "Tom Preston-Werner"
dob = 1979-05-27T07:32:00-08:00 # First class dates
`
	codec := &TomlCodec{}
	yamlMap, err := codec.Unmarshal([]byte(data))
	if err != nil {
		panic(err)
	}
	fmt.Println(yamlMap)
	switch yamlMap.(type) {
	case map[string]interface{}:
		fmt.Println("map[string]interface{}")
	case map[interface{}]interface{}:
		fmt.Println("map[interface{}]interface{}")
	}

}

func TestYamlCodec_Marshal(t *testing.T) {
	codec := &TomlCodec{}

	marshal, err := codec.Marshal(map[string]interface{}{
		"a": "Easy",
		"b": map[string]interface{}{
			"c": 2000000,
			"d": true},
	})

	assert.Nil(t, err)
	res, _ := codec.Unmarshal(marshal)
	fmt.Println(res)
	marshal, err = codec.Marshal(map[interface{}]interface{}{
		"a": "Easy",
		"b": map[interface{}]interface{}{
			"c": 2000000,
			"d": true},
	})
	assert.Nil(t, err)
	res, _ = codec.Unmarshal(marshal)
	fmt.Println(res)

}
