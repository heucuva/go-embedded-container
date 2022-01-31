package embedded

import (
	"constraints"
	"fmt"
	"hash/fnv"
	"unsafe"
)

// This is a combination double-linked list and hash table container - it allows
// for fast lookup via a hash value and linear iteration over its contents.
// This cointainer does not take ownership of its contents, so the application
// must remove items manually.

type HashListMap[TKey HashListMapKeyType, T any] interface {
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
	RemoveAllByKey(key TKey)
	RemoveAllByUniqueKey(key TKey)

	InsertFirst(key TKey, cur *T) *T
	InsertLast(key TKey, cur *T) *T
	InsertAfter(key TKey, prev, cur *T) *T
	InsertBefore(key TKey, after, cur *T) *T

	Move(obj *T, newKey TKey)
	MoveFirst(cur *T)
	MoveLast(cur *T)
	MoveAfter(dest, cur *T)
	MoveBefore(dest, cur *T)

	FindFirst(key TKey) *T
	FindNext(prevResult *T) *T
	GetKey(obj *T) int
	GetTableSize() int
	GetTableUsed() int
	Reserve(count int)
	IsEmpty() bool

	IsContained(cur *T) bool
}

type HashListMapKeyType interface {
	constraints.Ordered
}

func NewHashListMapStatic[TKey HashListMapKeyType, T any](linkField uintptr, tableSize int) HashListMap[TKey, T] {
	var hlml HashListMapLink[TKey, T]
	return &embeddedHashListMap[TKey, T]{
		hashList:  NewHashListStatic[T](linkField+unsafe.Offsetof(hlml.hashList), tableSize),
		linkField: linkField,
	}
}

func NewHashListMapDynamic[TKey HashListMapKeyType, T any](linkField uintptr) HashListMap[TKey, T] {
	var hlml HashListMapLink[TKey, T]
	return &embeddedHashListMap[TKey, T]{
		hashList:  NewHashListDynamic[T](linkField + unsafe.Offsetof(hlml.hashList)),
		linkField: linkField,
	}
}

type embeddedHashListMap[TKey HashListMapKeyType, T any] struct {
	hashList  HashList[T]
	linkField uintptr
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
	c.hashList.Remove(obj)
	return c.hashList.Remove(obj)
}

func (c *embeddedHashListMap[TKey, T]) RemoveFirst() *T {
	obj := c.hashList.RemoveFirst()
	if obj != nil {
		c.hashList.Remove(obj)
	}
	return obj
}

func (c *embeddedHashListMap[TKey, T]) RemoveLast() *T {
	obj := c.hashList.RemoveLast()
	if obj != nil {
		c.hashList.Remove(obj)
	}
	return obj
}

func (c *embeddedHashListMap[TKey, T]) RemoveAll() {
	c.hashList.RemoveAll()
}

func (c *embeddedHashListMap[TKey, T]) RemoveAllByKey(key TKey) {
	hashValue := c.getHashValue(key)
	cur := c.hashList.FindFirst(hashValue)
	for cur != nil {
		next := c.hashList.FindNext(cur)
		u := getHashListMapLink[TKey](cur, c.linkField)
		if u.key == key {
			c.hashList.Remove(cur)
		}
		cur = next
	}
}

func (c *embeddedHashListMap[TKey, T]) RemoveAllByUniqueKey(key TKey) {
	hashValue := c.getHashValue(key)
	cur := c.hashList.FindFirst(hashValue)
	for cur != nil {
		next := c.hashList.FindNext(cur)
		u := getHashListMapLink[TKey](cur, c.linkField)
		if u.key == key {
			c.hashList.Remove(cur)
			return
		}
		cur = next
	}
}

func (c *embeddedHashListMap[TKey, T]) InsertFirst(key TKey, cur *T) *T {
	hashValue := c.getHashValue(key)
	obj := c.hashList.InsertFirst(hashValue, cur)
	if obj == nil {
		return nil
	}

	u := getHashListMapLink[TKey](obj, c.linkField)
	u.key = key
	return obj
}

func (c *embeddedHashListMap[TKey, T]) InsertLast(key TKey, cur *T) *T {
	hashValue := c.getHashValue(key)
	obj := c.hashList.InsertLast(hashValue, cur)
	if obj == nil {
		return nil
	}

	u := getHashListMapLink[TKey](obj, c.linkField)
	u.key = key
	return obj
}

func (c *embeddedHashListMap[TKey, T]) InsertAfter(key TKey, prev, cur *T) *T {
	hashValue := c.getHashValue(key)
	obj := c.hashList.InsertAfter(hashValue, prev, cur)
	if obj == nil {
		return nil
	}

	u := getHashListMapLink[TKey](obj, c.linkField)
	u.key = key
	return obj
}

func (c *embeddedHashListMap[TKey, T]) InsertBefore(key TKey, after, cur *T) *T {
	hashValue := c.getHashValue(key)
	obj := c.hashList.InsertBefore(hashValue, after, cur)
	if obj == nil {
		return nil
	}

	u := getHashListMapLink[TKey](obj, c.linkField)
	u.key = key
	return obj
}

func (c *embeddedHashListMap[TKey, T]) Move(obj *T, newKey TKey) {
	newHashValue := c.getHashValue(newKey)
	c.hashList.Move(obj, newHashValue)
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
	hashValue := c.getHashValue(key)
	return c.hashList.FindFirst(hashValue)
}

func (c *embeddedHashListMap[TKey, T]) FindNext(prevResult *T) *T {
	return c.hashList.FindNext(prevResult)
}

func (c *embeddedHashListMap[TKey, T]) GetKey(obj *T) int {
	return c.hashList.GetKey(obj)
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

func (c *embeddedHashListMap[TKey, T]) getHashValue(key TKey) int {
	h := fnv.New64()
	h.Write([]byte(fmt.Sprint(key)))
	return int(h.Sum64())
}

// HashListMapLink is a link to the map container
type HashListMapLink[TKey HashListMapKeyType, T any] struct {
	hashList HashListLink[T]
	key      TKey
}

func getHashListMapLink[TKey HashListMapKeyType, T any](obj *T, linkFieldOfs uintptr) *HashListMapLink[TKey, T] {
	u := unsafe.Add(unsafe.Pointer(obj), linkFieldOfs)
	return (*HashListMapLink[TKey, T])(u)
}
