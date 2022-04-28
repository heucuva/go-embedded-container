package embedded

import (
	"unsafe"
)

// HashMapLink is a link to the list container
type HashMapLink[TKey HashMapKeyType, M any] struct {
	link HashLink[M]
	key  hashKey[TKey]
}

func getHashMapLink[TKey HashMapKeyType, T any](obj *T, linkFieldOfs uintptr) *HashMapLink[TKey, T] {
	if obj == nil {
		return nil
	}
	u := unsafe.Add(unsafe.Pointer(obj), linkFieldOfs)
	return (*HashMapLink[TKey, T])(u)
}
