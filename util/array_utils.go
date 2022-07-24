package util

import (
	"github.com/izern/zf/types"
)

func init() {

}

func ArrayCopy(src []interface{}, srcPos int, dest []interface{}, destPos int, length int) types.ZfError {

	srcLength := len(src)
	destLength := len(dest)
	if srcPos > srcLength {
		return types.NewIndexOutOfBoundError(src, "src", srcPos)
	}
	if srcPos+length > srcLength {
		return types.NewIndexOutOfBoundError(src, "src", srcPos+length)
	}

	if destPos > destLength {
		return types.NewIndexOutOfBoundError(src, "dest", srcPos)
	}
	if destPos+length > destLength {
		return types.NewIndexOutOfBoundError(src, "dest", srcPos+length)
	}

	j := destPos

	for i := srcPos; i < srcPos+length; i++ {
		dest[j] = src[i]
		j++
	}
	return nil
}
