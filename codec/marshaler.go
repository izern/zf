package codec

import "github.com/izern/zf/types"

func init() {

}

type Marshaler interface {
	Marshal(m interface{}) (result []byte, err types.ZfError)
}
