package codec

import "github.com/izern/zf/types"

func init() {

}

type Unmarshaler interface {
	Unmarshal(b []byte) (result interface{}, err types.ZfError)
}
