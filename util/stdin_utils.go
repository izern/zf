package util

import (
	"bufio"
	"io"
	"os"
)

func init() {

}

// IsPipe s判断当前环境是是否管道模式
func IsPipe() bool {
	fi, _ := os.Stdin.Stat()
	isPipe := fi.Mode()&os.ModeCharDevice == 0
	return isPipe
}

// ReadAll 从reader中读取所有内容，WARNING:
func ReadAll(r io.Reader) (string, error) {
	reader := bufio.NewReader(r)
	var data []byte
	p := make([]byte, 1024)
	for {
		n, err := reader.Read(p)
		if err != nil && err != io.EOF {
			return "", err
		} else if n == 0 {
			break
		} else {
			data = append(data, p[:n]...)
		}
	}
	return string(data), nil

}
