package embedded

import (
	"unsafe"
)

// ListLink is a link to the list container
type ListLink[M any] struct {
	prev *M
	next *M
}

func getListLink[T any](obj *T, linkFieldOfs uintptr) *ListLink[T] {
	u := unsafe.Add(unsafe.Pointer(obj), linkFieldOfs)
	return (*ListLink[T])(u)
}
