package embedded

import (
	"unsafe"
)

// This is a combination double-linked list and hash table container - it allows
// for fast lookup via a hash value and linear iteration over its contents.
// This cointainer does not take ownership of its contents, so the application
// must remove items manually.

type HashList[T any] interface {
	TableInterface
	ListInterface[HashedKeyValue, T]
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
	list List[T]
	hash Hash[T]
}

func (c *embeddedHashList[T]) IsStatic() bool {
	return c.hash.IsStatic()
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

func (c *embeddedHashList[T]) InsertFirst(hashValue HashedKeyValue, cur *T) *T {
	c.list.InsertFirst(cur)
	return c.hash.Insert(hashValue, cur)
}

func (c *embeddedHashList[T]) InsertLast(hashValue HashedKeyValue, cur *T) *T {
	c.list.InsertLast(cur)
	return c.hash.Insert(hashValue, cur)
}

func (c *embeddedHashList[T]) InsertAfter(hashValue HashedKeyValue, prev, cur *T) *T {
	c.list.InsertAfter(prev, cur)
	return c.hash.Insert(hashValue, cur)
}

func (c *embeddedHashList[T]) InsertBefore(hashValue HashedKeyValue, after, cur *T) *T {
	c.list.InsertBefore(after, cur)
	return c.hash.Insert(hashValue, cur)
}

func (c *embeddedHashList[T]) Move(obj *T, newHashValue HashedKeyValue) {
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

func (c *embeddedHashList[T]) FindFirst(hashValue HashedKeyValue) *T {
	return c.hash.FindFirst(hashValue)
}

func (c *embeddedHashList[T]) FindNext(prevResult *T) *T {
	return c.hash.FindNext(prevResult)
}

func (c *embeddedHashList[T]) GetKey(obj *T) HashedKeyValue {
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
