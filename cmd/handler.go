package cmd

import (
	"fmt"
	"math"
	"sort"

	"github.com/izern/zf/codec"
	"github.com/izern/zf/types"
	"github.com/izern/zf/util"
)

func init() {

}

type Codec interface {
}

type Handler struct {
	Marshaler   codec.Marshaler
	Unmarshaler codec.Unmarshaler
	Type        string
	Value       map[string]interface{}
}

func NewHandler(marshaler codec.Marshaler, unmarshaler codec.Unmarshaler, typeStr string) *Handler {
	return &Handler{
		Marshaler:   marshaler,
		Unmarshaler: unmarshaler,
		Type:        typeStr,
	}
}

func (receiver Handler) GetCurrType() string {
	return receiver.Type
}

func (receiver *Handler) parse(text string) types.ZfError {

	result, zfError := receiver.Unmarshaler.Unmarshal([]byte(text))
	if zfError != nil {
		return zfError
	}
	switch result.(type) {
	case map[interface{}]interface{}:
		receiver.Value = util.ConvertMap2String(result.(map[interface{}]interface{}))
	case map[string]interface{}:
		receiver.Value = result.(map[string]interface{})
	default:
		receiver.Value = make(map[string]interface{}, 1)
		receiver.Value[""] = result
	}
	return nil
}

func (receiver *Handler) Parse(text string) (string, types.ZfError) {
	err := receiver.parse(text)
	if err != nil {
		return "", err
	}
	return receiver.PrintToString()
}

func (receiver *Handler) PrintToString() (string, types.ZfError) {
	if receiver.Value == nil {
		return "", types.NewUnSupportError("未初始化，无法输出内容")
	}
	if len(receiver.Value) == 1 && receiver.Value[""] != nil {
		// 这种情况不是object，直接打印
		bytes, err := receiver.Marshaler.Marshal(receiver.Value[""])
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	} else {
		//bytes, err := receiver.Marshaler.Marshal(util.ConvertMap2String(receiver.Value))
		bytes, err := receiver.Marshaler.Marshal(receiver.Value)
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	}
}

func (receiver *Handler) Marshal(content interface{}) (string, types.ZfError) {
	if content != nil {
		res, err := receiver.Marshaler.Marshal(content)
		if err != nil {
			return "", err
		}
		return string(res), nil
	}
	return "", nil
}

func (receiver *Handler) Keys(from uint, to uint, path string, text string) ([]string, types.ZfError) {

	value, err := receiver.getValues(path, text)
	if err != nil {
		return nil, err
	}
	valueType, _ := types.GetType(value)
	if valueType != types.Object {
		return nil, types.NewUnSupportError("只有object支持此操作,当前类型:" + string(valueType))
	}
	v := value.(map[string]interface{})
	i := 0
	keys := make([]string, len(v))
	for k := range v {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	start := util.Max(0, int(from))
	end := util.Min(int(to), len(keys))

	return keys[start:end], nil
}

func (receiver *Handler) GetType(path string, text string) (types.ValueType, types.ZfError) {
	value, err := receiver.getValues(path, text)
	if err != nil {
		return types.Null, err
	}
	return types.GetType(value)
}

func (receiver *Handler) getValues(path string, text string) (interface{}, types.ZfError) {
	err := receiver.parse(text)
	if err != nil {
		return nil, err
	}
	paths, err := util.ParsePath(path)
	if err != nil {
		return nil, err
	}

	if paths[0].Type != types.RootNode {
		return nil, types.NewFormatError(path, "path")
	}
	return getValues(paths[1:], receiver.Value)
}

// 根据path解析值
func getValues(paths []*types.Path, v interface{}) (result interface{}, e types.ZfError) {
	result = v
	for i := 0; i < len(paths); i++ {
		p := paths[i]
		valueType, err := types.GetType(result)
		if err != nil {
			return nil, err
		}
		// 先前置处理path的key
		if valueType == types.Object {
			obj := result.(map[string]interface{})
			if childV, ok := obj[p.NodeKey]; ok {
				result = childV
			} else {
				return nil, types.NewKeyNotFoundError(p.NodeKey)
			}
		}
		valueType, err = types.GetType(result)
		if err != nil {
			return nil, err
		}
		// 根据path类型取值
		// 校验读取到的值是否满足要求
		if p.Type.IsSupportValue(valueType) {
			switch p.Type {
			case types.NormalNode:
				// 处理[].name情况，取出来数组，
				if i > 0 && (paths[i-1].Type == types.IndexNode || paths[i-1].Type == types.RangeNode) {
					array := result.([]interface{})
					tmpResult := make([]interface{}, 0)
					for _, item := range array {
						itemType, _ := types.GetType(item)
						if itemType == types.Object {
							tmpResult = append(tmpResult, item.(map[string]interface{})[p.NodeKey])
						} else {
							return nil, types.NewUnSupportError(fmt.Sprintf("当前节点下数组下元素类型是%s,不支持%s", itemType, p.OriginValue))
						}
					}
					result = tmpResult
				}
				// 如果是object类型,已经前置处理过了
			case types.IndexNode:
				array := result.([]interface{})
				if int(p.Index) >= len(array) {
					return nil, types.NewIndexOutOfBoundError(array, "array", int(p.Index))
				}
				result = array[int(p.Index)]
			case types.RangeNode:
				array := result.([]interface{})
				if p.To == math.MaxInt16 {
					p.To = uint(len(array))
				}
				if int(p.To) > len(array) {
					return nil, types.NewIndexOutOfBoundError(array, p.OriginValue, int(p.To))
				}
				for i := p.From; i < p.To; i++ {

				}
				result = array[int(p.From):int(p.To)]
			}
		} else {
			return nil, types.NewUnSupportError(fmt.Sprintf("当前值类型是%s,不支持%s", valueType, p.OriginValue))
		}
	}
	return result, nil
}

func (receiver *Handler) GetValues(from uint, to uint, path string, text string) (interface{}, types.ZfError) {

	res, err := receiver.getValues(path, text)
	if err != nil {
		return nil, err
	}

	valueType, err := types.GetType(res)
	if err != nil {
		return nil, err
	}
	if valueType == types.Array {
		switch res.(type) {
		case []interface{}:
			result := res.([]interface{})
			start := util.Max(0, int(from))
			end := util.Min(int(to), len(result))
			return result[start:end], nil
		case []map[string]interface{}:
			result := res.([]map[string]interface{})
			start := util.Max(0, int(from))
			end := util.Min(int(to), len(result))
			return result[start:end], nil
		case []map[interface{}]interface{}:
			result := res.([]map[interface{}]interface{})
			start := util.Max(0, int(from))
			end := util.Min(int(to), len(result))
			return result[start:end], nil
		}
	}

	return res, nil
}

func (receiver *Handler) Append(path string, key string, index uint, value string, text string) (string, types.ZfError) {
	err := receiver.parse(text)
	if err != nil {
		return "", err
	}
	paths, err := util.ParsePath(path)
	if err != nil {
		return "", err
	}
	if len(paths) <= 1 {
		return "", types.NewUnSupportError("路径最少要有两层，如 .a")
	}

	if paths[0].Type != types.RootNode {
		return "", types.NewFormatError(path, "path")
	}

	// 根据路径获取其父节点值
	parentValue, e := getValues(paths[1:len(paths)-1], receiver.Value)
	if e != nil {
		return "", nil
	}

	parentMap := parentValue.(map[string]interface{})
	lastPath := paths[len(paths)-1]

	lastPathV := parentMap[lastPath.NodeKey]
	lastPathVType, _ := types.GetType(lastPathV)

	v, err := receiver.Unmarshaler.Unmarshal([]byte(value))
	if err != nil {
		return "", err
	}

	vType, _ := types.GetType(v)

	switch lastPathVType {
	case types.Array:
		lastPathArrayV := lastPathV.([]interface{})
		actualIndex := util.Min(len(lastPathArrayV), int(index))

		size := 1
		var appendV []interface{}
		if types.Array == vType {
			appendV = v.([]interface{})
			size = len(appendV)
		} else {
			appendV = []interface{}{v}
		}
		result := make([]interface{}, len(lastPathArrayV)+size)

		err = util.ArrayCopy(lastPathArrayV, 0, result, 0, actualIndex)
		if err != nil {
			return "", err
		}
		err = util.ArrayCopy(appendV, 0, result, actualIndex, size)
		if err != nil {
			return "", err
		}

		err = util.ArrayCopy(lastPathArrayV, actualIndex, result, actualIndex+size, len(lastPathArrayV)-actualIndex)
		if err != nil {
			return "", err
		}

		parentMap[lastPath.NodeKey] = result

	case types.Object:
		lastPathMapV := lastPathV.(map[interface{}]interface{})

		var vMap map[string]interface{}
		switch v.(type) {
		case map[string]interface{}:
			vMap = v.(map[string]interface{})
		case map[interface{}]interface{}:
			vMap = util.ConvertMap2String(v.(map[interface{}]interface{}))
		}
		if vType == types.Object {
			for k, vItem := range vMap {
				lastPathMapV[k] = vItem
			}
		} else {
			if key == "" {
				return "", types.NewUnSupportError("当前节点类别为object，必须指定key")
			}
			lastPathMapV[key] = v
		}
	case types.Null:
		parentMap[lastPath.NodeKey] = v
	default:
		parentMap[lastPath.NodeKey] = fmt.Sprintf("%v%v", lastPathV, v)
	}

	return receiver.PrintToString()
}

func (receiver *Handler) SetValue(path string, value string, text string) (string, types.ZfError) {

	err := receiver.parse(text)
	if err != nil {
		return "", err
	}
	paths, err := util.ParsePath(path)
	if err != nil {
		return "", err
	}
	if len(paths) < 1 {
		return "", types.NewUnSupportError("路径最少要有两层，如 .a")
	}

	if paths[0].Type != types.RootNode {
		return "", types.NewFormatError(path, "path")
	}

	// 根据路径获取其父节点值
	parentValue, e := getValues(paths[1:len(paths)-1], receiver.Value)
	if e != nil {
		return "", e
	}
	v, err := receiver.Unmarshaler.Unmarshal([]byte(value))
	if err != nil {
		return "", err
	}

	e = receiver.setValue0(parentValue, v, 0, 0, paths)
	if e != nil {
		return "", e
	}

	return receiver.PrintToString()
}

func (receiver *Handler) setValue0(parentV interface{}, v interface{}, from uint, to uint, paths []*types.Path) types.ZfError {
	switch parentV.(type) {
	case map[string]interface{}:
		parentMap := parentV.(map[string]interface{})
		lastPath := paths[len(paths)-1]

		lastPathV := parentMap[lastPath.NodeKey]

		switch lastPath.Type {
		case types.IndexNode:
			err := receiver.setValue1(lastPathV, v, lastPath.Index, lastPath.Index+1)
			if err != nil {
				return err
			}
		case types.RangeNode:
			err := receiver.setValue1(lastPathV, v, lastPath.From, lastPath.To)
			if err != nil {
				return err
			}
		default:
			parentMap[lastPath.NodeKey] = v
		}
	case []interface{}:
		parentValue := parentV.([]interface{})
		for _, m := range parentValue {
			e := receiver.setValue0(m, v, from, to, paths)
			if e != nil {
				return e
			}
		}
	case []map[string]interface{}:
		parentValue := parentV.([]map[string]interface{})
		for _, m := range parentValue {
			e := receiver.setValue0(m, v, from, to, paths)
			if e != nil {
				return e
			}
		}
	default:
		parentValueType, _ := types.GetType(parentV)
		return types.NewUnSupportError(fmt.Sprintf("%s不支持的类型%s", paths[len(paths)-2].NodeKey, parentValueType))
	}
	return nil

}

func (receiver *Handler) setValue1(lastV interface{}, v interface{}, from uint, to uint) types.ZfError {

	parentVType, _ := types.GetType(lastV)
	if parentVType != types.Array {
		return types.NewUnSupportError(fmt.Sprintf("只支持array格式指定index，当前节点类别为:%s", parentVType))
	}
	switch lastV.(type) {
	case []map[string]interface{}:
		array := lastV.([]map[string]interface{})
		if int(from) > len(array)-1 {
			return types.NewIndexOutOfBoundError3(array, "array", int(from))
		}
		switch v.(type) {
		case map[string]interface{}:
			for i := from; i < to; i++ {
				array[i] = v.(map[string]interface{})
			}
		case map[interface{}]interface{}:
			for i := from; i < to; i++ {
				array[i] = util.ConvertMap2String(v.(map[interface{}]interface{}))
			}
		default:
			return types.NewUnSupportError(fmt.Sprintf("传入参数格式与指定节点格式不符，指定节点格式为%s", types.Object))
		}
	case []interface{}:
		switch v.(type) {
		case interface{}, nil:
			array := lastV.([]interface{})
			for i := from; i < to; i++ {
				array[i] = v
			}
		default:
			return types.NewUnSupportError(fmt.Sprintf("传入参数格式与指定节点格式不符，指定节点格式为普通节点"))
		}
	}
	return nil

}
