package cmd

import (
	"errors"
	"fmt"
	"github.com/izern/zf/types"
)

var extCommandMap = make(map[string]types.TypeCommand)

func Regist(cmd types.TypeCommand) {
	name := cmd.GetCurrType()
	if _, ok := extCommandMap[name]; ok {
		panic("内部错误,重复注册:" + name)
	}
	extCommandMap[name] = cmd
}

func GetCmd(name string) (types.TypeCommand, error) {
	if cmd, ok := extCommandMap[name]; ok {
		return cmd, nil
	}
	return nil, errors.New(fmt.Sprintf("%s未注册，无法使用", name))
}

func GetAllCmd() []types.TypeCommand {
	var result []types.TypeCommand
	for _, v := range extCommandMap {
		result = append(result, v)
	}
	return result
}
