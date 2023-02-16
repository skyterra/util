package common

import (
	"reflect"
	"unsafe"
)

// StrToBytes 无内存拷贝，将字符串转换成byte切片
func StrToBytes(s string) []byte {
	b := *(*[]byte)(unsafe.Pointer(&s))
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sliceHeader.Cap = sliceHeader.Len
	return b
}

// BytesToStr 无内存拷贝，将byte切片转换成字符串
func BytesToStr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}