package embedded

import (
	"unsafe"
)

// This is a combination double-linked list and hash table container - it allows
// for fast lookup via a hash value and linear iteration over its contents.
// This cointainer does not take ownership of its contents, so the application
// must remove items manually.

type HashListMap[TKey HashMapKeyType, T any] interface {
	TableInterface
	ListInterface[TKey, T]

	RemoveAllByKey(key TKey)
	RemoveAllByUniqueKey(key TKey)
}

func NewHashListMapStatic[TKey HashMapKeyType, T any](linkField uintptr, tableSize int) HashListMap[TKey, T] {
	var hlml HashListMapLink[TKey, T]
	return &embeddedHashListMap[TKey, T]{
		hashList:  NewHashListStatic[T](linkField+unsafe.Offsetof(hlml.hashList), tableSize),
		linkField: linkField,
	}
}

func NewHashListMapDynamic[TKey HashMapKeyType, T any](linkField uintptr) HashListMap[TKey, T] {
	var hlml HashListMapLink[TKey, T]
	return &embeddedHashListMap[TKey, T]{
		hashList:  NewHashListDynamic[T](linkField + unsafe.Offsetof(hlml.hashList)),
		linkField: linkField,
	}
}

type embeddedHashListMap[TKey HashMapKeyType, T any] struct {
	hashList  HashList[T]
	linkField uintptr
}

func (c *embeddedHashListMap[TKey, T]) getLink(obj *T) *HashListMapLink[TKey, T] {
	return getHashListMapLink[TKey](obj, c.linkField)
}

func (c *embeddedHashListMap[TKey, T]) IsStatic() bool {
	return c.hashList.IsStatic()
}

func (c *embeddedHashListMap[TKey, T]) First() *T {
	return c.hashList.First()
}

func (c *embeddedHashListMap[TKey, T]) Last() *T {
	return c.hashList.Last()
}

func (c *embeddedHashListMap[TKey, T]) Next(cur *T) *T {
	return c.hashList.Next(cur)
}

func (c *embeddedHashListMap[TKey, T]) Prev(cur *T) *T {
	return c.hashList.Prev(cur)
}

func (c *embeddedHashListMap[TKey, T]) Position(index int) *T {
	return c.hashList.Position(index)
}

func (c *embeddedHashListMap[TKey, T]) Count() int {
	return c.hashList.Count()
}

func (c *embeddedHashListMap[TKey, T]) Remove(obj *T) *T {
	return c.hashList.Remove(obj)
}

func (c *embeddedHashListMap[TKey, T]) RemoveFirst() *T {
	return c.hashList.RemoveFirst()
}

func (c *embeddedHashListMap[TKey, T]) RemoveLast() *T {
	return c.hashList.RemoveLast()
}

func (c *embeddedHashListMap[TKey, T]) RemoveAll() {
	c.hashList.RemoveAll()
}

func (c *embeddedHashListMap[TKey, T]) RemoveAllByKey(key TKey) {
	cur := c.FindFirst(key)
	for cur != nil {
		next := c.FindNext(cur)
		curLink := c.getLink(cur)
		if curLink.key.value == key {
			c.Remove(cur)
		}
		cur = next
	}
}

func (c *embeddedHashListMap[TKey, T]) RemoveAllByUniqueKey(key TKey) {
	cur := c.FindFirst(key)
	for cur != nil {
		next := c.FindNext(cur)
		curLink := c.getLink(cur)
		if curLink.key.value == key {
			c.hashList.Remove(cur)
			return
		}
		cur = next
	}
}

func (c *embeddedHashListMap[TKey, T]) InsertFirst(key TKey, cur *T) *T {
	hashedKey := newHashKey(key)
	obj := c.hashList.InsertFirst(hashedKey.hash, cur)
	if obj == nil {
		return nil
	}

	objLink := c.getLink(obj)
	objLink.key = hashedKey
	return obj
}

func (c *embeddedHashListMap[TKey, T]) InsertLast(key TKey, cur *T) *T {
	hashedKey := newHashKey(key)
	obj := c.hashList.InsertLast(hashedKey.hash, cur)
	if obj == nil {
		return nil
	}

	objLink := c.getLink(obj)
	objLink.key = hashedKey
	return obj
}

func (c *embeddedHashListMap[TKey, T]) InsertAfter(key TKey, prev, cur *T) *T {
	hashedKey := newHashKey(key)
	obj := c.hashList.InsertAfter(hashedKey.hash, prev, cur)
	if obj == nil {
		return nil
	}

	objLink := c.getLink(obj)
	objLink.key = hashedKey
	return obj
}

func (c *embeddedHashListMap[TKey, T]) InsertBefore(key TKey, after, cur *T) *T {
	hashedKey := newHashKey(key)
	obj := c.hashList.InsertBefore(hashedKey.hash, after, cur)
	if obj == nil {
		return nil
	}

	objLink := c.getLink(obj)
	objLink.key = hashedKey
	return obj
}

func (c *embeddedHashListMap[TKey, T]) Move(obj *T, newKey TKey) {
	objLink := c.getLink(obj)
	objLink.key = newHashKey(newKey)
	c.hashList.Move(obj, objLink.key.hash)
}

func (c *embeddedHashListMap[TKey, T]) MoveFirst(cur *T) {
	c.hashList.MoveFirst(cur)
}

func (c *embeddedHashListMap[TKey, T]) MoveLast(cur *T) {
	c.hashList.MoveLast(cur)
}

func (c *embeddedHashListMap[TKey, T]) MoveAfter(dest, cur *T) {
	c.hashList.MoveAfter(dest, cur)
}

func (c *embeddedHashListMap[TKey, T]) MoveBefore(dest, cur *T) {
	c.hashList.MoveBefore(dest, cur)
}

func (c *embeddedHashListMap[TKey, T]) FindFirst(key TKey) *T {
	return c.hashList.FindFirst(HashKey(key))
}

func (c *embeddedHashListMap[TKey, T]) FindNext(prevResult *T) *T {
	return c.hashList.FindNext(prevResult)
}

func (c *embeddedHashListMap[TKey, T]) GetKey(obj *T) TKey {
	return c.getLink(obj).key.value
}

func (c *embeddedHashListMap[TKey, T]) GetTableSize() int {
	return c.hashList.GetTableSize()
}

func (c *embeddedHashListMap[TKey, T]) GetTableUsed() int {
	return c.hashList.GetTableUsed()
}

func (c *embeddedHashListMap[TKey, T]) Reserve(count int) {
	c.hashList.Reserve(count)
}

func (c *embeddedHashListMap[TKey, T]) IsEmpty() bool {
	return c.hashList.IsEmpty()
}

func (c *embeddedHashListMap[TKey, T]) IsContained(cur *T) bool {
	return c.hashList.IsContained(cur)
}
