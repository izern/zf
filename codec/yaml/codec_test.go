package yaml

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {

}

func Test_Unmarshal(t *testing.T) {
	var data = `
a: Easy!
b:
  c: 22222222222
  d: [3, 4]
`
	codec := &YamlCodec{}
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
	codec := &YamlCodec{}

	marshal, err := codec.Marshal(map[string]interface{}{
		"a": "Easy",
		"b": map[string]interface{}{
			"c": 2000000,
			"d": true},
	})

	assert.Nil(t, err)
	fmt.Println(marshal)
	marshal, err = codec.Marshal(map[interface{}]interface{}{
		"a": "Easy",
		"b": map[interface{}]interface{}{
			"c": 2000000,
			"d": true},
	})
	assert.Nil(t, err)
	fmt.Println(marshal)

}
