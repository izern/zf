package types

func init() {

}

// TypeCommand 支持的类别
type TypeCommand interface {
	// GetCurrType path规则参考JSONPATH, .标最顶级
	// 返回当前的类别名称
	GetCurrType() string
	// Parse 格式化输出文本
	Parse(text string) (string, ZfError)
	// Marshal toString
	Marshal(content interface{}) (string, ZfError)
	// Keys 返回指定路径下的key，当类别为object时有效
	Keys(from uint, to uint, path string, text string) ([]string, ZfError)
	// GetType 返回指定路径的值的类型
	GetType(path string, text string) (ValueType, ZfError)
	// GetValues 获取指定路径的值，如果类型是array，则支持指定顺序的值
	GetValues(from uint, to uint, path string, text string) (interface{}, ZfError)
	// Append 对指定路径的值进行追加内容，如果类型是object，可以指定key增加键值对,返回更新后的值
	Append(path string, key string, index uint, value string, text string) (string, ZfError)
	// SetValue 对指定路径的值进行覆盖更新，返回更新后的值
	SetValue(path string, value string, text string) (string, ZfError)
}
