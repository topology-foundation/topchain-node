package utils

import (
	"unsafe"
)

func Int64ToByteArray(i int64) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(&i)), 8)
}

func ByteArrayToInt(b []byte) int64 {
	return *(*int64)(unsafe.Pointer(&b[0]))
}
