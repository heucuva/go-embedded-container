package embedded

import (
	"unsafe"
)

// HashLink is a link to the map container
type HashLink[M any] struct {
	hashNext  *M
	hashValue HashedKeyValue
}

func getHashLink[T any](obj *T, linkFieldOfs uintptr) *HashLink[T] {
	if obj == nil {
		return nil
	}
	u := unsafe.Add(unsafe.Pointer(obj), linkFieldOfs)
	return (*HashLink[T])(u)
}
