package embedded

import (
	"unsafe"
)

// This is a hash map container - it allows for fast lookups of its contents.
// This cointainer does not take ownership of its contents, so the application
// must remove items manually.

type HashMap[TKey HashMapKeyType, T any] interface {
	TableInterface

	Remove(obj *T) *T

	Insert(key TKey, obj *T) *T

	Move(obj *T, newKey TKey)

	FindFirst(key TKey) *T
	FindNext(prevResult *T) *T
	GetKey(obj *T) TKey

	IsContained(obj *T) bool

	RemoveAllByKey(key TKey)
	RemoveAllByUniqueKey(key TKey)

	WalkFirst() *T
	WalkNext(prevResult *T) *T
}

func NewHashMapStatic[TKey HashMapKeyType, T any](linkField uintptr, tableSize int) HashMap[TKey, T] {
	var hml HashMapLink[TKey, T]
	return &embeddedHashMap[TKey, T]{
		hash:      NewHashStatic[T](linkField+unsafe.Offsetof(hml.link), tableSize),
		linkField: linkField,
	}
}

func NewHashMapDynamic[TKey HashMapKeyType, T any](linkField uintptr) HashMap[TKey, T] {
	var hml HashMapLink[TKey, T]
	return &embeddedHashMap[TKey, T]{
		hash:      NewHashDynamic[T](linkField + unsafe.Offsetof(hml.link)),
		linkField: linkField,
	}
}

type embeddedHashMap[TKey HashMapKeyType, T any] struct {
	hash      Hash[T]
	linkField uintptr
}

func (c *embeddedHashMap[TKey, T]) getLink(obj *T) *HashMapLink[TKey, T] {
	return getHashMapLink[TKey](obj, c.linkField)
}

func (c *embeddedHashMap[TKey, T]) Insert(key TKey, obj *T) *T {
	hashedKey := newHashKey(key)
	o := c.hash.Insert(hashedKey.hash, obj)
	if o == nil {
		return nil
	}
	oLink := c.getLink(o)
	oLink.key = hashedKey
	return o
}

func (c *embeddedHashMap[TKey, T]) Remove(obj *T) *T {
	return c.hash.Remove(obj)
}

func (c *embeddedHashMap[TKey, T]) Move(obj *T, newKey TKey) {
	hashedKey := newHashKey(newKey)
	c.hash.Move(obj, hashedKey.hash)
	if obj == nil {
		return
	}
	objLink := c.getLink(obj)
	objLink.key = hashedKey
}

func (c *embeddedHashMap[TKey, T]) RemoveAll() {
	c.hash.RemoveAll()
}

func (c *embeddedHashMap[TKey, T]) RemoveAllByKey(key TKey) {
	hashedKey := newHashKey(key)
	cur := c.hash.FindFirst(hashedKey.hash)
	for cur != nil {
		next := c.hash.FindNext(cur)
		curLink := c.getLink(cur)
		if curLink.key.value == key {
			c.hash.Remove(cur)
		}
		cur = next
	}
}

func (c *embeddedHashMap[TKey, T]) RemoveAllByUniqueKey(key TKey) {
	hashedKey := newHashKey(key)
	cur := c.hash.FindFirst(hashedKey.hash)
	for cur != nil {
		next := c.hash.FindNext(cur)
		curLink := c.getLink(cur)
		if curLink.key.value == key {
			c.hash.Remove(cur)
			return
		}
		cur = next
	}
}

func (c *embeddedHashMap[TKey, T]) Reserve(count int) {
	c.hash.Reserve(count)
}

func (c *embeddedHashMap[TKey, T]) GetKey(obj *T) TKey {
	objLink := c.getLink(obj)
	return objLink.key.value
}

func (c *embeddedHashMap[TKey, T]) Count() int {
	return c.hash.Count()
}

func (c *embeddedHashMap[TKey, T]) GetTableSize() int {
	return c.hash.GetTableSize()
}

func (c *embeddedHashMap[TKey, T]) GetTableUsed() int {
	return c.hash.GetTableUsed()
}

func (c *embeddedHashMap[TKey, T]) IsEmpty() bool {
	return c.hash.IsEmpty()
}

func (c *embeddedHashMap[TKey, T]) FindFirst(key TKey) *T {
	hashedKey := newHashKey(key)
	cur := c.hash.FindFirst(hashedKey.hash)
	for cur != nil {
		next := c.hash.FindNext(cur)
		curLink := c.getLink(cur)
		if curLink.key.value == key {
			return cur
		}
		cur = next
	}
	return nil
}

func (c *embeddedHashMap[TKey, T]) FindNext(prevResult *T) *T {
	prevResultLink := c.getLink(prevResult)
	cur := c.hash.FindNext(prevResult)
	for cur != nil {
		next := c.hash.FindNext(cur)
		curLink := c.getLink(cur)
		if prevResultLink.key.value == curLink.key.value {
			return cur
		}
		cur = next
	}
	return nil
}

func (c *embeddedHashMap[TKey, T]) WalkFirst() *T {
	return c.hash.WalkFirst()
}

func (c *embeddedHashMap[TKey, T]) WalkNext(prevResult *T) *T {
	return c.hash.WalkNext(prevResult)
}

func (c *embeddedHashMap[TKey, T]) IsContained(obj *T) bool {
	return c.hash.IsContained(obj)
}
