package embedded

import (
	"unsafe"
)

// MapLink is a link to the map container
type MapLink[TKey, T any] struct {
	key      TKey
	parent   *T
	left     *T
	right    *T
	red      bool
	position int
}

func getMapLink[TKey, T any](obj *T, linkFieldOfs uintptr) *MapLink[TKey, T] {
	u := unsafe.Add(unsafe.Pointer(obj), linkFieldOfs)
	return (*MapLink[TKey, T])(u)
}
