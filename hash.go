package embedded

import (
	"unsafe"
)

// This is a hash table container - it allows for fast lookup via a hash value.
// This cointainer does not take ownership of its contents, so the application
// must remove items manually.

type Hash[T any] interface {
	Insert(hashValue int, obj *T) *T
	Remove(obj *T) *T
	Move(obj *T, newHashValue int)
	Reserve(count int)
	GetKey(obj *T) int
	Count() int
	GetTableSize() int
	GetTableUsed() int
	IsEmpty() bool
	FindFirst(hashValue int) *T
	FindNext(prevResult *T) *T
	WalkFirst() *T
	WalkNext(prevResult *T) *T
	RemoveAll()
	IsContained(cur *T) bool
}

// HashLink is a link to the map container
type HashLink[M any] struct {
	hashNext *M
	hashValue int
}

func getHashLink[T any](obj *T, linkFieldOfs uintptr) *HashLink[T] {
	u := unsafe.Add(unsafe.Pointer(obj), linkFieldOfs)
	return (*HashLink[T])(u)
}

