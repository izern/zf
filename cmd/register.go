package cmd

import (
	"github.com/izern/zf/codec/json"
	"github.com/izern/zf/codec/toml"
	yaml2 "github.com/izern/zf/codec/yaml"
)

func init() {

	yamlHandler := NewHandler(&yaml2.YamlCodec{}, &yaml2.YamlCodec{}, "yaml")
	jsonHandler := NewHandler(&json.JSONCodec{}, &json.JSONCodec{}, "json")
	tomlHandler := NewHandler(&toml.TomlCodec{}, &toml.TomlCodec{}, "toml")

	Regist(yamlHandler)
	Regist(jsonHandler)
	Regist(tomlHandler)

}
