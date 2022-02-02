package embedded

import (
	"github.com/heucuva/go-embedded-container/internal/array"
)

// This is a hash table container - it allows for fast lookup via a hash value.
// This cointainer does not take ownership of its contents, so the application
// must remove items manually.

type Hash[T any] interface {
	Insert(hashValue HashedKeyValue, obj *T) *T
	Remove(obj *T) *T
	Move(obj *T, newHashValue HashedKeyValue)
	Reserve(count int)
	GetKey(obj *T) HashedKeyValue
	Count() int
	GetTableSize() int
	GetTableUsed() int
	IsEmpty() bool
	FindFirst(hashValue HashedKeyValue) *T
	FindNext(prevResult *T) *T
	WalkFirst() *T
	WalkNext(prevResult *T) *T
	RemoveAll()
	IsContained(cur *T) bool
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
	entryCount int
	linkField  uintptr
	table      array.Array[*T]
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

func (c *embeddedHash[T]) Insert(hashValue HashedKeyValue, obj *T) *T {
	if !c.table.IsStatic() {
		c.Reserve(c.entryCount + 1)
	}
	spot := c.calcSpot(hashValue)
	objU := c.getLink(obj)
	objU.hashValue = hashValue
	objU.hashNext = c.table.Slice()[spot]
	c.table.Slice()[spot] = obj
	c.entryCount++
	return obj
}

func (c *embeddedHash[T]) Remove(obj *T) *T {
	spot := c.calcSpot(c.getLink(obj).hashValue)
	cur := c.table.Slice()[spot]
	prev := &c.table.Slice()[spot]

	for cur != nil {
		objU := c.getLink(cur)
		if cur == obj {
			*prev = objU.hashNext
			objU.hashNext = nil
			objU.hashValue = 0
			c.entryCount--
			return cur
		}
		prev = &objU.hashNext
		cur = objU.hashNext
	}
	return nil
}

func (c *embeddedHash[T]) Move(obj *T, newHashValue HashedKeyValue) {
	c.Remove(obj)
	c.Insert(newHashValue, obj)
}

func (c *embeddedHash[T]) Reserve(count int) {
	if c.table.IsStatic() {
		panic("cannot reserve a count with a static table size")
	} else {
		c.table.Reserve(count)
	}
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
	entry := c.table.Slice()[spot]
	for entry != nil {
		objU := c.getLink(entry)
		if objU.hashValue == hashValue {
			return entry
		}
		entry = objU.hashNext
	}
	return nil
}

func (c *embeddedHash[T]) FindNext(prevResult *T) *T {
	entry := prevResult
	encryU := c.getLink(entry)
	hashValue := encryU.hashValue
	entry = encryU.hashNext
	for entry != nil {
		encryU = c.getLink(entry)
		if encryU.hashValue == hashValue {
			return entry
		}
		entry = encryU.hashNext
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
	entryU := c.getLink(entry)
	spot := c.calcSpot(entryU.hashValue)
	entry = entryU.hashNext
	if entry != nil {
		return entry
	}

	for spot++; spot < c.table.Size(); spot++ {
		entry = c.table.Slice()[spot]
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
	curU := c.getLink(cur)
	if curU.hashValue != 0 || curU.hashNext != nil {
		return true
	}

	walk := c.table.Slice()[0]
	for walk != nil {
		if walk == cur {
			return true
		}
		walk = c.getLink(walk).hashNext
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
			currentU := c.getLink(current)
			next := currentU.hashNext
			currentU.hashNext = tempBucketRoot
			tempBucketRoot = current
			current = next
		}

		current = tempBucketRoot
		for current != nil {
			currentU := c.getLink(current)
			next := currentU.hashNext
			spot := c.calcSpotForSize(currentU.hashValue, dynamicSize)
			currentU.hashNext = dest[spot]
			dest[spot] = current
			current = next
		}
	}
}
