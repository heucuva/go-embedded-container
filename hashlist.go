package embedded

import (
	"unsafe"
)

// This is a combination double-linked list and hash table container - it allows
// for fast lookup via a hash value and linear iteration over its contents.
// This cointainer does not take ownership of its contents, so the application
// must remove items manually.

type HashList[T any] interface {
	First() *T
	Last() *T
	Next(cur *T) *T
	Prev(cur *T) *T
	Position(index int) *T
	Count() int

	Remove(obj *T) *T
	RemoveFirst() *T
	RemoveLast() *T
	RemoveAll()

	InsertFirst(hashValue int, cur *T) *T
	InsertLast(hashValue int, cur *T) *T
	InsertAfter(hashValue int, prev, cur *T) *T
	InsertBefore(hashValue int, after, cur *T) *T

	Move(obj *T, newHashValue int)
	MoveFirst(cur *T)
	MoveLast(cur *T)
	MoveAfter(dest, cur *T)
	MoveBefore(dest, cur *T)

	FindFirst(hashValue int) *T
	FindNext(prevResult *T) *T
	GetKey(obj *T) int
	GetTableSize() int
	GetTableUsed() int
	Reserve(count int)
	IsEmpty() bool

	IsContained(cur *T) bool
}

func NewHashListStatic[T any](linkField uintptr, tableSize int) HashList[T] {
	var hll HashListLink[T]
	return &embeddedHashList[T]{
		hash: NewHashStatic[T](linkField+unsafe.Offsetof(hll.hash), tableSize),
		list: NewList[T](linkField + unsafe.Offsetof(hll.list)),
	}
}

func NewHashListDynamic[T any](linkField uintptr) HashList[T] {
	var hll HashListLink[T]
	return &embeddedHashList[T]{
		hash: NewHashDynamic[T](linkField + unsafe.Offsetof(hll.hash)),
		list: NewList[T](linkField + unsafe.Offsetof(hll.list)),
	}
}

type embeddedHashList[T any] struct {
	hash Hash[T]
	list List[T]
}

func (c *embeddedHashList[T]) First() *T {
	return c.list.First()
}

func (c *embeddedHashList[T]) Last() *T {
	return c.list.Last()
}

func (c *embeddedHashList[T]) Next(cur *T) *T {
	return c.list.Next(cur)
}

func (c *embeddedHashList[T]) Prev(cur *T) *T {
	return c.list.Prev(cur)
}

func (c *embeddedHashList[T]) Position(index int) *T {
	return c.list.Position(index)
}

func (c *embeddedHashList[T]) Count() int {
	return c.hash.Count()
}

func (c *embeddedHashList[T]) Remove(obj *T) *T {
	c.hash.Remove(obj)
	return c.list.Remove(obj)
}

func (c *embeddedHashList[T]) RemoveFirst() *T {
	obj := c.list.RemoveFirst()
	if obj != nil {
		c.hash.Remove(obj)
	}
	return obj
}

func (c *embeddedHashList[T]) RemoveLast() *T {
	obj := c.list.RemoveLast()
	if obj != nil {
		c.hash.Remove(obj)
	}
	return obj
}

func (c *embeddedHashList[T]) RemoveAll() {
	c.list.RemoveAll()
	c.hash.RemoveAll()
}

func (c *embeddedHashList[T]) InsertFirst(hashValue int, cur *T) *T {
	c.list.InsertFirst(cur)
	return c.hash.Insert(hashValue, cur)
}

func (c *embeddedHashList[T]) InsertLast(hashValue int, cur *T) *T {
	c.list.InsertLast(cur)
	return c.hash.Insert(hashValue, cur)
}

func (c *embeddedHashList[T]) InsertAfter(hashValue int, prev, cur *T) *T {
	c.list.InsertAfter(prev, cur)
	return c.hash.Insert(hashValue, cur)
}

func (c *embeddedHashList[T]) InsertBefore(hashValue int, after, cur *T) *T {
	c.list.InsertBefore(after, cur)
	return c.hash.Insert(hashValue, cur)
}

func (c *embeddedHashList[T]) Move(obj *T, newHashValue int) {
	c.hash.Move(obj, newHashValue)
}

func (c *embeddedHashList[T]) MoveFirst(cur *T) {
	c.list.MoveFirst(cur)
}

func (c *embeddedHashList[T]) MoveLast(cur *T) {
	c.list.MoveLast(cur)
}

func (c *embeddedHashList[T]) MoveAfter(dest, cur *T) {
	c.list.MoveAfter(dest, cur)
}

func (c *embeddedHashList[T]) MoveBefore(dest, cur *T) {
	c.list.MoveBefore(dest, cur)
}

func (c *embeddedHashList[T]) FindFirst(hashValue int) *T {
	return c.hash.FindFirst(hashValue)
}

func (c *embeddedHashList[T]) FindNext(prevResult *T) *T {
	return c.hash.FindNext(prevResult)
}

func (c *embeddedHashList[T]) GetKey(obj *T) int {
	return c.hash.GetKey(obj)
}

func (c *embeddedHashList[T]) GetTableSize() int {
	return c.hash.GetTableSize()
}

func (c *embeddedHashList[T]) GetTableUsed() int {
	return c.hash.GetTableUsed()
}

func (c *embeddedHashList[T]) Reserve(count int) {
	c.hash.Reserve(count)
}

func (c *embeddedHashList[T]) IsEmpty() bool {
	return c.hash.IsEmpty()
}

func (c *embeddedHashList[T]) IsContained(cur *T) bool {
	return c.hash.IsContained(cur)
}

// HashListLink is a link to the map container
type HashListLink[M any] struct {
	hash HashLink[M]
	list ListLink[M]
}

func getHashListLink[T any](obj *T, linkFieldOfs uintptr) *HashListLink[T] {
	u := unsafe.Add(unsafe.Pointer(obj), linkFieldOfs)
	return (*HashListLink[T])(u)
}
