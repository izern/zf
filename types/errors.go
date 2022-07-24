package types

import (
	"errors"
	"fmt"
)

func init() {

}

type ZfError interface {
	Error() error
}

type KeyNotFoundError struct {
	K string
}

type FormatError struct {
	V    string
	Type string
}

type UnSupportError struct {
	Type string
}

type IndexOutOfBoundError struct {
	size      int
	index     int
	arrayName string
}

func NewKeyNotFoundError(key string) *KeyNotFoundError {
	return &KeyNotFoundError{K: key}
}

func (err *KeyNotFoundError) Error() error {
	return errors.New("找不到键:" + err.K)
}

func NewFormatError(value string, typeStr string) *FormatError {
	return &FormatError{V: value, Type: typeStr}
}

func (err *FormatError) Error() error {
	return errors.New(fmt.Sprintf("无法将解析化为%s格式, %s", err.Type, err.V))
}

func NewUnSupportError(typeStr string) *UnSupportError {
	return &UnSupportError{Type: typeStr}
}

func (err *UnSupportError) Error() error {
	return errors.New(fmt.Sprintf("不支持的操作：%s", err.Type))
}

func NewIndexOutOfBoundError(array []interface{}, arrayName string, index int) *IndexOutOfBoundError {
	return &IndexOutOfBoundError{index: index, size: len(array), arrayName: arrayName}
}
func NewIndexOutOfBoundError2(array []map[string]interface{}, arrayName string, index int) *IndexOutOfBoundError {
	return &IndexOutOfBoundError{index: index, size: len(array), arrayName: arrayName}
}

func NewIndexOutOfBoundError3(array []map[string]interface{}, arrayName string, index int) *IndexOutOfBoundError {
	return &IndexOutOfBoundError{index: index, size: len(array), arrayName: arrayName}
}

func (err *IndexOutOfBoundError) Error() error {
	return errors.New(fmt.Sprintf("数组%s越界，最大:%d，请求值:%d", err.arrayName, err.size, err.index))
}
