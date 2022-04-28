package embedded

import (
	"unsafe"
)

// HashListMapLink is a link to the map container
type HashListMapLink[TKey HashMapKeyType, T any] struct {
	hashList HashListLink[T]
	key      hashKey[TKey]
}

func getHashListMapLink[TKey HashMapKeyType, T any](obj *T, linkFieldOfs uintptr) *HashListMapLink[TKey, T] {
	if obj == nil {
		return nil
	}
	u := unsafe.Add(unsafe.Pointer(obj), linkFieldOfs)
	return (*HashListMapLink[TKey, T])(u)
}
