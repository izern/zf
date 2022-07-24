package types

func init() {

}

type PathType uint8

const (
	_          PathType = iota
	RootNode            // 根节点，一般是$
	NormalNode          // 普通定位节点
	IndexNode           // 索引，如a[1]
	RangeNode           // 范围索引，如 a[0,20]

)

type Path struct {
	OriginValue string
	NodeKey     string
	Type        PathType
	From        uint // 类型为RANGE_NODE时有值
	To          uint // 类型为RANGE_NODE时有值
	Index       uint // 类型为INDEX_NODE时有值
}

func (receiver PathType) IsSupportValue(v ValueType) bool {
	if receiver == RootNode {
		return true
	}
	if receiver == NormalNode {
		return true
	}

	if receiver == IndexNode && (v == Array) {
		return true
	}
	if receiver == RangeNode && (v == Array) {
		return true
	}
	return false
}
