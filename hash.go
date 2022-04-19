package embedded

import (
	"github.com/heucuva/go-embedded-container/internal/array"
)

// This is a hash table container - it allows for fast lookup via a hash value.
// This cointainer does not take ownership of its contents, so the application
// must remove items manually.

type Hash[T any] interface {
	TableInterface

	Remove(obj *T) *T

	Insert(hashValue HashedKeyValue, obj *T) *T

	Move(obj *T, newHashValue HashedKeyValue)

	FindFirst(hashValue HashedKeyValue) *T
	FindNext(prevResult *T) *T
	GetKey(obj *T) HashedKeyValue

	IsContained(obj *T) bool

	WalkFirst() *T
	WalkNext(prevResult *T) *T
}

func NewHashStatic[T any](linkField uintptr, tableSize int) Hash[T] {
	return &embeddedHash[T]{
		linkField: linkField,
		table:     array.NewStaticArray[*T](tableSize),
	}
}

func NewHashDynamic[T any](linkField uintptr) Hash[T] {
	h := &embeddedHash[T]{
		linkField: linkField,
	}
	h.table = array.NewDynamicArray(minDynamicHashSize, h.onResize)
	return h
}

const (
	minDynamicHashSize = 8
)

type embeddedHash[T any] struct {
	table      array.Array[*T]
	linkField  uintptr
	entryCount int
}

func (c *embeddedHash[T]) getLink(obj *T) *HashLink[T] {
	return getHashLink(obj, c.linkField)
}

func (c *embeddedHash[T]) calcSpot(hashValue HashedKeyValue) int {
	return c.calcSpotForSize(hashValue, c.table.Size())
}

func (c *embeddedHash[T]) calcSpotForSize(hashValue HashedKeyValue, tableSize int) int {
	return int(hashValue % HashedKeyValue(tableSize))
}

func (c *embeddedHash[T]) IsStatic() bool {
	return c.table.IsStatic()
}

func (c *embeddedHash[T]) Insert(hashValue HashedKeyValue, obj *T) *T {
	if !c.table.IsStatic() {
		c.Reserve(c.entryCount + 1)
	}
	spot := c.calcSpot(hashValue)
	entryLink := c.getLink(obj)
	entryLink.hashValue = hashValue
	entryLink.hashNext = c.table.At(spot)
	c.table.Set(spot, obj)
	c.entryCount++
	return obj
}

func (c *embeddedHash[T]) Remove(obj *T) *T {
	spot := c.calcSpot(c.getLink(obj).hashValue)
	prev := c.table.Ptr(spot)
	cur := *prev

	for cur != nil {
		entryLink := c.getLink(cur)
		if cur == obj {
			*prev = entryLink.hashNext
			entryLink.hashNext = nil
			entryLink.hashValue = 0
			c.entryCount--
			return cur
		}
		prev = &entryLink.hashNext
		cur = entryLink.hashNext
	}
	return nil
}

func (c *embeddedHash[T]) Move(obj *T, newHashValue HashedKeyValue) {
	c.Remove(obj)
	c.Insert(newHashValue, obj)
}

func (c *embeddedHash[T]) Reserve(count int) {
	c.table.Reserve(count)
}

func (c *embeddedHash[T]) GetKey(obj *T) HashedKeyValue {
	return c.getLink(obj).hashValue
}

func (c *embeddedHash[T]) Count() int {
	return c.entryCount
}

func (c *embeddedHash[T]) GetTableSize() int {
	return c.table.Size()
}

func (c *embeddedHash[T]) GetTableUsed() int {
	if c.entryCount <= 1 {
		return c.entryCount
	}

	var tableUsed int
	for _, entry := range c.table.Slice() {
		if entry != nil {
			tableUsed++
		}
	}
	return tableUsed
}

func (c *embeddedHash[T]) IsEmpty() bool {
	return c.entryCount == 0
}

func (c *embeddedHash[T]) FindFirst(hashValue HashedKeyValue) *T {
	spot := c.calcSpot(hashValue)
	var next *T
	for entry := c.table.At(spot); entry != nil; entry = next {
		entryLink := c.getLink(entry)
		if entryLink.hashValue == hashValue {
			return entry
		}
		next = entryLink.hashNext
	}
	return nil
}

func (c *embeddedHash[T]) FindNext(prevResult *T) *T {
	entry := prevResult
	entryLink := c.getLink(entry)
	if entryLink == nil {
		return nil
	}
	hashValue := entryLink.hashValue
	entry = entryLink.hashNext
	for entry != nil {
		entryLink = c.getLink(entry)
		if entryLink.hashValue == hashValue {
			return entry
		}
		entry = entryLink.hashNext
	}
	return nil
}

func (c *embeddedHash[T]) WalkFirst() *T {
	if c.entryCount == 0 {
		return nil
	}

	for _, entry := range c.table.Slice() {
		if entry != nil {
			return entry
		}
	}
	return nil
}

func (c *embeddedHash[T]) WalkNext(prevResult *T) *T {
	entry := prevResult
	entryLink := c.getLink(entry)
	entry = entryLink.hashNext
	if entry != nil {
		return entry
	}

	spot := c.calcSpot(entryLink.hashValue)
	for spot++; spot < c.table.Size(); spot++ {
		entry = c.table.At(spot)
		if entry != nil {
			return entry
		}
	}
	return nil
}

func (c *embeddedHash[T]) RemoveAll() {
	cur := c.WalkFirst()
	for cur != nil {
		next := c.WalkNext(cur)
		c.Remove(cur)
		cur = next
	}
}

func (c *embeddedHash[T]) IsContained(cur *T) bool {
	curLink := c.getLink(cur)
	if curLink == nil {
		return false
	}
	if curLink.hashValue != 0 || curLink.hashNext != nil {
		return true
	}

	for obj := c.FindFirst(curLink.hashValue); obj != nil; obj = c.FindNext(obj) {
		if obj == cur {
			return true
		}
	}
	return false
}

func (c *embeddedHash[T]) onResize(dest, src []*T) {
	if c.entryCount == 0 {
		return
	}

	dynamicSize := len(dest)
	for _, current := range src {
		if current == nil {
			continue
		}

		var tempBucketRoot *T
		for current != nil {
			currentLink := c.getLink(current)
			next := currentLink.hashNext
			currentLink.hashNext = tempBucketRoot
			tempBucketRoot = current
			current = next
		}

		current = tempBucketRoot
		for current != nil {
			currentLink := c.getLink(current)
			next := currentLink.hashNext
			spot := c.calcSpotForSize(currentLink.hashValue, dynamicSize)
			currentLink.hashNext = dest[spot]
			dest[spot] = current
			current = next
		}
	}
}
