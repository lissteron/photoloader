package helpers

import "unsafe"

func IsNil(v any) bool {
	return (*[2]uintptr)(unsafe.Pointer(&v))[1] == 0
}

func IsNotNil(v any) bool {
	return !IsNil(v)
}
