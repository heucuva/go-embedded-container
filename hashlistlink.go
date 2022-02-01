package embedded

import (
	"unsafe"
)

// HashListLink is a link to the map container
type HashListLink[M any] struct {
	hash HashLink[M]
	list ListLink[M]
}

func getHashListLink[T any](obj *T, linkFieldOfs uintptr) *HashListLink[T] {
	u := unsafe.Add(unsafe.Pointer(obj), linkFieldOfs)
	return (*HashListLink[T])(u)
}
