package util

import (
	"fmt"
	testing2 "testing"
)

func init() {

}

func Test_ParsePath(t *testing2.T) {

	paths, zfError := ParsePath(".[1].a.b[1].c")
	if zfError != nil {
		fmt.Println(zfError.Error())
		panic(zfError)
	}
	for _, path := range paths {
		fmt.Println(path)
	}
}
